package rulecache

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

// RuleCache caches results of individual lint rule executions.
// Each cache entry maps a (rule, file) pair to its []CachedFailure result.
//
// On disk, entries are grouped per package: one file per package containing
// all (rule, file) action IDs for that package. This reduces disk I/O from
// tens of thousands of tiny files to one file per package.
type RuleCache struct {
	ruleConfigs sync.Map // ruleName → [32]byte (config hash)
	fileHashes  sync.Map // filePath → [32]byte (content hash)
	results     sync.Map // [32]byte (action ID) → []CachedFailure (hot in-memory layer)
	dir         string   // disk cache directory; empty means in-memory only

	// pendingMu protects pendingPkg.
	pendingMu  sync.Mutex
	pendingPkg map[string]map[[32]byte]bool // pkgKey → set of action IDs to flush

	// loadedPkgs tracks which package cache files have been loaded into memory.
	loadedPkgs sync.Map // pkgKey (string) → bool
}

// New creates a new in-memory-only RuleCache.
func New() *RuleCache {
	return &RuleCache{
		pendingPkg: map[string]map[[32]byte]bool{},
	}
}

// NewWithDir creates a RuleCache backed by the given directory for disk persistence.
// The directory is created if it does not exist.
func NewWithDir(dir string) (*RuleCache, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("rulecache: create dir: %w", err)
	}
	return &RuleCache{
		dir:        dir,
		pendingPkg: map[string]map[[32]byte]bool{},
	}, nil
}

// RegisterRule pre-computes and caches the config hash for a rule.
// Must be called once per rule before any Get/Put.
func (rc *RuleCache) RegisterRule(ruleName string, args []any) error {
	h := sha256.New()
	fmt.Fprintf(h, "rule %s\n", ruleName)

	if len(args) > 0 {
		var buf bytes.Buffer
		if err := gob.NewEncoder(&buf).Encode(args); err != nil {
			return fmt.Errorf("failed to hash rule config for %q: %w", ruleName, err)
		}
		h.Write(buf.Bytes())
	}

	var hash [32]byte
	copy(hash[:], h.Sum(nil))
	rc.ruleConfigs.Store(ruleName, hash)

	return nil
}

// CacheFileHash returns the cached content hash for filePath, computing and
// storing it only when no entry exists yet. Within a single lint run the
// content for a given path never changes, so the first call wins and all
// subsequent calls (e.g. from actionID for sibling files) skip SHA-256.
func (rc *RuleCache) CacheFileHash(filePath string, content []byte) [32]byte {
	if existing, ok := rc.fileHashes.Load(filePath); ok {
		return existing.([32]byte)
	}
	hash := sha256.Sum256(content)
	rc.fileHashes.Store(filePath, hash)
	return hash
}

// Get retrieves cached failures for a (rule, file) pair.
// Returns (failures, true) on hit, (nil, false) on miss.
//
// siblingDigest is an optional pre-computed digest of the sibling files
// (from SiblingDigest). When non-nil and tier >= TierPackageAware, the
// digest is used directly instead of re-sorting and re-hashing siblings.
// Pass nil to have actionID compute it on the fly (original behaviour).
func (rc *RuleCache) Get(
	ruleName string,
	filePath string,
	tier CacheTier,
	siblingFiles map[string][]byte, // filePath → content, for TierPackageAware+
	siblingDigest *[32]byte, // optional pre-computed sibling digest
) ([]CachedFailure, bool) {
	actionID, ok := rc.actionID(ruleName, filePath, tier, siblingFiles, siblingDigest)
	if !ok {
		return nil, false
	}

	// Hot in-memory check.
	if v, ok := rc.results.Load(actionID); ok {
		return v.([]CachedFailure), true
	}

	// Try loading the package cache file from disk (once per package).
	if rc.dir != "" {
		pkgKey := packageKey(filePath)
		if _, loaded := rc.loadedPkgs.Load(pkgKey); !loaded {
			rc.loadPkgFromDisk(pkgKey)
			// Re-check memory after loading.
			if v, ok := rc.results.Load(actionID); ok {
				return v.([]CachedFailure), true
			}
		}
	}

	return nil, false
}

