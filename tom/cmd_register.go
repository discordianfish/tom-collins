package main

import "fmt"

func init() {
	commands["register"] = &noopCommand{}
}

type registerCommand struct{}

func (c *registerCommand) usage() string {
	return "<hostname>"
}

func (c *registerCommand) help() string {
	return "Return register"
}

func (c *registerCommand) run(args []string) error {
	fmt.Println("register")
	return nil
}
