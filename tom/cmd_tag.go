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

	it, err := NewFindIterator(&collins.AssetFindOpts{Attribute: "hostname;^" + args[0] + "$"})
	if err != nil {
		return err
	}
	for it.Next() {
		fmt.Println(it.Value().Tag)
	}
	return it.Err()
}
