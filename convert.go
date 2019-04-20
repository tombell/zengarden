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
	"time"

	"github.com/BurntSushi/toml"
)

var funcMap = template.FuncMap{
	"dateToString": dateToString,
	"filter":       filter,
	"slice":        slice,
}

func dateToString(date time.Time) string {
	return date.Format("2 Jan 2006")
}

func filter(data []Context, key string, val interface{}) []Context {
	var result []Context

	for _, ctx := range data {
		if ctx[key] == val {
			result = append(result, ctx)
		}
	}

	return nil
}

func slice(data []Context, offset, count int) []Context {
	return data[offset:count]
}

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

func copyFile(src, dst string) error {
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

func parseFile(file string, vars Context) (string, error) {
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

func renderTemplate(src, content string, vars Context) (string, error) {
	tmpl, err := template.New(src).Funcs(funcMap).Parse(content)
	if err != nil {
		return "", err
	}

	if _, err := tmpl.ParseGlob(filepath.Join(includesDir, "*.html")); err != nil {
		return "", err
	}

	var output bytes.Buffer
	if err := tmpl.Execute(&output, vars); err != nil {
		return "", err
	}

	return output.String(), nil
}

func convertFile(src, dst, url string, site *Site, postVars Context) error {
	dir := filepath.Dir(dst)

	if _, err := os.Stat(dir); err != nil {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	if !isConvertable(src) {
		return copyFile(src, dst)
	}

	if isMarkdown(src) {
		dst = dst[0:len(dst)-len(filepath.Ext(dst))] + ".html"
	}

	vars := Context{"content": ""}

	for {
		for k, v := range site.vars {
			vars[k] = v
		}

		pageVars := Context{}

		content, err := parseFile(src, pageVars)
		if err != nil {
			return err
		}

		for k, v := range pageVars {
			vars[k] = v
		}

		if content != "" {
			output, err := renderTemplate(src, content, vars)
			if err != nil {
				return err
			}

			content = output
		}

		vars["page"] = Context{
			"date":     toDate(src, vars),
			"url":      url,
			"title":    str(vars["title"]),
			"next":     postVars["next"],
			"previous": postVars["previous"],
		}

		if isMarkdown(src) {
			vars["content"] = renderMarkdown(content, site.cfg.Style)
		} else {
			vars["content"] = content
		}

		if str(vars["layout"]) == "" || str(vars["layout"]) == "nil" {
			break
		}

		src = filepath.ToSlash(filepath.Join(layoutsDir, str(vars["layout"])+".html"))
		content = str(vars["content"])

		vars["content"] = content
		vars["page"].(Context)["content"] = content
		vars["layout"] = ""
	}

	return ioutil.WriteFile(dst, []byte(str(vars["content"])), 0644)
}
