package main

import "fmt"

//procm:use=doesnotexists
type delayed struct{}

//procm:use=eval
func evalTest() {
	n := eval(`
		1 + 1
    `)
	fmt.Print(n)
}

func main() {
	evalTest()
}
