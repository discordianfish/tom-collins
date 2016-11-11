package main

import (
	"fmt"

	"gopkg.in/tumblr/go-collins.v0/collins"
)

const (
	queryStrictFS = "%s = ^%s$"
)

func query(opts *collins.AssetFindOpts, f func(collins.Asset)) error {
	for {
		assets, resp, err := client.Assets.Find(opts)
		if err != nil {
			return err
		}
		for _, asset := range assets {
			f(asset)
		}
		if resp.NextPage == resp.CurrentPage {
			break
		}
		opts.PageOpts.Page++
	}
	return nil
}

func findAssets(tag, hostname, queryStr string, f func(collins.Asset)) error {
	q := ""
	switch {
	case tag != "":
		q = fmt.Sprintf(queryStrictFS, "tag", tag)
	case hostname != "":
		q = fmt.Sprintf(queryStrictFS, "hostname", hostname)
	case queryStr != "":
		q = queryStr
	default:
		return fmt.Errorf("Require -h, -t or -q")
	}
	return query(&collins.AssetFindOpts{Query: q}, f)
}
