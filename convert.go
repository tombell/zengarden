package zengarden

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
)

func isConvertable(src string) bool {
	if isMarkdown(src) {
		return true
	}

	switch filepath.Ext(src) {
	case ".html", ".xml":
		return true
	}

	return false
}

func parseFile(file string, vars context) (string, error) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	content := string(buf)
	lines := strings.Split(content, "\n")

	if len(lines) > 2 && lines[0] == "+++" {
		var n int
		var line string

		for n, line = range lines[1:] {
			if line == "+++" {
				break
			}
		}

		frontmatter := []byte(strings.Join(lines[1:n+1], "\n"))

		if err := toml.Unmarshal(frontmatter, &vars); err != nil {
			return "", fmt.Errorf("%s: %v", file, err)
		}

		content = strings.Join(lines[n+2:], "\n")
	} else if isMarkdown(file) {
		vars["title"] = ""
		vars["date"] = ""
	}

	return content, nil
}

func convertFile(src, dst, url string, siteVars context) error {
	dir := filepath.Dir(dst)

	if _, err := os.Stat(dir); err != nil {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	if !isConvertable(src) {
		sf, err := os.Open(src)
		if err != nil {
			return err
		}
		defer sf.Close()

		df, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer df.Close()

		_, err = io.Copy(df, sf)
		return err
	}

	if isMarkdown(src) {
		dst = dst[0:len(dst)-len(filepath.Ext(dst))] + ".html"
	}

	vars := context{"content": ""}

	for {
		for k, v := range siteVars {
			vars[k] = v
		}

		pageVars := context{}

		content, err := parseFile(src, pageVars)
		if err != nil {
			return err
		}

		for k, v := range pageVars {
			vars[k] = v
		}

		if content != "" {
			tmpl, err := template.New("tmpl").Parse(content)
			if err != nil {
				return err
			}

			var output bytes.Buffer
			if err := tmpl.Execute(&output, vars); err != nil {
				return err
			}

			content = output.String()
		}

		vars["page"] = context{
			"date":  toDate(src, vars),
			"url":   url,
			"title": str(vars["title"]),
		}

		if isMarkdown(src) {
			vars["content"] = renderMarkdown(content)
		} else {
			vars["content"] = content
		}

		if str(vars["layout"]) == "" || str(vars["layout"]) == "nil" {
			break
		}

		src = filepath.ToSlash(filepath.Join("_layouts", str(vars["layout"])+".html"))
		content = str(vars["content"])

		vars["content"] = content
		vars["page"].(context)["content"] = content
		vars["layout"] = ""
	}

	return ioutil.WriteFile(dst, []byte(str(vars["content"])), 0644)
}
