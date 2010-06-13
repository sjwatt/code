package main

import (
	"fmt"
	"time"
)

const (
	second int64 = 1e9
)


func thread0(ch chan string) {
	fmt.Println(time.Nanoseconds(),"Thread 0 online.")
	printer := <- ch
	fmt.Println(time.Nanoseconds(),printer,"This is thread0")
}

func thread1(ch chan string) {
	fmt.Println(time.Nanoseconds(),"Thread 1 online.")
	printer := <- ch
	fmt.Println(time.Nanoseconds(),printer,"This is thread1")
} 

func main() {
	fmt.Println(time.Nanoseconds(),"Begin Program")
	ch := make(chan string)
	go thread0(ch)	
	go thread1(ch)

	for i := 0; i < 2; i++ {
		ch <- "Hello World! " + string(i)
	}

	time.Sleep(second)	
	fmt.Println(time.Nanoseconds(),"Exiting")
}
