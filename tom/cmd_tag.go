package main

import (
	"fmt"

	"gopkg.in/tumblr/go-collins.v0/collins"
)

func init() {
	commands["tag"] = &tagCommand{}
}

type tagCommand struct{}

func (c *tagCommand) usage() string {
	return "<hostname>"
}

func (c *tagCommand) help() string {
	return "returns tag for hostname"
}

func (c *tagCommand) run(args []string) error {
	if len(args) != 1 {
		return errUsage
	}
	opts := &collins.AssetFindOpts{
		Attribute: "hostname;^" + args[0] + "$",
	}
	return forAssets(opts, func(asset *collins.Asset) {
		fmt.Println(asset.Tag)
	})
}
