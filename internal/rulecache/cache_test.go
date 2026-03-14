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
	_, hit := rc.Get("test-rule", "/tmp/test.go", TierFileOnly, nil)
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
	rc.Put("test-rule", "/tmp/test.go", TierFileOnly, nil, failures)

	got, hit := rc.Get("test-rule", "/tmp/test.go", TierFileOnly, nil)
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
	rc.Put("rule", "/tmp/test.go", TierFileOnly, nil, []CachedFailure{{Message: "A"}})

	got, hit := rc.Get("rule", "/tmp/test.go", TierFileOnly, nil)
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

	_, hit = rc.Get("rule", "/tmp/test.go", TierFileOnly, nil)
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
	rc.Put("rule", "/tmp/a.go", TierPackageAware, siblings, []CachedFailure{{Message: "pkg"}})

	got, hit := rc.Get("rule", "/tmp/a.go", TierPackageAware, siblings)
	if !hit {
		t.Fatal("expected hit")
	}
	if got[0].Message != "pkg" {
		t.Errorf("expected pkg, got %s", got[0].Message)
	}

	// Change sibling content — should invalidate.
	content2Changed := []byte("package p\nvar Y = 3\n")
	rc.CacheFileHash("/tmp/b.go", content2Changed) // overwrite cached hash
	siblings["/tmp/b.go"] = content2Changed

	_, hit = rc.Get("rule", "/tmp/a.go", TierPackageAware, siblings)
	if hit {
		t.Fatal("expected miss after sibling change")
	}
}

func TestRuleCache_UnregisteredRule(t *testing.T) {
	rc := New()
	rc.CacheFileHash("/tmp/test.go", []byte("package main\n"))

	_, hit := rc.Get("unregistered", "/tmp/test.go", TierFileOnly, nil)
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
	rc.Put("clean-rule", "/tmp/clean.go", TierFileOnly, nil, []CachedFailure{})

	got, hit := rc.Get("clean-rule", "/tmp/clean.go", TierFileOnly, nil)
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
	rc1.Put("rule", "/tmp/test.go", TierFileOnly, nil, failures)
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

	got, hit := rc2.Get("rule", "/tmp/test.go", TierFileOnly, nil)
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
	rc1.Put("rule", "/tmp/test.go", TierFileOnly, nil, failures)
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

	_, hit := rc2.Get("rule", "/tmp/test.go", TierFileOnly, nil)
	if hit {
		t.Fatal("expected miss after file content change")
	}
}

func TestCacheFileHash_UpdatesOnNewContent(t *testing.T) {
	rc := New()

	content := []byte("package main\n")
	h1 := rc.CacheFileHash("/tmp/test.go", content)
	h2 := rc.CacheFileHash("/tmp/test.go", []byte("different content"))

	// Second call with different content should produce different hash.
	if h1 == h2 {
		t.Error("expected different hash for different content")
	}

	// Same content should produce same hash.
	h3 := rc.CacheFileHash("/tmp/test.go", content)
	if h1 != h3 {
		t.Error("expected same hash for same content")
	}
}
