package main

import "fmt"

func main() {

	array_b := []byte{5, 8, 9, 4, 7}

	histogramPrint(array_b)


}


// func print histogram

func histogramPrint(array []byte) {

	for i := 0; i < len(array); i++ {
		fmt.Println(PrintSymbolString( array[i]) )
	}

}


// function create string with  '*' , length as parameter
func PrintSymbolString(count byte) string {

	var (
		str string
		i byte = 0
	)
	for ; i < count; i++ {
		str = str + "*"
	}

	return str
}