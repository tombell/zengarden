package zengarden

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// Config is the configuration data for building the site.
type Config struct {
	Source    string
	Target    string
	BaseURL   string
	Permalink string
	Excludes  []string
}

// Run will build the site using the given configuration.
func Run(cfg *Config) error {
	source, _ := filepath.Abs(cfg.Source)
	target, _ := filepath.Abs(cfg.Target)

	s := Site{
		source:    filepath.ToSlash(source),
		target:    filepath.ToSlash(target),
		baseURL:   cfg.BaseURL,
		permalink: cfg.Permalink,
		excludes:  cfg.Excludes,

		vars: context{
			"site": context{},
		},
	}

	return s.build()
}

func joinURL(l, r string) string {
	r = path.Clean(r)
	ls := strings.HasSuffix(l, "/")
	rp := strings.HasPrefix(r, "/")

	if ls && rp {
		return l + r[1:]
	}

	if !ls && !rp {
		return l + "/" + r
	}

	return l + r
}

func toDate(from string, vars context) time.Time {
	if v, ok := vars["date"]; ok {
		date, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", str(v))
		if err == nil {
			return date
		}
	}

	info, err := os.Stat(from)
	if err != nil {
		return time.Now()
	}

	name := filepath.Base(from)
	if len(name) <= 11 {
		return info.ModTime()
	}

	date, err := time.Parse("2006-01-02-", name[:11])
	if err != nil {
		return info.ModTime()
	}

	return date
}
