package main

import (
	"flag"
	"fmt"
	"log"

	"gopkg.in/tumblr/go-collins.v0/collins"
)

func init() {
	commands["update"] = &updateCommand{}
}

type updateCommand struct{}

func (c *updateCommand) usage() string {
	return "<attribute> [value]"
}

func (c *updateCommand) help() string {
	return "Update attributes on assets"
}

func (c *updateCommand) run(args []string) error {
	var (
		flags    = flag.NewFlagSet("set", flag.ExitOnError)
		remove   = flags.Bool("d", false, "Remove attribute")
		hostname = flags.String("h", "", "Hostname of asset to update attribute on")
		tag      = flags.String("t", "", "Tag of asset to update attributes on")
		query    = flags.String("q", "", "Tag of asset to update attributes on")
	)
	if err := flags.Parse(args); err != nil {
		return err
	}
	args = flags.Args()
	if len(args) < 1 {
		return errUsage
	}
	attribute := args[0]
	if *remove {
		return findAssets(*tag,
			*hostname,
			*query,
			func(a collins.Asset) {
				_, err := client.Assets.DeleteAttribute(a.Metadata.Tag, attribute)
				if err != nil {
					log.Fatal(err)
				}
			},
		)
	}
	if len(args) != 2 {
		return fmt.Errorf("Value required for setting attributes")
	}
	opts := &collins.AssetUpdateOpts{
		Attribute: attribute + ";" + args[1],
	}
	return findAssets(
		*tag,
		*hostname,
		*query,
		func(a collins.Asset) {
			log.Println("Setting", attribute, "=", args[1], "on", a.Metadata.Tag)
			_, err := client.Assets.Update(a.Metadata.Tag, opts)
			if err != nil {
				log.Fatal(err)
			}
		},
	)
}
