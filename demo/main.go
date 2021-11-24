package main

import "fmt"

func main() {
	s := "()"
	isValid(s)
}

func isValid(s string) bool {
	aa := map[string]int{
		"(": 1,
		")": 1,
		"{": 2,
		"}": 2,
		"[": 3,
		"]": 3,
	}
	fmt.Println(aa)

	return true
}
