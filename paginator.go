package zengarden

import (
	"fmt"
	"math"
	"path/filepath"
)

// Paginator handles paginating the posts into the correct pagination size, and
// generates the additional pagination pages.
type Paginator struct {
	site  *Site
	pages Pages
	posts Posts
	vars  context
}

func newPaginator(site *Site, pages Pages, posts Posts) *Paginator {
	p := &Paginator{site, pages, posts, context{}}

	p.site.vars["paginator"] = context{}

	if p.site.paginate > 0 {
		npages := int(math.Floor(float64(len(posts)) / float64(p.site.paginate)))

		p.vars["per_page"] = p.site.paginate
		p.vars["total_posts"] = len(p.posts)
		p.vars["total_pages"] = npages
		p.vars["page"] = 1

		nni := p.site.paginate
		if nni > len(p.posts) {
			nni = len(p.posts)
		}

		p.vars["posts"] = p.posts.context()[:nni]

		if len(posts) > p.site.paginate {
			p.vars["previous_page"] = false
			p.vars["previous_page_path"] = nil

			p.vars["next_page"] = 2
			p.vars["next_page_path"] = "/page2/"
		}

		p.site.vars["paginator"] = p.vars
	}

	return p
}

func (p *Paginator) context() context {
	return p.vars
}

func (p *Paginator) generate() error {
	if p.site.paginate <= 0 {
		return nil
	}

	index := p.pages.findIndex()

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

		nni := p.site.paginate * i
		if nni > len(p.posts) {
			nni = len(p.posts)
		}

		p.vars["posts"] = p.posts.context()[p.site.paginate*(i-1) : nni]

		p.site.vars["paginator"] = p.vars

		to := filepath.ToSlash(filepath.Join(p.site.target, fmt.Sprintf("page%d", i), "index.html"))
		url := joinURL(p.site.baseURL, filepath.ToSlash(from[len(p.site.source):]))

		if err := convertFile(from, to, url, p.site.vars); err != nil {
			return err
		}
	}

	return nil
}
