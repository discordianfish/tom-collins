package main

import "fmt"

func init() {
	commands["register"] = &registerCommand{}
}

type registerCommand struct{}

func (c *registerCommand) usage() string {
	return "<hostname>"
}

func (c *registerCommand) help() string {
	return "is to be implement"
}

func (c *registerCommand) run(args []string) error {
	fmt.Println("register")
	return nil
}
