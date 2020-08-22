package main

import "testing"

func Test_run(t *testing.T) {
	err := run(`./testdata/proto`, `/testdata/proto_out`)
	if err != nil {
		t.Fatal(err)
	}
}
