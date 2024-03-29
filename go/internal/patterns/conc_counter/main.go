package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	data  chan int
	total chan int
}

func NewCounter() *Counter {
	data := make(chan int)
	total := make(chan int)

	c := Counter{data, total}

	go func() {
		var count int
		for {
			select {
			case increment := <-data:
				count += increment
			case total <- count:
			}
		}
	}()

	return &c
}

func (c *Counter) Add(v int) {
	c.data <- v
}

func (c *Counter) Total() int {
	return <-c.total
}

func main() {
	counter := NewCounter()

	var wg sync.WaitGroup
	wg.Add(1000000)

	for i := 0; i < 1000000; i++ {
		go func() {
			counter.Add(1) // goroutine
			wg.Done()
		}()
	}

	counter.Add(1) // main goroutine

	wg.Wait()

	fmt.Println(counter.Total())
}
