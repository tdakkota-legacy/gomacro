package main

import "time"

//procm:use=derive_zero
type M struct {
	i int
	t time.Duration
}

//procm:use=derive_zero
type T struct {
	v             int
	s             string
	abc           [5]int
	abcd          [5]time.Duration
	slice         []int
	importedSlice []time.Duration
	m             map[string]struct{}
	m2            map[string]time.Duration
	em            M
}

//procm:use=derive_zero
type Comparable struct {
	s string
}

func main() {
	var t T
	if !t.Zero() {
		panic("expected t is zero")
	}
	t.s = "abc"
	if t.Zero() {
		panic("expected t is not zero")
	}

	var c Comparable
	if !c.Zero() {
		panic("expected c is zero")
	}
	c.s = "abc"
	if c.Zero() {
		panic("expected c is not zero")
	}
}
