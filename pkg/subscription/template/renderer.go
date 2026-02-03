package template

import (
	"bytes"
	"encoding/json"
	"net/url"
	"strings"
	"text/template"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	titleCaser = cases.Title(language.Und)

	funcMap = template.FuncMap{
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"title": func(s string) string {
			return titleCaser.String(s)
		},
		"trim": strings.TrimSpace,
		"join": strings.Join,
		"now":  time.Now,
		"toJSON": func(v any) (string, error) {
			buf, err := json.MarshalIndent(v, "", "  ")
			if err != nil {
				return "", err
			}
			return string(buf), nil
		},
		"urlquery": func(v string) string {
			return url.QueryEscape(v)
		},
	}
)

// Render 根据模板格式渲染订阅内容。
func Render(format, content string, data map[string]any) (string, error) {
	_ = format
	return renderGoTemplate(content, data)
}

func renderGoTemplate(content string, data map[string]any) (string, error) {
	tmpl, err := template.New("subscription").Funcs(funcMap).Parse(content)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
