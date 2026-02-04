package template

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
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
		"split": func(s, sep string) []string {
			return strings.Split(s, sep)
		},
		"replace": func(s, old, new string) string {
			return strings.ReplaceAll(s, old, new)
		},
		"trimPrefix": strings.TrimPrefix,
		"trimSuffix": strings.TrimSuffix,
		"b64enc": func(v any) (string, error) {
			switch value := v.(type) {
			case string:
				return base64.StdEncoding.EncodeToString([]byte(value)), nil
			case []byte:
				return base64.StdEncoding.EncodeToString(value), nil
			default:
				return "", fmt.Errorf("b64enc: unsupported type %T", v)
			}
		},
		"b64dec": func(v string) (string, error) {
			buf, err := base64.StdEncoding.DecodeString(v)
			if err != nil {
				return "", err
			}
			return string(buf), nil
		},
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
