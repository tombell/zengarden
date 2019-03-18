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
	vars context
}

func (p *Post) toURL() string {
	if v, ok := p.vars["permalink"]; ok {
		return filepath.ToSlash(filepath.Join(p.site.baseURL, str(v)))
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

			url := p.site.permalink
			url = strings.Replace(url, ":title", title, -1)
			url = strings.Replace(url, ":category", category, -1)
			url = strings.Replace(url, ":year", fmt.Sprintf("%d", date.Year()), -1)
			url = strings.Replace(url, ":month", fmt.Sprintf("%02d", date.Month()), -1)
			url = strings.Replace(url, ":day", fmt.Sprintf("%02d", date.Day()), -1)

			return joinURL(p.site.baseURL, url)
		}
	}

	return joinURL(p.site.baseURL, name+".html")
}

func (p *Post) toPath() string {
	if v, ok := p.vars["permalink"]; ok {
		return filepath.ToSlash(filepath.Join(p.site.target, str(v)))
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

			url := p.site.permalink
			url = strings.Replace(url, ":title", title, -1)
			url = strings.Replace(url, ":category", category, -1)
			url = strings.Replace(url, ":year", fmt.Sprintf("%d", date.Year()), -1)
			url = strings.Replace(url, ":month", fmt.Sprintf("%02d", date.Month()), -1)
			url = strings.Replace(url, ":day", fmt.Sprintf("%02d", date.Day()), -1)

			if p.site.permalink[len(p.site.permalink)-1:len(p.site.permalink)] == "/" {
				url += "/index"
			}

			return filepath.ToSlash(filepath.Clean(filepath.Join(p.site.target, url)))
		}
	}

	return filepath.ToSlash(filepath.Join(p.site.target, name))
}

// Posts is a collection of posts.
type Posts []Post

func (p Posts) context() []context {
	ctx := make([]context, 0, len(p))

	for _, post := range p {
		ctx = append(ctx, post.vars)
	}

	return ctx
}

func (p Posts) convert(siteVars context) error {
	for _, post := range p {
		src := post.vars["path"].(string)

		if err := convertFile(src, post.toPath(), post.toURL(), siteVars); err != nil {
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

func (c Categories) context() map[string][]context {
	ctx := map[string][]context{}

	for category, posts := range c {
		ctxCategory := make([]context, 0, len(posts))

		for _, p := range posts {
			ctxCategory = append(ctxCategory, p.vars)
		}

		ctx[category] = ctxCategory
	}

	return ctx
}
