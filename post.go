package zengarden

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// Post is a single post from the _posts directory.
type Post struct {
	site *Site
	vars Context
}

func (p *Post) toURL() string {
	if v, ok := p.vars["permalink"]; ok {
		return filepath.ToSlash(filepath.Join(p.site.cfg.BaseURL, str(v)))
	}

	from := p.vars["path"].(string)

	ext := filepath.Ext(from)
	name := filepath.Base(from)
	name = name[0 : len(name)-len(ext)]

	if len(name) > 11 {
		date, err := time.Parse("2006-01-02-", name[:11])
		if err == nil {
			category := ""

			if v, ok := p.vars["category"]; ok {
				category, _ = v.(string)
			}

			title := name[11:]

			url := p.site.cfg.Permalink
			url = strings.Replace(url, ":title", title, -1)
			url = strings.Replace(url, ":category", category, -1)
			url = strings.Replace(url, ":year", fmt.Sprintf("%d", date.Year()), -1)
			url = strings.Replace(url, ":month", fmt.Sprintf("%02d", date.Month()), -1)
			url = strings.Replace(url, ":day", fmt.Sprintf("%02d", date.Day()), -1)

			return joinURL(p.site.cfg.BaseURL, url)
		}
	}

	return joinURL(p.site.cfg.BaseURL, name+".html")
}

func (p *Post) toPath() string {
	if v, ok := p.vars["permalink"]; ok {
		return filepath.ToSlash(filepath.Join(p.site.cfg.Target, str(v)))
	}

	from := p.vars["path"].(string)

	ext := filepath.Ext(from)
	name := filepath.Base(from)
	name = name[0 : len(name)-len(ext)]

	if len(name) > 11 {
		date, err := time.Parse("2006-01-02-", name[:11])
		if err == nil {
			category := ""

			if v, ok := p.vars["category"]; ok {
				category, _ = v.(string)
			}

			title := name[11:]

			url := p.site.cfg.Permalink
			url = strings.Replace(url, ":title", title, -1)
			url = strings.Replace(url, ":category", category, -1)
			url = strings.Replace(url, ":year", fmt.Sprintf("%d", date.Year()), -1)
			url = strings.Replace(url, ":month", fmt.Sprintf("%02d", date.Month()), -1)
			url = strings.Replace(url, ":day", fmt.Sprintf("%02d", date.Day()), -1)

			if p.site.cfg.Permalink[len(p.site.cfg.Permalink)-1:len(p.site.cfg.Permalink)] == "/" {
				url += "/index"
			}

			return filepath.ToSlash(filepath.Clean(filepath.Join(p.site.cfg.Target, url)))
		}
	}

	return filepath.ToSlash(filepath.Join(p.site.cfg.Target, name))
}

// Posts is a collection of posts.
type Posts []*Post

func (p Posts) context() []Context {
	ctx := make([]Context, 0, len(p))

	for _, post := range p {
		ctx = append(ctx, post.vars)
	}

	return ctx
}

func (p Posts) convert(siteVars Context) error {
	for _, post := range p {
		src := post.vars["path"].(string)

		if err := convertFile(src, post.toPath(), post.toURL(), post.site); err != nil {
			return err
		}
	}

	return nil
}

func (p Posts) Len() int {
	return len(p)
}

func (p Posts) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p Posts) Less(i, j int) bool {
	return p[i].vars["date"].(time.Time).UnixNano() < p[j].vars["date"].(time.Time).UnixNano()
}

// Categories is a map of category names to the posts of that category.
type Categories map[string]Posts

func (c Categories) context() map[string][]Context {
	ctx := map[string][]Context{}

	for category, posts := range c {
		ctxCategory := make([]Context, 0, len(posts))

		for _, p := range posts {
			ctxCategory = append(ctxCategory, p.vars)
		}

		ctx[category] = ctxCategory
	}

	return ctx
}
