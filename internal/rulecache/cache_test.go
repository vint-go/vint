package rulecache

import (
	"path/filepath"
	"testing"
)

func TestRuleCache_RegisterAndGet(t *testing.T) {
	rc := New()

	// Register a rule.
	if err := rc.RegisterRule("test-rule", []any{"arg1", 42}); err != nil {
		t.Fatalf("RegisterRule failed: %v", err)
	}

	// Cache file hash.
	content := []byte("package main\nfunc main() {}\n")
	rc.CacheFileHash("/tmp/test.go", content)

	// Miss before Put.
	_, hit := rc.Get("test-rule", "/tmp/test.go", TierFileOnly, nil, nil)
	if hit {
		t.Fatal("expected cache miss before Put")
	}

	// Put and then Get.
	failures := []CachedFailure{
		{
			Message:       "found an issue",
			RuleName:      "test-rule",
			Category:      "style",
			StartFilename: "/tmp/test.go",
			StartLine:     2,
			StartColumn:   1,
			Confidence:    1.0,
		},
	}
	rc.Put("test-rule", "/tmp/test.go", TierFileOnly, nil, nil, failures)

	got, hit := rc.Get("test-rule", "/tmp/test.go", TierFileOnly, nil, nil)
	if !hit {
		t.Fatal("expected cache hit after Put")
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 failure, got %d", len(got))
	}
	if got[0].Message != "found an issue" {
		t.Errorf("expected message %q, got %q", "found an issue", got[0].Message)
	}
}

func TestRuleCache_DifferentConfig(t *testing.T) {
	rc := New()

	content := []byte("package main\n")

	// Register with config A.
	if err := rc.RegisterRule("rule", []any{"configA"}); err != nil {
		t.Fatal(err)
	}
	rc.CacheFileHash("/tmp/test.go", content)
	rc.Put("rule", "/tmp/test.go", TierFileOnly, nil, nil, []CachedFailure{{Message: "A"}})

	got, hit := rc.Get("rule", "/tmp/test.go", TierFileOnly, nil, nil)
	if !hit {
		t.Fatal("expected hit")
	}
	if got[0].Message != "A" {
		t.Errorf("expected A, got %s", got[0].Message)
	}

	// Re-register with config B — should miss because config hash changed.
	if err := rc.RegisterRule("rule", []any{"configB"}); err != nil {
		t.Fatal(err)
	}

	_, hit = rc.Get("rule", "/tmp/test.go", TierFileOnly, nil, nil)
	if hit {
		t.Fatal("expected cache miss after config change")
	}
}

func TestRuleCache_PackageAwareTier(t *testing.T) {
	rc := New()

	if err := rc.RegisterRule("rule", nil); err != nil {
		t.Fatal(err)
	}

	content1 := []byte("package p\nvar X = 1\n")
	content2 := []byte("package p\nvar Y = 2\n")

	rc.CacheFileHash("/tmp/a.go", content1)
	rc.CacheFileHash("/tmp/b.go", content2)

	siblings := map[string][]byte{
		"/tmp/a.go": content1,
		"/tmp/b.go": content2,
	}

	// Put with package-aware tier.
	rc.Put("rule", "/tmp/a.go", TierPackageAware, siblings, nil, []CachedFailure{{Message: "pkg"}})

	got, hit := rc.Get("rule", "/tmp/a.go", TierPackageAware, siblings, nil)
	if !hit {
		t.Fatal("expected hit")
	}
	if got[0].Message != "pkg" {
		t.Errorf("expected pkg, got %s", got[0].Message)
	}

	// Simulate a second lint run where sibling content has changed.
	// A new cache instance means CacheFileHash stores the updated hash.
	rc2 := New()
	if err := rc2.RegisterRule("rule", nil); err != nil {
		t.Fatal(err)
	}
	rc2.CacheFileHash("/tmp/a.go", content1)
	content2Changed := []byte("package p\nvar Y = 3\n")
	rc2.CacheFileHash("/tmp/b.go", content2Changed)

	siblings2 := map[string][]byte{
		"/tmp/a.go": content1,
		"/tmp/b.go": content2Changed,
	}

	// Different sibling hash → different action ID → cache miss.
	_, hit = rc2.Get("rule", "/tmp/a.go", TierPackageAware, siblings2, nil)
	if hit {
		t.Fatal("expected miss after sibling change in new run")
	}
}

func TestRuleCache_UnregisteredRule(t *testing.T) {
	rc := New()
	rc.CacheFileHash("/tmp/test.go", []byte("package main\n"))

	_, hit := rc.Get("unregistered", "/tmp/test.go", TierFileOnly, nil, nil)
	if hit {
		t.Fatal("expected miss for unregistered rule")
	}
}

