package main

import "gopkg.in/tumblr/go-collins.v0/collins"

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
