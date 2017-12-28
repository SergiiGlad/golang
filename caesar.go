/*
 * Program does string as encrypted Caesar Cipher with some key
 * if you want decrypte you should use key = 26 - key
 *	for example
 *  asdf 10 --> kcnp
 *  kcnp 16 --> asdf
 *
 */
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const UP rune = 65
const LOW rune = 97

func main() {

	var (
		v    int // key number
		err  error
		text string // input string
		key  string // key number string
	)

	args := os.Args[1:] // argument command line

	if len(args) == 0 {
		//Parse input key
		text = ReadString("Caught key:")
		key = text

	} else {

		key = args[0]
	}

	if v, err = strconv.Atoi(key); err != nil {
		fmt.Errorf("%T, %v\n", v, v)
		fmt.Println(" key shoud be number")
		return
	}

	//Parse input text
	text = ReadString("Caught text which to crypt:")

	fmt.Printf("Encrypted text : %s \n ", FunctionCrypt(text, rune(v)))

}

func ReadString(quest string) string {

	reader := bufio.NewReader(os.Stdin)

	fmt.Println(quest)
	text, _ := reader.ReadString('\n')

	// convert CRLF to LF

	text = strings.Replace(text, "\n", "", -1)

	return text
}

func FunctionCrypt(s string, k rune) (result string) {

	length := len(s)

	for i := 0; i < length; i++ {

		ch := rune(s[i])

		if ch >= 'A' && ch <= 'Z' {
			ch = UP + (ch-UP+k)%26
		}

		if ch >= 'a' && ch <= 'z' {
			ch = LOW + (ch-LOW+k)%26
		}

		result += string(ch)

	} // for

	return
}
