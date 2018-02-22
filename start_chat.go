package main

import (
	"fmt"
	"time"
	"math/rand"
)

const BYE int = 1
const WELCOME int = 2

var (
	state int 
	c chan int 

)


func chat() {

		defer func() {
			c<-BYE
		}()

	

}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func start(n int) {

	
	r := rand.Intn(5) + 1
	
	go chat()

	fmt.Printf("chat % d time chatting : %v\n", n ,r)

	time.Sleep(time.Second * time.Duration(r))
	
	
}

func main () {

	c = make(chan int )
	
	var f func(int)

	f = start

	fmt.Printf("f : %v\n", f)

				
	for i :=0 ; i < 5; i++ {
		go f(i+1)
	}

	
	fmt.Printf("from chan : %v\n", <- c) 
	fmt.Printf("from chan : %v\n", <- c) 
	fmt.Printf("from chan : %v\n", <- c) 
	fmt.Printf("from chan : %v\n", <- c) 
	fmt.Printf("from chan : %v\n", <- c) 
	
	
	
	
	time.Sleep(time.Second*3)
	
	

}