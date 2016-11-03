package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/tumblr/go-collins.v0/collins"
)

var (
	uri      = flag.String("uri", first(os.Getenv("COLLINS_URL"), "http://localhost:9000/api"), "url to collins api")
	user     = flag.String("user", first(os.Getenv("COLLINS_USER"), "blake"), "collins user")
	password = flag.String("password", first(os.Getenv("COLLINS_PASSWORD"), "admin:first"), "collins password")

	client *collins.Client

	commands = map[string]command{}
	errUsage = errors.New("Invalid usage")
)

type command interface {
	usage() string
	help() string
	run(args []string) error
}

func first(ss ...string) string {
	for _, s := range ss {
		if s != "" {
			return s
		}
	}
	return ""
}

func forAssets(opts *collins.AssetFindOpts, f func(*collins.Asset)) error {
	for {
		assets, resp, err := client.Assets.Find(opts)
		if err != nil {
			return err
		}
		for _, asset := range assets {
			f(&asset)
		}
		if resp.NextPage == resp.CurrentPage {
			break
		}
		opts.PageOpts.Page++
	}
	return nil
}

func main() {
	flag.Usage = func() { printUsage("") }
	flag.Parse()
	var err error
	client, err = collins.NewClient(*user, *password, *uri)
	if err != nil {
		log.Fatal(err)
	}
	args := flag.Args()
	if len(args) < 1 {
		printUsage("")
	}
	cmd, ok := commands[args[0]]
	if !ok {
		fmt.Fprintln(os.Stderr, "Invalid command", args[0])
		printUsage("")
	}
	if err := cmd.run(args[1:]); err != nil {
		if err == errUsage {
			printUsage("Invalid argument. Expected: " + args[0] + " " + cmd.usage())
		}
		log.Fatal(err)
	}
}

func printUsage(message string) {
	if message != "" {
		message = message + "\n"
	}
	fmt.Fprintln(os.Stderr, message+"Usage: collins-cli [options] sub-command")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "\nSub commands:")
	for n, c := range commands {
		fmt.Fprintln(os.Stderr, "  -", n, c.usage())
		fmt.Fprintln(os.Stderr, "\t", n, c.help())
	}
	os.Exit(1)
}
