package zengarden

import (
	"fmt"
	"math"
	"path/filepath"
)

// Paginator handles paginating the posts into the correct pagination size, and
// generates the additional pagination pages.
type Paginator struct {
	site *Site
	vars Context
}

func newPaginator(site *Site) *Paginator {
	p := &Paginator{site, Context{}}

	p.site.vars["paginator"] = Context{}

	if p.site.cfg.Paginate > 0 {
		npages := int(math.Floor(float64(len(p.site.posts)) / float64(p.site.cfg.Paginate)))

		p.vars["per_page"] = p.site.cfg.Paginate
		p.vars["total_posts"] = len(p.site.posts)
		p.vars["total_pages"] = npages
		p.vars["page"] = 1

		nni := p.site.cfg.Paginate
		if nni > len(p.site.posts) {
			nni = len(p.site.posts)
		}

		p.vars["posts"] = p.site.posts.context()[:nni]

		if len(p.site.posts) > p.site.cfg.Paginate {
			p.vars["previous_page"] = false
			p.vars["previous_page_path"] = nil

			p.vars["next_page"] = 2
			p.vars["next_page_path"] = "/page2/"
		}

		p.site.vars["paginator"] = p.vars
	}

	return p
}

func (p *Paginator) context() Context {
	return p.vars
}

func (p *Paginator) generate() error {
	if p.site.cfg.Paginate <= 0 {
		return nil
	}

	index := p.site.pages.findIndex()

	if index == nil {
		return nil
	}

	from := index.vars["path"].(string)

	p.vars["previous_page"] = nil
	p.vars["previous_page_path"] = nil

	p.vars["next_page"] = nil
	p.vars["next_page_path"] = nil

	npages := p.vars["total_pages"].(int)

	for i := 2; i <= npages; i++ {
		if i-1 > 0 {
			p.vars["previous_page"] = i - 1
			p.vars["previous_page_path"] = fmt.Sprintf("/page%d/", i-1)
		} else {
			p.vars["previous_page"] = nil
			p.vars["previous_page_path"] = nil
		}

		if i+1 <= npages {
			p.vars["next_page"] = i + 1
			p.vars["next_page_path"] = fmt.Sprintf("/page%d/", i+1)
		} else {
			p.vars["next_page"] = nil
			p.vars["next_page_path"] = nil
		}

		p.vars["page"] = i

		nni := p.site.cfg.Paginate * i
		if nni > len(p.site.posts) {
			nni = len(p.site.posts)
		}

		p.vars["posts"] = p.site.posts.context()[p.site.cfg.Paginate*(i-1) : nni]

		p.site.vars["paginator"] = p.vars

		to := filepath.ToSlash(filepath.Join(p.site.cfg.Target, fmt.Sprintf("page%d", i), "index.html"))
		url := joinURL(p.site.cfg.BaseURL, filepath.ToSlash(from[len(p.site.cfg.Source):]))

		if err := convertFile(from, to, url, p.site.vars); err != nil {
			return err
		}
	}

	return nil
}