func TestRuleCache_EmptyFailures(t *testing.T) {
	rc := New()

	if err := rc.RegisterRule("clean-rule", nil); err != nil {
		t.Fatal(err)
	}
	rc.CacheFileHash("/tmp/clean.go", []byte("package clean\n"))

	// Cache empty result (file has no issues).
	rc.Put("clean-rule", "/tmp/clean.go", TierFileOnly, nil, nil, []CachedFailure{})

	got, hit := rc.Get("clean-rule", "/tmp/clean.go", TierFileOnly, nil, nil)
	if !hit {
		t.Fatal("expected hit for empty result")
	}
	if len(got) != 0 {
		t.Errorf("expected 0 failures, got %d", len(got))
	}
}

func TestRuleCache_DiskPersistence(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "cache")

	content := []byte("package main\nfunc main() {}\n")
	failures := []CachedFailure{{Message: "disk issue", RuleName: "rule", Confidence: 1.0}}

	// Write to disk via first cache instance.
	rc1, err := NewWithDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := rc1.RegisterRule("rule", nil); err != nil {
		t.Fatal(err)
	}
	rc1.CacheFileHash("/tmp/test.go", content)
	rc1.Put("rule", "/tmp/test.go", TierFileOnly, nil, nil, failures)
	if err := rc1.FlushAll(); err != nil {
		t.Fatal(err)
	}

	// New cache instance pointing at the same directory — should find on disk.
	rc2, err := NewWithDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := rc2.RegisterRule("rule", nil); err != nil {
		t.Fatal(err)
	}
	rc2.CacheFileHash("/tmp/test.go", content)

	got, hit := rc2.Get("rule", "/tmp/test.go", TierFileOnly, nil, nil)
	if !hit {
		t.Fatal("expected disk cache hit from second instance")
	}
	if len(got) != 1 || got[0].Message != "disk issue" {
		t.Errorf("unexpected result: %+v", got)
	}
}

func TestRuleCache_DiskMissOnContentChange(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "cache")

	content := []byte("package main\n")
	failures := []CachedFailure{{Message: "old"}}

	rc1, err := NewWithDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := rc1.RegisterRule("rule", nil); err != nil {
		t.Fatal(err)
	}
	rc1.CacheFileHash("/tmp/test.go", content)
	rc1.Put("rule", "/tmp/test.go", TierFileOnly, nil, nil, failures)
	if err := rc1.FlushAll(); err != nil {
		t.Fatal(err)
	}

	// Second instance with changed file content.
	rc2, err := NewWithDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := rc2.RegisterRule("rule", nil); err != nil {
		t.Fatal(err)
	}
	rc2.CacheFileHash("/tmp/test.go", []byte("package main\n// changed\n"))

	_, hit := rc2.Get("rule", "/tmp/test.go", TierFileOnly, nil, nil)
	if hit {
		t.Fatal("expected miss after file content change")
	}
}

func TestCacheFileHash_ReturnsCachedOnSubsequentCall(t *testing.T) {
	rc := New()

	content := []byte("package main\n")
	h1 := rc.CacheFileHash("/tmp/test.go", content)

	// Second call with different content returns the cached hash (first wins).
	h2 := rc.CacheFileHash("/tmp/test.go", []byte("different content"))
	if h1 != h2 {
		t.Error("expected cached hash to be returned on subsequent call")
	}

	// Same content should still produce the same hash.
	h3 := rc.CacheFileHash("/tmp/test.go", content)
	if h1 != h3 {
		t.Error("expected same hash for same content")
	}

	// Different path should produce a different hash.
	h4 := rc.CacheFileHash("/tmp/other.go", []byte("different content"))
	if h1 == h4 {
		t.Error("expected different hash for different path with different content")
	}
}

func TestRuleCache_PackageLevelCache(t *testing.T) {
	rc := New()

	if err := rc.RegisterRule("rule-a", nil); err != nil {
		t.Fatal(err)
	}
	if err := rc.RegisterRule("rule-b", []any{"x"}); err != nil {
		t.Fatal(err)
	}

	content1 := []byte("package p\nvar X = 1\n")
	content2 := []byte("package p\nvar Y = 2\n")
	rc.CacheFileHash("/tmp/a.go", content1)
	rc.CacheFileHash("/tmp/b.go", content2)

	configHash := [32]byte{1, 2, 3} // arbitrary fixed config hash

	sortedFiles := []string{"/tmp/a.go", "/tmp/b.go"}
	pkgKey := "/tmp"

	pkgID, ok := rc.PackageActionID(sortedFiles, configHash)
	if !ok {
		t.Fatal("expected PackageActionID to succeed")
	}

	// Miss before Put.
	_, hit := rc.GetPackage(pkgID, pkgKey)
	if hit {
		t.Fatal("expected package cache miss before PutPackage")
	}

	// Put and then Get.
	failures := []CachedFailure{
		{Message: "issue1", RuleName: "rule-a"},
		{Message: "issue2", RuleName: "rule-b"},
	}
	rc.PutPackage(pkgID, pkgKey, failures)

	got, hit := rc.GetPackage(pkgID, pkgKey)
	if !hit {
		t.Fatal("expected package cache hit after PutPackage")
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 failures, got %d", len(got))
	}
	if got[0].Message != "issue1" || got[1].Message != "issue2" {
		t.Errorf("unexpected failures: %+v", got)
	}
}

