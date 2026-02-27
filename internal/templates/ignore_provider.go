package templates

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/TerraFaster/scaf/internal/cache"
)

const (
	gitignoreBaseURL = "https://www.toptal.com/developers/gitignore/api"
	gitignoreCacheNS = "gitignore"
	gitignoreListKey = "_list"
)

// IgnoreProvider fetches and caches gitignore templates.
type IgnoreProvider struct {
	cache      *cache.Cache
	httpClient *http.Client
	userDir    string
}

func NewIgnoreProvider() (*IgnoreProvider, error) {
	return NewIgnoreProviderWithTTL(0)
}

func NewIgnoreProviderWithTTL(ttlHours int) (*IgnoreProvider, error) {
	c, err := cache.NewWithOptions(ttlHours)
	if err != nil {
		return nil, err
	}

	home, _ := os.UserHomeDir()
	userDir := filepath.Join(home, ".scaf", "templates", "gitignore")

	return &IgnoreProvider{
		cache:      c,
		httpClient: newHTTPClient(),
		userDir:    userDir,
	}, nil
}

// List returns all available gitignore template names.
func (p *IgnoreProvider) List() ([]string, error) {
	cached, ok, err := p.cache.Get(gitignoreCacheNS, gitignoreListKey)
	if err == nil && ok {
		return parseCSVList(cached), nil
	}

	resp, err := p.httpClient.Get(gitignoreBaseURL + "/list")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch gitignore list: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	_ = p.cache.Set(gitignoreCacheNS, gitignoreListKey, string(body))
	return parseCSVList(string(body)), nil
}

// Get fetches combined gitignore content for given templates.
func (p *IgnoreProvider) Get(names []string) (string, error) {
	if len(names) == 0 {
		return "", fmt.Errorf("no templates specified")
	}

	// Check for user local overrides and combine
	var parts []string
	var apiNames []string

	for _, name := range names {
		name = strings.ToLower(name)
		if content, err := p.getUserTemplate(name); err == nil {
			parts = append(parts, fmt.Sprintf("### %s (local) ###\n%s", name, content))
		} else {
			apiNames = append(apiNames, name)
		}
	}

	if len(apiNames) > 0 {
		// Check cache for combined key
		cacheKey := strings.Join(apiNames, ",")
		cached, ok, err := p.cache.Get(gitignoreCacheNS, cacheKey)
		if err == nil && ok {
			parts = append([]string{cached}, parts...)
		} else {
			content, err := p.fetchFromAPI(apiNames)
			if err != nil {
				return "", err
			}
			_ = p.cache.Set(gitignoreCacheNS, cacheKey, content)
			parts = append([]string{content}, parts...)
		}
	}

	return strings.Join(parts, "\n\n"), nil
}

func (p *IgnoreProvider) fetchFromAPI(names []string) (string, error) {
	joined := strings.Join(names, ",")
	url := fmt.Sprintf("%s/%s", gitignoreBaseURL, joined)

	resp, err := p.httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch gitignore templates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("gitignore.io API returned %d for templates %s", resp.StatusCode, joined)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (p *IgnoreProvider) getUserTemplate(key string) (string, error) {
	path := filepath.Join(p.userDir, key)
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

func parseCSVList(raw string) []string {
	raw = strings.TrimSpace(raw)
	var result []string
	for _, line := range strings.Split(raw, "\n") {
		for _, item := range strings.Split(line, ",") {
			item = strings.TrimSpace(item)
			if item != "" {
				result = append(result, item)
			}
		}
	}
	return result
}
