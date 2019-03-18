package zengarden

import (
	"path/filepath"
	"time"
)

// Page is a single non-post page not in the _posts directory.
type Page struct {
	site *Site
	vars context
}

func (p *Page) toURL() string {
	from := p.vars["path"].(string)
	return joinURL(p.site.baseURL, filepath.ToSlash(from[len(p.site.source):]))
}

func (p *Page) toPath() string {
	from := p.vars["path"].(string)
	return filepath.ToSlash(filepath.Join(p.site.target, from[len(p.site.source):]))
}

// Pages is a collection of non-post pages.
type Pages []Page

func (p Pages) context() []context {
	ctx := make([]context, 0, len(p))

	for _, page := range p {
		ctx = append(ctx, page.vars)
	}

	return ctx
}

func (p Pages) convert(siteVars context) error {
	for _, page := range p {
		src := page.vars["path"].(string)

		if err := convertFile(src, page.toPath(), page.toURL(), siteVars); err != nil {
			return err
		}
	}

	return nil
}

func (p Pages) Len() int {
	return len(p)
}

func (p Pages) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p Pages) Less(i, j int) bool {
	return p[i].vars["date"].(time.Time).UnixNano() < p[j].vars["date"].(time.Time).UnixNano()
}
