package main

import (
	"fmt"

	"gopkg.in/tumblr/go-collins.v0/collins"
)

const (
	queryStrictFS = "%s = ^%s$"
)

type findIterator struct {
	assets   []collins.Asset
	response *collins.Response
	opts     *collins.AssetFindOpts
	current  int
	err      error
}

// NewFindIterator returns a new findIterator
func NewFindIterator(opts *collins.AssetFindOpts) (*findIterator, error) {
	i := &findIterator{opts: opts, current: -1} // -1 so first Next() sets to 0
	return i, i.find()
}

// Value returns the current asset
func (i *findIterator) Value() collins.Asset {
	return i.assets[i.current]
}

func (i *findIterator) find() (err error) {
	i.assets, i.response, err = client.Assets.Find(i.opts)
	return err
}

// Next advances to the next asset and returns true. If there is no next asset
// or an error occurs, it returns false.
func (i *findIterator) Next() bool {
	if i.current < len(i.assets)-1 {
		i.current++
		return true
	}
	if i.response.NextPage == i.response.CurrentPage {
		return false
	}
	i.opts.PageOpts.Page++
	i.current = 0
	i.err = i.find()
	return i.err == nil
}

func (i *findIterator) Err() error {
	return i.err
}

func newFindOpts(tag, hostname, query string) (*collins.AssetFindOpts, error) {
	q := ""
	switch {
	case tag != "":
		q = fmt.Sprintf(queryStrictFS, "tag", tag)
	case hostname != "":
		q = fmt.Sprintf(queryStrictFS, "hostname", hostname)
	case query != "":
		q = query
	default:
		return nil, fmt.Errorf("Require -h, -t or -q")
	}
	return &collins.AssetFindOpts{Query: q}, nil
}
