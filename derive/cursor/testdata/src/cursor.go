package main

import (
	"bytes"
	"fmt"
	"image"
	"time"

	"github.com/tdakkota/cursor"
)

func equalBytes(a, b []byte) {
	if !bytes.Equal(a, b) {
		fmt.Println(a, "\n", b)
		panic("expected equal")
	}
}

type codec interface {
	cursor.Appender
	cursor.Scanner
}

func testCursor(s codec, data []byte) {
	cur := cursor.NewCursor(nil)
	err := s.Append(cur)
	if err != nil {
		panic(err)
	}

	equalBytes(data, cur.Buffer())

	// test unmarshal
	cur = cursor.NewCursor(cur.Buffer())
	err = s.Scan(cur)
	if err != nil {
		panic(err)
	}
}

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

func testStructData() (*testStruct, []byte) {
	return &testStruct{
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

func testStructTest() {
	testCursor(testStructData())
}

//procm:use=derive_binary
type attr struct {
	predicate int8
	value     int8 `if:"$m.predicate >= 10"`
	skip      int  `skip:""`
}

//procm:use=derive_binary
type attrArray struct {
	length int8
	array  []int8 `length:"uint8($m.length)"`
}

func testConditional() {
	testCursor(
		&attr{10, 42, 11},
		[]byte{10, 42},
	)
	testCursor(
		&attr{4, 0, 11},
		[]byte{4},
	)
}

func testAtrrArray() {
	testCursor(&attrArray{2, []int8{1, 2}}, []byte{2, 1, 2})
	testCursor(&attrArray{0, []int8{}}, []byte{0})
}

func main() {
	testStructTest()
	testConditional()
}
