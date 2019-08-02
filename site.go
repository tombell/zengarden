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
	cfg *Config

	vars       Context
	pages      Pages
	posts      Posts
	categories Categories
}

func (s *Site) build() error {
	if err := s.buildPages(); err != nil {
		return err
	}

	if err := s.buildPosts(); err != nil {
		return err
	}

	if err := os.RemoveAll(s.cfg.Target); err != nil {
		return err
	}

	if _, err := os.Stat(s.cfg.Target); err != nil {
		if err := os.MkdirAll(s.cfg.Target, 0755); err != nil {
			return err
		}
	}

	sort.Sort(sort.Reverse(s.pages))
	sort.Sort(sort.Reverse(s.posts))

	for _, category := range s.categories {
		sort.Sort(sort.Reverse(category))
	}

	s.vars["site"] = Context{}

	s.vars["site"].(Context)["baseurl"] = s.cfg.BaseURL
	s.vars["site"].(Context)["time"] = time.Now()

	s.vars["site"].(Context)["pages"] = s.pages.context()
	s.vars["site"].(Context)["posts"] = s.posts.context()
	s.vars["site"].(Context)["categories"] = s.categories.context()

	for k, v := range s.cfg.Vars {
		s.vars["site"].(Context)[k] = v
	}

	if err := s.posts.convert(s.vars); err != nil {
		return err
	}

	paginator := newPaginator(s)

	if err := s.pages.convert(s.vars); err != nil {
		return err
	}

	return paginator.generate()
}

func (s *Site) buildPages() error {
	err := filepath.Walk(s.cfg.Source, func(name string, info os.FileInfo, err error) error {
		if info == nil || name == s.cfg.Source {
			return err
		}

		from := filepath.ToSlash(name)
		dot := filepath.Base(name)[0]

		if info.IsDir() {
			if from == s.cfg.Target || dot == '.' || dot == '_' {
				return filepath.SkipDir
			}

			for _, exclude := range s.cfg.Excludes {
				if strings.HasSuffix(from, exclude) {
					return filepath.SkipDir
				}
			}
		} else {
			for _, exclude := range s.cfg.Excludes {
				if strings.HasSuffix(from, exclude) {
					return err
				}
			}

			if dot != '.' && dot != '_' {
				p := &Page{site: s, vars: Context{}}
				p.vars["path"] = from
				p.vars["url"] = p.toURL()
				p.vars["date"] = info.ModTime()

				s.pages = append(s.pages, p)
			}
		}

		return err
	})

	return err
}

func (s *Site) buildPosts() error {
	if _, err := os.Stat(postsDir); err != nil {
		return nil
	}

	err := filepath.Walk(postsDir, func(name string, info os.FileInfo, err error) error {
		if info == nil || name == postsDir {
			return err
		}

		if info.IsDir() {
			return err
		}

		from := filepath.ToSlash(name)

		if !isConvertable(from) {
			return err
		}

		p := &Post{s, Context{}}

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
			categorizedPosts := s.categories[cname]

			if categorizedPosts == nil {
				categorizedPosts = Posts{}
			}

			categorizedPosts = append(categorizedPosts, p)
			s.categories[cname] = categorizedPosts
		}

		s.posts = append(s.posts, p)

		return err
	})

	return err
}
