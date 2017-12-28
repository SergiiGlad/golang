package main

import "fmt"

func main() {

	test := "aFGAEJH B SFHWREHaaaaaaaaaBDSFLGB;QARGHERhm\\ bDFSG,SDBG"

	m := make(map[rune]int)

	ch := 'a'
	count := 1

	m[ch] = count

	for i := 0; i < len(test); i++ {
		c := test[i]
		if c != ' ' {
			// Increment existing value in HashMap.
			// ... Start with zero if no key exists.
			// retrieve an element
			elem := m[rune(c)]
			elem++
			m[rune(c)] = elem
		}
	}

	for key := range m {
		fmt.Printf("%c : %d\n", key, m[rune(key)])
	}

}
