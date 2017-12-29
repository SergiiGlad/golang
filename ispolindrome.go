package main

import (
	"fmt"
	basicfunc "github.com/SergiiGlad/golang/basicfunc"
	"strings"
)

func main() {

	str := "Radar"

	fmt.Printf("Word %q is polindrom : ", str)
	fmt.Println(isPalindromeCompare(str))

}

func isPalindromeCompare(s string) bool {

	result := strings.Compare(s, basicfunc.Reverse(s))

	for result == 0 {
		return true
	}

	return false

}
