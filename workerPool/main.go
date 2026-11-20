package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

func process(jobChan <-chan int, resultChan chan<- int, wg *sync.WaitGroup) {
	// Done runs once, when the worker exits (jobChan closed and drained)
	defer wg.Done()
	// a worker is long-lived: it keeps taking jobs until jobChan is closed
	for val := range jobChan {
		// create a random time, 0-499ms
		randomTime := rand.IntN(500)
		time.Sleep(time.Duration(randomTime) * time.Millisecond)
		resultChan <- val
	}
}

func main() {
	var wg sync.WaitGroup
	
	var input int
	fmt.Print("Enter the number of images: ")
	fmt.Scanln(&input)
	fmt.Println("the number of images are: ", input)
	
	jobChan := make(chan int, 8)
	resultChan := make(chan int)
	
	// start the fixed pool of 8 workers once, before any job is sent
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go process(jobChan, resultChan, &wg)
	}

	// producer in its own goroutine so main is free to read results;
	// only the sender closes jobChan
	go func() {
		for i := 0; i < input; i++ {
			jobChan <- i
		}
		close(jobChan)
	}()

	// coordinator closes resultChan only after every worker has exited
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// range ends when resultChan is closed
	for res := range resultChan {
		fmt.Println("Processed image: ", res)
	}
}