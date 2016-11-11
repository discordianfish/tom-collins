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
	findOpts, err := newFindOpts(*tag, *hostname, *query)
	if err != nil {
		return err
	}
	if *remove {
		return c.remove(findOpts, attribute)
	}

	if len(args) != 2 {
		return fmt.Errorf("Value required for setting attributes")
	}
	return c.set(findOpts, attribute, args[1])
}

func (c *updateCommand) remove(findOpts *collins.AssetFindOpts, attribute string) error {
	it, err := NewFindIterator(findOpts)
	if err != nil {
		return err
	}
	for it.Next() {
		asset := it.Value()
		_, err := client.Assets.DeleteAttribute(asset.Metadata.Tag, attribute)
		if err != nil {
			return err
		}
	}
	return it.Err()
}

func (c *updateCommand) set(findOpts *collins.AssetFindOpts, attribute, value string) error {
	opts := &collins.AssetUpdateOpts{
		Attribute: attribute + ";" + value,
	}
	it, err := NewFindIterator(findOpts)
	if err != nil {
		return err
	}
	for it.Next() {
		asset := it.Value()
		log.Println("Setting", attribute, "=", value, "on", asset.Metadata.Tag)
		_, err := client.Assets.Update(asset.Metadata.Tag, opts)
		if err != nil {
			return err
		}
	}
	return it.Err()
}
