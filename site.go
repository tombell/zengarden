package zengarden

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Site is the site being built.
type Site struct {
	source    string
	target    string
	baseURL   string
	permalink string

	excludes []string

	vars context
}

func (s *Site) build() error {
	var err error

	pages := Pages{}

	err = filepath.Walk(s.source, func(name string, info os.FileInfo, err error) error {
		if info == nil || name == s.source {
			return err
		}

		from := filepath.ToSlash(name)
		dot := filepath.Base(name)[0]

		if info.IsDir() {
			if from == s.target || dot == '.' || dot == '_' {
				return filepath.SkipDir
			}
		} else {
			for _, exclude := range s.excludes {
				if strings.HasSuffix(from, exclude) {
					return err
				}
			}

			if dot != '.' && dot != '_' {
				p := Page{site: s, vars: context{}}
				p.vars["path"] = from
				p.vars["url"] = p.toURL()
				p.vars["date"] = info.ModTime()

				pages = append(pages, p)
			}
		}

		return err
	})

	if err != nil {
		return err
	}

	categories := Categories{}
	posts := Posts{}

	err = filepath.Walk("_posts", func(name string, info os.FileInfo, err error) error {
		if info == nil || name == "_posts" {
			return err
		}

		if info.IsDir() {
			return err
		}

		from := filepath.ToSlash(name)

		if !isConvertable(from) {
			return err
		}

		p := Post{
			site: s,
			vars: context{},
		}

		content, err := parseFile(from, p.vars)
		if err != nil {
			return err
		}

		if _, err := os.Stat(from); err != nil {
			return err
		}

		p.vars["path"] = from
		p.vars["url"] = p.toURL()
		p.vars["date"] = toDate(from, p.vars)
		p.vars["content"] = content

		if category, ok := p.vars["category"]; ok {
			cname := str(category)
			categorizedPosts := categories[cname]

			if categorizedPosts == nil {
				categorizedPosts = Posts{}
			}

			categorizedPosts = append(categorizedPosts, p)
			categories[cname] = categorizedPosts
		}

		posts = append(posts, p)

		return err
	})

	if err != nil {
		return err
	}

	if _, err := os.Stat(s.target); err != nil {
		if err := os.MkdirAll(s.target, 0755); err != nil {
			return err
		}
	}

	sort.Sort(sort.Reverse(pages))
	sort.Sort(sort.Reverse(posts))

	for _, category := range categories {
		sort.Sort(sort.Reverse(category))
	}

	s.vars["site"].(context)["url"] = s.baseURL
	s.vars["site"].(context)["baseurl"] = s.baseURL
	s.vars["site"].(context)["time"] = time.Now()

	s.vars["site"].(context)["pages"] = pages.context()
	s.vars["site"].(context)["posts"] = posts.context()
	s.vars["site"].(context)["categories"] = categories.context()

	if err := pages.convert(s.vars); err != nil {
		return err
	}

	if err := posts.convert(s.vars); err != nil {
		return err
	}

	return nil
}
