package main

import "testing"

func TestFirst(t *testing.T) {
	if first("foo") != "foo" {
		t.Fatal("err")
	}
	if first("foo", "bar") != "foo" {
		t.Fatal("err")
	}
	if first("", "foo", "bar") != "foo" {
		t.Fatal("err")
	}
	if first("", "") != "" {
		t.Fatal("err")
	}
}
