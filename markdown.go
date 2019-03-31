package zengarden

import (
	"io"
	"path/filepath"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/russross/blackfriday/v2"
)

func isMarkdown(src string) bool {
	switch filepath.Ext(src) {
	case ".md", ".mkd", ".markdown":
		return true
	}

	return false
}

func renderMarkdown(content, style string) string {
	r := &Renderer{
		base: blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
			Flags: blackfriday.CommonHTMLFlags,
		}),
		formatter: html.New(),
		style:     styles.Get(style),
	}

	renderer := blackfriday.WithRenderer(r)
	extensions := blackfriday.WithExtensions(blackfriday.CommonExtensions)

	return string(blackfriday.Run([]byte(content), renderer, extensions))
}

// Renderer is a custom blackfriday HTML renderer that uses the chroma library
// to highlight code in fenced code blocks.
type Renderer struct {
	base      *blackfriday.HTMLRenderer
	formatter *html.Formatter
	style     *chroma.Style
}

// RenderWithChroma renders the given text to the w io.Writer
func (r *Renderer) RenderWithChroma(w io.Writer, text []byte, data blackfriday.CodeBlockData) error {
	lexer := lexers.Fallback

	lang := string(data.Info)

	if lang != "" {
		lexer = lexers.Get(lang)
	}

	iterator, err := lexer.Tokenise(nil, string(text))
	if err != nil {
		return err
	}

	return r.formatter.Format(w, r.style, iterator)
}

// RenderNode satisfies the Renderer interface
func (r *Renderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	switch node.Type {
	case blackfriday.CodeBlock:
		if err := r.RenderWithChroma(w, node.Literal, node.CodeBlockData); err != nil {
			return r.base.RenderNode(w, node, entering)
		}

		return blackfriday.SkipChildren
	default:
		return r.base.RenderNode(w, node, entering)
	}
}

// RenderHeader satisfies the Renderer interface
func (r *Renderer) RenderHeader(w io.Writer, ast *blackfriday.Node) {
	r.base.RenderHeader(w, ast)
}

// RenderFooter satisfies the Renderer interface
func (r *Renderer) RenderFooter(w io.Writer, ast *blackfriday.Node) {
	r.base.RenderFooter(w, ast)
}
