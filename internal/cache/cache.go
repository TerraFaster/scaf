package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const defaultTTLHours = 168 // 7 days

// Cache manages file-based caching in ~/.scaf/cache/
type Cache struct {
	baseDir    string
	ttlHours   int
}

// New creates a Cache using the default directory (~/.scaf/cache).
func New() (*Cache, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("cannot determine home directory: %w", err)
	}
	baseDir := filepath.Join(home, ".scaf", "cache")
	return &Cache{
		baseDir:  baseDir,
		ttlHours: defaultTTLHours,
	}, nil
}

// NewWithOptions creates a Cache with custom TTL.
func NewWithOptions(ttlHours int) (*Cache, error) {
	c, err := New()
	if err != nil {
		return nil, err
	}
	if ttlHours > 0 {
		c.ttlHours = ttlHours
	}
	return c, nil
}

func (c *Cache) dir(namespace string) string {
	return filepath.Join(c.baseDir, namespace)
}

// Get retrieves cached content. Returns ("", false, nil) if not found or expired.
func (c *Cache) Get(namespace, key string) (string, bool, error) {
	path := c.keyPath(namespace, key)
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}

	// Check TTL
	age := time.Since(info.ModTime())
	if age > time.Duration(c.ttlHours)*time.Hour {
		return "", false, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", false, err
	}
	return string(data), true, nil
}

// Set stores content in cache.
func (c *Cache) Set(namespace, key, content string) error {
	dir := c.dir(namespace)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("cannot create cache dir: %w", err)
	}

	path := c.keyPath(namespace, key)
	// Validate path (prevent path traversal)
	if !isSubPath(c.baseDir, path) {
		return fmt.Errorf("invalid cache key: path traversal detected")
	}

	return os.WriteFile(path, []byte(content), 0644)
}

// Clear removes all cached files.
func (c *Cache) Clear() error {
	return os.RemoveAll(c.baseDir)
}

func (c *Cache) keyPath(namespace, key string) string {
	// Sanitize key: replace special chars with underscore
	safe := sanitizeKey(key)
	return filepath.Join(c.dir(namespace), safe)
}

func sanitizeKey(key string) string {
	result := make([]byte, len(key))
	for i := 0; i < len(key); i++ {
		ch := key[i]
		if isAlphanumeric(ch) || ch == '-' || ch == '_' || ch == '.' {
			result[i] = ch
		} else {
			result[i] = '_'
		}
	}
	return string(result)
}

func isAlphanumeric(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9')
}

// isSubPath checks path is under base (path traversal protection).
func isSubPath(base, path string) bool {
	base = filepath.Clean(base)
	path = filepath.Clean(path)
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return false
	}
	return rel != ".." && len(rel) > 1 || (len(rel) == 1 && rel != "..")
}
