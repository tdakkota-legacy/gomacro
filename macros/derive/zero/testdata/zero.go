package main

import "time"

//procm:use=derive_zero
type Impl struct {
	value string
	t     time.Duration
}

type NotImpl struct {
	value       string
	array       [5]int
	darray      [5][5]int
	arrayImport [5]time.Duration
	slice       []int
	dslice      [][]int
}

//procm:use=derive_zero
type T struct {
	value        string
	array        [5]int
	darray       [5][5]int
	arrayImport  [5]time.Duration
	slice        []int
	dslice       [][]int
	dsliceStruct [][]struct {
		slice  []T
		dslice [][]int
	}
	sliceImport []time.Duration
	m           map[string]struct{}
	mImport     map[string]time.Duration
	notImpl     NotImpl
	pNotImpl    *NotImpl
	impl        Impl
	pImpl       *Impl
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
	t.value = "abc"
	if t.Zero() {
		panic("expected t is not zero")
	}

	var t2 T
	if !t2.Zero() {
		panic("expected t is zero")
	}
	t2.notImpl.value = "abc"
	if t2.Zero() {
		panic("expected t is not zero")
	}

	var t3 T
	if !t3.Zero() {
		panic("expected t is zero")
	}
	t3.impl.value = "abc"
	if t3.Zero() {
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
