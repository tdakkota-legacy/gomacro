package main

import "fmt"

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