// Put stores failures for a (rule, file) pair.
// Writes to memory immediately; call FlushPackage to persist to disk.
//
// siblingDigest is an optional pre-computed digest (see Get).
func (rc *RuleCache) Put(
	ruleName string,
	filePath string,
	tier CacheTier,
	siblingFiles map[string][]byte,
	siblingDigest *[32]byte,
	failures []CachedFailure,
) {
	actionID, ok := rc.actionID(ruleName, filePath, tier, siblingFiles, siblingDigest)
	if !ok {
		return
	}

	rc.results.Store(actionID, failures)

	// Track this action ID for later batch flush.
	if rc.dir != "" {
		pkgKey := packageKey(filePath)
		rc.pendingMu.Lock()
		if rc.pendingPkg[pkgKey] == nil {
			rc.pendingPkg[pkgKey] = map[[32]byte]bool{}
		}
		rc.pendingPkg[pkgKey][actionID] = true
		rc.pendingMu.Unlock()
	}
}

// FlushAll persists all pending cache entries to disk, grouped by package.
// Should be called after linting completes.
func (rc *RuleCache) FlushAll() error {
	if rc.dir == "" {
		return nil
	}

	rc.pendingMu.Lock()
	pending := rc.pendingPkg
	rc.pendingPkg = map[string]map[[32]byte]bool{}
	rc.pendingMu.Unlock()

	for pkgKey, actionIDs := range pending {
		if err := rc.flushPkg(pkgKey, actionIDs); err != nil {
			return err
		}
	}
	return nil
}

// flushPkg writes all entries for a package to a single file on disk.
// It merges with any existing entries already on disk.
func (rc *RuleCache) flushPkg(pkgKey string, newActionIDs map[[32]byte]bool) error {
	p := rc.pkgDiskPath(pkgKey)

	// Load existing entries from disk to merge.
	existing := map[[32]byte][]CachedFailure{}
	if data, err := os.ReadFile(p); err == nil {
		_ = gob.NewDecoder(bytes.NewReader(data)).Decode(&existing)
	}

	// Merge new entries.
	for id := range newActionIDs {
		if v, ok := rc.results.Load(id); ok {
			existing[id] = v.([]CachedFailure)
		}
	}

	// Encode and write.
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(existing); err != nil {
		return err
	}

	return os.WriteFile(p, buf.Bytes(), 0o644)
}

// loadPkgFromDisk loads all entries from a package's cache file into memory.
func (rc *RuleCache) loadPkgFromDisk(pkgKey string) {
	// Mark as loaded regardless of success to avoid repeated attempts.
	if _, alreadyLoaded := rc.loadedPkgs.LoadOrStore(pkgKey, true); alreadyLoaded {
		return
	}

	p := rc.pkgDiskPath(pkgKey)
	data, err := os.ReadFile(p)
	if err != nil {
		return
	}

	var entries map[[32]byte][]CachedFailure
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&entries); err != nil {
		return
	}

	for id, failures := range entries {
		rc.results.Store(id, failures)
	}
}

// PackageActionID computes a single cache key for an entire package.
// sortedFilenames must already be sorted. configHash is an opaque hash of all
// configuration that affects lint output (rule configs, confidence, etc.).
// File content hashes are looked up from the internal fileHashes map (populated
// by prior CacheFileHash calls).
func (rc *RuleCache) PackageActionID(sortedFilenames []string, configHash [32]byte) ([32]byte, bool) {
	h := sha256.New()
	fmt.Fprintf(h, "pkg\n")
	h.Write(configHash[:])
	for _, f := range sortedFilenames {
		fh, ok := rc.fileHashes.Load(f)
		if !ok {
			return [32]byte{}, false
		}
		fmt.Fprintf(h, "file %s %x\n", f, fh.([32]byte))
	}
	var id [32]byte
	copy(id[:], h.Sum(nil))
	return id, true
}

// GetPackage retrieves cached failures for an entire package.
// pkgKey is the package directory (used to load from disk on first access).
func (rc *RuleCache) GetPackage(actionID [32]byte, pkgKey string) ([]CachedFailure, bool) {
	if v, ok := rc.results.Load(actionID); ok {
		return v.([]CachedFailure), true
	}
	if rc.dir != "" {
		if _, loaded := rc.loadedPkgs.Load(pkgKey); !loaded {
			rc.loadPkgFromDisk(pkgKey)
			if v, ok := rc.results.Load(actionID); ok {
				return v.([]CachedFailure), true
			}
		}
	}
	return nil, false
}

