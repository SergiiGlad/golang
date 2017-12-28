	package main

	import 	"fmt"



	func main() {

		array := []int{1,2,3,4,5}

		fmt.Println("Sum ", array, " = ", sum(array))

		fmt.Println("Multiply ", array, " = ", multiply(array))

	}





	// sums all the number in a slice of number and return result

	func sum(array1 []int) (sum int) {

		i := len( array1 );

		for	 i > 0 {

			sum += array1[i-1]
			i--
		}

		return
	}


	// multiplies all the number in a slice of number and return result
	func multiply(array2 []int) int {

		multi := 1

		for i := 0; i < len(array2); i++ {

			multi *= array2[i]

		}

		return multi

	}

