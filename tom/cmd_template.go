package main

import "flag"

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
		remote = flags.Bool("r", false, "Use remote templates")
	)
	if err := flags.Parse(args); err != nil {
		return err
	}
	args = flags.Args()
	if len(args) < 1 {
		return errUsage
	}

	templater, err := newTemplater(args, *remote)
	if err != nil {
		return err
	}
	return templater.execute()
}
