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

func buffer(buf int, c chan int) {

	for i := 0; i < buf; i++ {
		fmt.Printf(" %d send to channel\n", i)
		c <- i
	}
}

func main() {

	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}

	c := make(chan int, 2)

	go summary(arr, c)              //45 long running
	go summary(arr[len(arr)/2:], c) //30
	go summary(arr[:1], c)          //1
	go summary(arr[:2], c)          //3

	x := <-c // 45
	y := <-c // 30
	x = <-c  //1

	fmt.Printf("x :%d  y : %d\n", x, y) //

	//Sends to a buffered channel block only when the buffer is full. Receives block when the buffer is empty.
	go buffer(5, c)
	fmt.Printf("read from chan %d\n", <-c)
	fmt.Printf("read from chan %d\n", <-c)

}
