/*
 * Function print binary form of integer
 *
 *
 */
package main

import "fmt"

func main() {

	PrintBinaryform( 13 )
	fmt.Println()

}

func PrintBinaryform(number int) {

	if number <= 1 {
		fmt.Print(number)
		return // Kick out of the recursion
	}

	remainder := number % 2
	PrintBinaryform(number >> 1)
	fmt.Print(remainder)

}
