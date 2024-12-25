package main

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"
)

func chanReadWriteExample() {
	intCh := make(chan int)
	go func() {
		intCh <- 1
		close(intCh)
	}()

	for _ = range 2 {
		val, ok := <-intCh
		fmt.Printf("(%v): %d\n", ok, val)
	}
}

func chanReadWriteExample_Deadline() {
	intCh := make(chan int)
	defer close(intCh)

	intCh <- 1

	val, ok := <-intCh
	fmt.Printf("(%v): %d\n", ok, val)
}

// func invalildReadAndWrite() {
//     readOnlyCh := make(chan<- int)
//     writeOnlyCh := make(<-chan int)

//     // invalid op: write to read-only channel
//     readOnlyCh <- 1

//     // invalid op: read from write-only channel
//     <-writeOnlyCh

// }

func closeChannelAsSignal() {
	startRace := make(chan interface{})
	wg := sync.WaitGroup{}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			<-startRace // Wait for start signal
			fmt.Printf("Runner %d started (%s)\n", id, time.Now().Format("15:04:05.000"))
		}(i)
	}

	fmt.Printf("Preparing the race... (%s)\n", time.Now().Format("15:04:05.000"))
	time.Sleep(2 * time.Second)
	fmt.Printf("...GO! (%s)\n", time.Now().Format("15:04:05.000"))

	close(startRace) // Signal all runners to start
	wg.Wait()        // Wait for all runners to finish
}

func bufferedChannelExample() {
	// Create a buffer to store output for synchronized printing
	var stdoutBuff bytes.Buffer
	defer stdoutBuff.WriteTo(os.Stdout)

	// Create buffered channel with capacity 4
	intStream := make(chan int, 4)

	// Producer goroutine
	go func() {
		defer close(intStream)
		defer fmt.Fprintln(&stdoutBuff, "Producer Done.")
		for i := 0; i < 6; i++ {
			fmt.Fprintf(&stdoutBuff, "+ Sending: %d (%s)\n", i, time.Now().Format("15:04:05.000"))
			intStream <- i // Will block when buffer is full (after 4 items)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	// Consumer: read until channel is closed and all buffered values are received
	for integer := range intStream {
		fmt.Fprintf(&stdoutBuff, "- Received %v (%s).\n", integer, time.Now().Format("15:04:05.000"))
	}
}

func main() {
	chanReadWriteExample()
	// invalildReadAndWrite()
	closeChannelAsSignal()
	// bufferedChannelExample()
}
