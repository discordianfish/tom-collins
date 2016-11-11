package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/tumblr/go-collins.v0/collins"
)

func tmplFindByTag(tag string) *collins.Asset {
	asset, _, err := client.Assets.Get(tag)
	if err != nil {
		log.Fatal(err)
	}
	return asset
}

func tmplQuery(q string) (assets []collins.Asset) {
	it, err := NewFindIterator(&collins.AssetFindOpts{Query: q})
	if err != nil {
		log.Fatal(err)
	}
	for it.Next() {
		assets = append(assets, it.Value())
	}
	if err := it.Err(); err != nil {
		log.Fatal(err)
	}
	return assets
}

type templater struct {
	templates map[string]*template.Template
}

func newTemplater(args []string, remote bool) (*templater, error) {
	t := &templater{
		templates: make(map[string]*template.Template, len(args)),
	}
	funcMap := template.FuncMap{
		"FindByTag": tmplFindByTag,
		"Query":     tmplQuery,
		"now":       time.Now,
		"add":       func(a, b int) int { return a + b },
		"sub":       func(a, b int) int { return a + b },
		"prefix":    strings.HasPrefix,
		"suffix":    strings.HasSuffix,
		"split":     strings.Split,
		"splitN":    strings.SplitN,
		"lower":     strings.ToLower,
		"upper":     strings.ToUpper,
	}
	for _, arg := range args {
		parts := strings.SplitN(arg, ":", 2)
		file := parts[0]
		var destination string
		if len(parts) == 2 {
			destination = parts[1]
		} else {
			destination = strings.TrimSuffix(parts[0], ".tmpl")
			if parts[0] == destination {
				return nil, fmt.Errorf("template argument requires destination or source file with .tmpl suffix")
			}
		}
		tmplStr := ""
		if remote {
			str, err := loadTemplateCollins(arg)
			if err != nil {
				return nil, err
			}
			tmplStr = str
		} else {
			b, err := ioutil.ReadFile(file)
			if err != nil {
				return nil, err
			}
			tmplStr = string(b)
		}
		tmpl, err := template.New("").Funcs(funcMap).Parse(tmplStr)
		if err != nil {
			return nil, err
		}
		t.templates[destination] = tmpl
	}
	return t, nil
}

func loadTemplateCollins(path string) (string, error) {
	parts := strings.SplitN(path, "/", 2) // format: asset-tag/attribute
	assetS := "default"
	attrib := parts[0]
	if len(parts) > 1 {
		assetS = parts[0]
		attrib = parts[1]
	}
	attrib = strings.ToUpper(attrib)
	asset, _, err := client.Assets.Get(assetS)
	if err != nil {
		return "", err
	}
	template, ok := asset.Attributes["0"][attrib]
	if !ok {
		return "", fmt.Errorf("Attribute %s doesn't exist", attrib)
	}
	return template, nil
}

func (t *templater) execute() error {
	for destination, template := range t.templates {
		fh := os.Stdout
		if destination != "" {
			var err error
			fh, err = os.Create(destination)
			if err != nil {
				return err
			}
		}
		if err := template.Execute(fh, t); err != nil {
			_ = fh.Close()
			return err
		}
		if err := fh.Close(); err != nil {
			return err
		}
	}
	return nil
}
