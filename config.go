package zengarden

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const (
	// DefaultConfigPath is the default path for the zen-garden config file.
	DefaultConfigPath = "_garden.toml"
)

// Config is the configuration data for building the site.
type Config struct {
	Source    string
	Target    string
	BaseURL   string
	Permalink string
	Paginate  int
	Style     string
	Excludes  []string

	Vars Context
}

// LoadConfig loads the configuration file at the given path, and adds any
// non-config variables to the Vars context.
func LoadConfig(path string) (*Config, error) {
	cfg := &Config{
		Permalink: "/posts/:title",
		Style:     "fallback",
		Vars:      Context{},
	}

	cfg.Source, _ = filepath.Abs(".")
	cfg.Target, _ = filepath.Abs("_site")

	if _, err := os.Stat(path); err == nil {
		vars := Context{}

		if _, err := toml.DecodeFile(path, &vars); err != nil {
			return nil, err
		}

		if src, ok := vars["source"].(string); ok {
			cfg.Source, _ = filepath.Abs(src)
			delete(vars, "source")
		}

		if dst, ok := vars["target"].(string); ok {
			cfg.Target, _ = filepath.Abs(dst)
			delete(vars, "target")
		}

		if baseURL, ok := vars["baseurl"].(string); ok {
			cfg.BaseURL = baseURL
			delete(vars, "baseurl")
		}

		if permalink, ok := vars["permalink"].(string); ok {
			cfg.Permalink = permalink
			delete(vars, "permalink")
		}

		if paginate, ok := vars["paginate"].(int64); ok {
			cfg.Paginate = int(paginate)
			delete(vars, "paginate")
		}

		if style, ok := vars["syntax_highlight"].(string); ok {
			cfg.Style = style
			delete(vars, "syntax_highlight")
		}

		if excludes, ok := vars["excludes"].([]interface{}); ok {
			for _, exclude := range excludes {
				cfg.Excludes = append(cfg.Excludes, exclude.(string))
			}
			delete(vars, "excludes")
		}

		cfg.Vars = vars
	}

	return cfg, nil
}
