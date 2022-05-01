package main

import (
	"fmt"
	"time"
)

func main() {
	items := 100
	concurrency := 10
	in := make(chan int)
	ret := make(chan error)

	for x := 0; x < concurrency; x++ {
		go worker(x, in, ret)
	}

	go func() {
		for x := 0; x < items; x++ {
			in <- x
		}
		close(in)
	}()

	for err := range ret {
		if err != nil {
			fmt.Println(err.Error())
			break
		}
	}

}

func worker(workerNumber int, in chan int, ret chan error) {
	fmt.Println(fmt.Sprintf("workder %d started working...", workerNumber))
	for x := range in {
		if x == 95 {
			ret <- fmt.Errorf("something not right")
		} else {
			ret <- nil
		}
		fmt.Println(fmt.Sprintf("Input value: %d", x))
		fmt.Println("waiting 5s...")
		time.Sleep(5 * time.Second)
		fmt.Println("end waiting...")
	}

	ret <- fmt.Errorf("the end")
}
