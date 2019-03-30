package zengarden

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const (
	includesDir = "_includes"
	layoutsDir  = "_layouts"
	postsDir    = "_posts"
)

// Run will build the site using the given configuration.
func Run(cfg *Config) error {
	s := Site{
		cfg:  cfg,
		vars: Context{"site": Context{}},
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

func toDate(from string, vars Context) time.Time {
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