// PutPackage stores failures for an entire package.
// pkgKey is the package directory (used for disk persistence grouping).
func (rc *RuleCache) PutPackage(actionID [32]byte, pkgKey string, failures []CachedFailure) {
	rc.results.Store(actionID, failures)
	if rc.dir != "" {
		rc.pendingMu.Lock()
		if rc.pendingPkg[pkgKey] == nil {
			rc.pendingPkg[pkgKey] = map[[32]byte]bool{}
		}
		rc.pendingPkg[pkgKey][actionID] = true
		rc.pendingMu.Unlock()
	}
}

// SiblingDigest computes a single SHA-256 hash that represents all sibling
// files (every file in siblingFiles except filePath itself). The result is
// deterministic: sibling paths are sorted before hashing.
//
// Call this once per (filePath, siblingFiles) pair, then pass the result as
// siblingDigest to Get/Put to avoid re-sorting and re-hashing for every rule.
func (rc *RuleCache) SiblingDigest(filePath string, siblingFiles map[string][]byte) [32]byte {
	h := sha256.New()
	paths := make([]string, 0, len(siblingFiles))
	for p := range siblingFiles {
		if p == filePath {
			continue
		}
		paths = append(paths, p)
	}
	sort.Strings(paths)
	for _, p := range paths {
		sibHash := rc.CacheFileHash(p, siblingFiles[p])
		fmt.Fprintf(h, "sibling %s %x\n", p, sibHash)
	}
	var digest [32]byte
	copy(digest[:], h.Sum(nil))
	return digest
}

// actionID computes the unique cache key for a (rule, file) pair.
//
// When siblingDigest is non-nil and tier >= TierPackageAware, the provided
// digest is folded in directly, skipping the per-sibling sort+hash work.
func (rc *RuleCache) actionID(
	ruleName string,
	filePath string,
	tier CacheTier,
	siblingFiles map[string][]byte,
	siblingDigest *[32]byte,
) ([32]byte, bool) {
	confHash, ok := rc.ruleConfigs.Load(ruleName)
	if !ok {
		return [32]byte{}, false
	}
	ch := confHash.([32]byte)

	fileHash, ok := rc.fileHashes.Load(filePath)
	if !ok {
		return [32]byte{}, false
	}
	fh := fileHash.([32]byte)

	h := sha256.New()
	fmt.Fprintf(h, "rule %s %x\n", ruleName, ch)
	fmt.Fprintf(h, "file %s %x\n", filePath, fh)

	// For package-aware and cross-package tiers, include sibling file hashes.
	if tier >= TierPackageAware {
		if siblingDigest != nil {
			// Fast path: use pre-computed digest.
			fmt.Fprintf(h, "siblings %x\n", *siblingDigest)
		} else if siblingFiles != nil {
			// Slow path: sort paths and hash each sibling individually.
			paths := make([]string, 0, len(siblingFiles))
			for p := range siblingFiles {
				if p == filePath {
					continue // already included above
				}
				paths = append(paths, p)
			}
			sort.Strings(paths)

			for _, p := range paths {
				sibHash := rc.CacheFileHash(p, siblingFiles[p])
				fmt.Fprintf(h, "sibling %s %x\n", p, sibHash)
			}
		}
	}

	var id [32]byte
	copy(id[:], h.Sum(nil))

	return id, true
}

// packageKey returns a stable key for the package containing filePath.
// Uses the directory of the file path.
func packageKey(filePath string) string {
	return filepath.Dir(filePath)
}

// pkgDiskPath returns the on-disk path for a package cache file,
// using the first byte of the hash as a subdirectory (like Go's cache: 00/ .. ff/).
func (rc *RuleCache) pkgDiskPath(pkgKey string) string {
	h := sha256.Sum256([]byte(pkgKey))
	hex := fmt.Sprintf("%x", h)
	return filepath.Join(rc.dir, hex[:2], hex[2:])
}
