package main

import (
	"fmt"
)

func main() {
	m := map[string]string{
		"a":  "b",
		"a1": "b1",
	}
	fmt.Println("before", m)
	test(m)
	fmt.Println("after", m)
}

func test(m map[string]string) {
	m["a"] = "b0"
	fmt.Println("changed", m)
}