func TestRuleCache_PackageLevelCacheMissOnFileChange(t *testing.T) {
	configHash := [32]byte{1, 2, 3}
	sortedFiles := []string{"/tmp/a.go", "/tmp/b.go"}
	pkgKey := "/tmp"

	// First run.
	rc1 := New()
	_ = rc1.RegisterRule("rule", nil)
	rc1.CacheFileHash("/tmp/a.go", []byte("package p\nvar X = 1\n"))
	rc1.CacheFileHash("/tmp/b.go", []byte("package p\nvar Y = 2\n"))
	pkgID1, _ := rc1.PackageActionID(sortedFiles, configHash)
	rc1.PutPackage(pkgID1, pkgKey, []CachedFailure{{Message: "old"}})

	// Second run with changed file content.
	rc2 := New()
	_ = rc2.RegisterRule("rule", nil)
	rc2.CacheFileHash("/tmp/a.go", []byte("package p\nvar X = 1\n"))
	rc2.CacheFileHash("/tmp/b.go", []byte("package p\nvar Y = 999\n")) // changed
	pkgID2, _ := rc2.PackageActionID(sortedFiles, configHash)

	if pkgID1 == pkgID2 {
		t.Fatal("expected different package action IDs after file change")
	}
}

func TestRuleCache_PackageLevelCacheMissOnConfigChange(t *testing.T) {
	sortedFiles := []string{"/tmp/a.go"}
	pkgKey := "/tmp"

	rc := New()
	_ = rc.RegisterRule("rule", nil)
	rc.CacheFileHash("/tmp/a.go", []byte("package p\n"))

	configA := [32]byte{1, 2, 3}
	configB := [32]byte{4, 5, 6}

	pkgIDA, _ := rc.PackageActionID(sortedFiles, configA)
	pkgIDB, _ := rc.PackageActionID(sortedFiles, configB)

	if pkgIDA == pkgIDB {
		t.Fatal("expected different package action IDs for different configs")
	}

	rc.PutPackage(pkgIDA, pkgKey, []CachedFailure{{Message: "A"}})
	_, hit := rc.GetPackage(pkgIDB, pkgKey)
	if hit {
		t.Fatal("expected miss for different config")
	}
}

func TestRuleCache_PackageLevelDiskPersistence(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "cache")

	configHash := [32]byte{9, 8, 7}
	sortedFiles := []string{"/tmp/a.go"}
	pkgKey := "/tmp"
	content := []byte("package p\n")
	failures := []CachedFailure{{Message: "pkg-disk", RuleName: "rule"}}

	// Write via first instance.
	rc1, err := NewWithDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	_ = rc1.RegisterRule("rule", nil)
	rc1.CacheFileHash("/tmp/a.go", content)
	pkgID, _ := rc1.PackageActionID(sortedFiles, configHash)
	rc1.PutPackage(pkgID, pkgKey, failures)
	if err := rc1.FlushAll(); err != nil {
		t.Fatal(err)
	}

	// Read via second instance.
	rc2, err := NewWithDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	_ = rc2.RegisterRule("rule", nil)
	rc2.CacheFileHash("/tmp/a.go", content)
	pkgID2, _ := rc2.PackageActionID(sortedFiles, configHash)
	if pkgID != pkgID2 {
		t.Fatal("expected same package action ID for same inputs")
	}

	got, hit := rc2.GetPackage(pkgID2, pkgKey)
	if !hit {
		t.Fatal("expected disk cache hit for package-level entry")
	}
	if len(got) != 1 || got[0].Message != "pkg-disk" {
		t.Errorf("unexpected result: %+v", got)
	}
}

func TestCacheFileHash_DifferentInstancesCanDiffer(t *testing.T) {
	// Separate cache instances (simulating separate lint runs) compute
	// independent hashes, so changed file content is correctly detected.
	rc1 := New()
	rc2 := New()

	h1 := rc1.CacheFileHash("/tmp/test.go", []byte("package main\n"))
	h2 := rc2.CacheFileHash("/tmp/test.go", []byte("package main\n// changed\n"))

	if h1 == h2 {
		t.Error("expected different hash from different cache instances with different content")
	}
}
