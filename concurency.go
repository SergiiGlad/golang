package main

import (
	"fmt"
)

func summary(arr []int, c chan int) {

	sum := 0
	for i := 0; i < len(arr); i++ {
		sum += arr[i]
	}

	c <- sum

}

func main() {

	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}

	c := make(chan int)

	go summary(arr, c)              //45
	go summary(arr[:1], c)          //1
	go summary(arr[:2], c)          //3
	go summary(arr[len(arr)/2:], c) //30

	x := <-c // 45
	y := <-c // 1
	x = <-c  //3

	fmt.Printf("x :%x  y : %d\n", x, y) // 3 1 or 1 45

}
