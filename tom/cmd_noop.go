package main

import "fmt"

func init() {
	commands["noop"] = &noopCommand{}
}

type noopCommand struct{}

func (c *noopCommand) usage() string {
	return "<hostname>"
}

func (c *noopCommand) help() string {
	return "Return noop"
}

func (c *noopCommand) run(args []string) error {
	fmt.Println("noop")
	return nil
}
