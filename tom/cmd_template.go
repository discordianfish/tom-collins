package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/tumblr/go-collins.v0/collins"
)

func init() {
	commands["template"] = &templateCommand{}
}

type templateCommand struct{}

func (c *templateCommand) usage() string {
	return "<template:[destination]> [template...]"
}

func (c *templateCommand) help() string {
	return "renders templates with assets returned by given query."
}

func (c *templateCommand) run(args []string) error {
	var (
		flags  = flag.NewFlagSet("template", flag.ExitOnError)
		tquery = flags.String("q", "type = SERVER_NODE", "CQL query for available assets")
		remote = flags.Bool("r", false, "Use remote templates")
	)
	if err := flags.Parse(args); err != nil {
		return err
	}
	args = flags.Args()
	if len(args) < 1 {
		return errUsage
	}

	templates, err := c.loadTemplates(args, *remote)
	if err != nil {
		return err
	}

	assets := []collins.Asset{}
	if err := query(&collins.AssetFindOpts{Query: *tquery}, func(asset collins.Asset) {
		assets = append(assets, asset)
	}); err != nil {
		return err
	}

	td := &templateData{
		Assets: assets,
	}
	return c.renderTemplates(td, templates)
}

func (c *templateCommand) loadTemplates(args []string, remote bool) (map[string]*template.Template, error) {
	templates := make(map[string]*template.Template, len(args))
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
			str, err := c.loadTemplateCollins(arg)
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
		log.Printf("using '%s'", tmplStr)
		tmpl, err := template.New("").Parse(tmplStr)
		if err != nil {
			return nil, err
		}
		templates[destination] = tmpl
	}
	return templates, nil
}
func (c *templateCommand) loadTemplateCollins(path string) (string, error) {
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

func (c *templateCommand) renderTemplates(td *templateData, templates map[string]*template.Template) error {
	for destination, template := range templates {
		fh := os.Stdout
		if destination != "" {
			var err error
			fh, err = os.Create(destination)
			if err != nil {
				return err
			}
		}
		if err := template.Execute(fh, td); err != nil {
			_ = fh.Close()
			return err
		}
		if err := fh.Close(); err != nil {
			return err
		}
	}
	return nil
}
