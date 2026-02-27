package templates

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/TerraFaster/scaf/internal/cache"
)

const (
	githubLicensesURL = "https://api.github.com/licenses"
	licensesCacheNS   = "licenses"
	licensesListKey   = "_list"
)

// LicenseInfo holds basic info about a license.
type LicenseInfo struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	SpdxID string `json:"spdx_id"`
	URL    string `json:"url"`
}

// LicenseDetail holds full license details.
type LicenseDetail struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	Body string `json:"body"`
}

// LicenseProvider fetches and caches licenses.
type LicenseProvider struct {
	cache      *cache.Cache
	httpClient *http.Client
	userDir    string
}

func NewLicenseProvider() (*LicenseProvider, error) {
	return NewLicenseProviderWithTTL(0)
}

func NewLicenseProviderWithTTL(ttlHours int) (*LicenseProvider, error) {
	c, err := cache.NewWithOptions(ttlHours)
	if err != nil {
		return nil, err
	}

	home, _ := os.UserHomeDir()
	userDir := filepath.Join(home, ".scaf", "templates", "licenses")

	return &LicenseProvider{
		cache:      c,
		httpClient: newHTTPClient(),
		userDir:    userDir,
	}, nil
}

// List returns all available licenses.
func (p *LicenseProvider) List() ([]LicenseInfo, error) {
	// Check cache
	cached, ok, err := p.cache.Get(licensesCacheNS, licensesListKey)
	if err == nil && ok {
		var licenses []LicenseInfo
		if json.Unmarshal([]byte(cached), &licenses) == nil {
			return licenses, nil
		}
	}

	// Fetch from GitHub
	resp, err := p.httpClient.Get(githubLicensesURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch licenses list: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var licenses []LicenseInfo
	if err := json.Unmarshal(body, &licenses); err != nil {
		return nil, fmt.Errorf("failed to parse licenses: %w", err)
	}

	// Cache it
	_ = p.cache.Set(licensesCacheNS, licensesListKey, string(body))

	return licenses, nil
}

// Get fetches the body of a specific license by key.
func (p *LicenseProvider) Get(key string) (string, error) {
	key = strings.ToLower(key)

	// Check user local templates first
	if body, err := p.getUserTemplate(key); err == nil {
		return body, nil
	}

	// Check cache
	cached, ok, err := p.cache.Get(licensesCacheNS, key)
	if err == nil && ok {
		// cached is the full body text
		return cached, nil
	}

	// Fetch from GitHub
	url := fmt.Sprintf("%s/%s", githubLicensesURL, key)
	resp, err := p.httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch license %s: %w", key, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("license %q not found", key)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned %d for license %s", resp.StatusCode, key)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var detail LicenseDetail
	if err := json.Unmarshal(data, &detail); err != nil {
		return "", fmt.Errorf("failed to parse license detail: %w", err)
	}

	if detail.Body == "" {
		return "", fmt.Errorf("empty license body for %s", key)
	}

	// Cache just the body
	_ = p.cache.Set(licensesCacheNS, key, detail.Body)

	return detail.Body, nil
}

// FuzzySearch returns licenses with keys similar to query.
func (p *LicenseProvider) FuzzySearch(query string, licenses []LicenseInfo) []LicenseInfo {
	query = strings.ToLower(query)
	keys := make([]string, len(licenses))
	for i, l := range licenses {
		keys[i] = l.Key
	}

	matches := fuzzy.RankFindNormalizedFold(query, keys)
	result := make([]LicenseInfo, 0, len(matches))
	for _, m := range matches {
		for _, l := range licenses {
			if l.Key == m.Target {
				result = append(result, l)
				break
			}
		}
		if len(result) >= 5 {
			break
		}
	}
	return result
}

func (p *LicenseProvider) getUserTemplate(key string) (string, error) {
	path := filepath.Join(p.userDir, key)
	// path traversal check
	clean := filepath.Clean(path)
	if !strings.HasPrefix(clean, filepath.Clean(p.userDir)) {
		return "", fmt.Errorf("invalid key")
	}
	data, err := os.ReadFile(clean)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
