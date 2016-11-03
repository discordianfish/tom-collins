package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"gopkg.in/tumblr/go-collins.v0/collins"
)

const defaultTemplate = `{{ index .Attributes "0" "HOSTNAME" }}`

func init() {
	commands["query"] = &queryCommand{}
}

type queryCommand struct{}

func (c *queryCommand) usage() string {
	return "<cql query>"
}

// Is there a way to define multiple functions on a struct at once?
func (c *queryCommand) help() string {
	return "runs a CQL query and returns matching assets"
}

func (c *queryCommand) run(args []string) error {
	var (
		flags     = flag.NewFlagSet("query", flag.ExitOnError)
		tmplStr   = flags.String("t", defaultTemplate, "template to use when printing assets")
		noNewline = flags.Bool("n", true, "don't print newlines between assets")
	)
	flags.Parse(args)
	args = flags.Args()
	if len(args) < 1 {
		return errUsage
	}
	tmpl, err := template.New("").Parse(*tmplStr)
	if err != nil {
		return err
	}
	return query(&collins.AssetFindOpts{Query: strings.Join(args, " ")}, func(asset collins.Asset) {
		err := tmpl.Execute(os.Stdout, asset)
		if *noNewline {
			fmt.Println("")
		}
		if err != nil {
			log.Println(err)
		}
	})
}
