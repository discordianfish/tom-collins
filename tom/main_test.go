package main

import "testing"

func TestFirst(t *testing.T) {
	if first("foo") != "foo" {
		t.Fatal()
	}
	if first("foo", "bar") != "foo" {
		t.Fatal()
	}
	if first("", "foo", "bar") != "foo" {
		t.Fatal()
	}
	if first("", "") != "" {
		t.Fatal()
	}
}
