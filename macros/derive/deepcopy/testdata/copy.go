package main

import "time"

//procm:use=derive_copy
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

//procm:use=derive_copy
type T struct {
	value        string
	array        [5]int
	darray       [5][5]int
	arrayImport  [5]time.Duration
	slice        []int
	sliceEmpty   []struct{}
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

func main() {
	t := T{
		value:       "abc",
		array:       [5]int{1, 2, 3, 4, 5},
		darray:      [5][5]int{{1, 2, 3, 4, 5}},
		arrayImport: [5]time.Duration{1, 2, 3, 4, 5},
		slice:       []int{1, 2, 3},
		dslice:      [][]int{{1, 2, 3}},
		sliceEmpty:  []struct{}{{}},
		dsliceStruct: [][]struct {
			slice  []T
			dslice [][]int
		}{
			{
				struct {
					slice  []T
					dslice [][]int
				}{slice: []T{}, dslice: [][]int{{1, 2, 3}}},
			},
		},
		sliceImport: nil,
		m: map[string]struct{}{
			"abc": {},
		},
		mImport: map[string]time.Duration{
			"abc": 1,
		},
		notImpl: NotImpl{
			value: "abc",
		},
		pNotImpl: &NotImpl{
			value: "abc",
		},
		impl: Impl{
			value: "abc",
		},
		pImpl: &Impl{
			value: "abc",
		},
	}

	var t2 T
	t2 = t.Copy()
	t.value = "abcd"
	t.array[0] = 10
	t.slice[0] = 10
	t.dslice[0][0] = 10
	t.sliceEmpty = append(t.sliceEmpty, struct{}{})
	t.dsliceStruct[0][0].dslice[0][0] = 10
	t.m["abcd"] = struct{}{}
	t.mImport["abc"] = 10
	t.notImpl.value = "abcd"
	t.pNotImpl.value = "abcd"
	t.impl.value = "abcd"
	t.pImpl.value = "abcd"

	if t2.value != "abc" {
		panic("expected t2.value is deep copied")
	}
	if t2.array[0] != 1 {
		panic("expected t.array[0] is deep copied")
	}
	if t2.slice[0] != 1 {
		panic("expected t.slice[0] is deep copied")
	}
	if len(t2.sliceEmpty) != 1 {
		panic("expected t.sliceEmpty is deep copied")
	}
	if t2.dslice[0][0] != 1 {
		panic("expected t.dslice[0][0] is deep copied")
	}
	if t2.dsliceStruct[0][0].dslice[0][0] != 1 {
		panic("expected t.dsliceStruct[0][0].dslice[0][0] is deep copied")
	}
	if len(t2.m) != 1 {
		panic("expected t2.m is deep copied")
	}
	if t2.mImport["abc"] != 1 {
		panic("expected t2.mImport[\"abc\"] is deep copied")
	}
	if t2.notImpl.value != "abc" {
		panic("expected t2.notImpl.value is deep copied")
	}
	if t2.pNotImpl.value != "abc" {
		panic("expected t2.pNotImpl.value is deep copied")
	}
	if t2.impl.value != "abc" {
		panic("expected t2.impl.value is deep copied")
	}
	if t2.pImpl.value != "abc" {
		panic("expected t2.pImpl.value is deep copied")
	}
}
