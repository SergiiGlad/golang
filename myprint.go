	package main

	import (
		"fmt"
		"os"

	)

	func main() {

		// Print arguments from command line if they exists

		args := os.Args[1:]

		if  len(args) > 0 {

			fmt.Printf("\"Hello, %v!\"\n", args )


		} else {

			fmt.Println("Usage: myprint argument1 [argument2] ....")
			MyPrint("Hello, Rob Pike")
		}
	}


	// the function takes string as argument and prints some frase

	func MyPrint ( str string ) {
		fmt.Printf("%q\n", str )
	}