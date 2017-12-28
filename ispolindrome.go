package main

import (
	"basicfunc"
	"fmt"
	"strings"
)

func main() {

	str := "Radar"

	fmt.Printf("Word %q is polindrom : ", str)
	fmt.Println( isPalindromeCompare( str ) )

}

func isPalindromeCompare(s string) bool {

	result := strings.Compare(s, basicfunc.Reverse(s))

	for result == 0 {
		return true
	}

	return false

}