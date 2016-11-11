package main

import "gopkg.in/tumblr/go-collins.v0/collins"

type templateData struct {
	Assets []collins.Asset
}

// FindByTag returns the asset with the given tag or an empty asset if not found.
func (t *templateData) FindByTag(tag string) collins.Asset {
	for _, asset := range t.Assets {
		if asset.Metadata.Tag == tag {
			return asset
		}
	}
	return collins.Asset{}
}
