package main

import (
	"bytes"
	"fmt"
	"image"
	"time"

	"github.com/tdakkota/cursor"
)

//procm:use=derive_binary
type Flag byte

type cycle struct {
	s []testStruct
}

//procm:use=derive_binary
type testStruct struct {
	uint
	byte
	uint16
	uint32
	uint64
	int
	int8
	int16
	int32
	int64
	float32
	float64
	bytes []byte
	ints  [2]int8
	string
	selfCycle []testStruct

	dur        time.Duration
	importable image.Point
	cycle      cycle
	flag       Flag
}

func data() (testStruct, []byte) {
	return testStruct{
			10,
			1,
			2,
			3,
			4,
			11,
			5,
			6,
			7,
			8,
			0,
			0,
			[]byte{'x', 'y'},
			[2]int8{1, 2},
			"abc",
			nil,
			10,
			image.Point{X: 3, Y: 4},
			cycle{},
			9,
		}, []byte{
			10, 0, 0, 0, 0, 0, 0, 0, // uint
			1,    // byte
			2, 0, // uint16
			3, 0, 0, 0, // uint32
			4, 0, 0, 0, 0, 0, 0, 0, // uint64
			11, 0, 0, 0, 0, 0, 0, 0, // int
			5,    // int8
			6, 0, // int16
			7, 0, 0, 0, // int32
			8, 0, 0, 0, 0, 0, 0, 0, // int64
			0, 0, 0, 0, // float32
			0, 0, 0, 0, 0, 0, 0, 0, // float64
			2, 'x', 'y', // bytes
			1, 2, // ints
			3, 'a', 'b', 'c', // string
			0,                       // selfCycle
			10, 0, 0, 0, 0, 0, 0, 0, // dur
			3, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, // importable
			0, // cycle
			9, // flag
		}
}

func main() {
	s, b := data()

	cur := cursor.NewCursor(nil)
	err := s.Append(cur)
	if err != nil {
		panic(err)
	}

	if !bytes.Equal(b, cur.Buffer()) {
		fmt.Println(b)
		fmt.Println(cur.Buffer())
		panic("expected equal")
	}
}
