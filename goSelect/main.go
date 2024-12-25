package main

import (
	"fmt"
	"time"
)

func simpleSelectExample() {
	ch := make(chan any)

	start := time.Now()
	go func() {
		time.Sleep(3 * time.Second)
		close(ch)
	}()

	fmt.Printf("Blocking...\n")
	select {
	case <-ch:
		fmt.Printf("Unblocked %.2v later!\n", time.Since(start))
	}
}

func selectEqualChanceExample() {
	c1 := make(chan interface{})
	close(c1)
	c2 := make(chan interface{})
	close(c2)

	var c1Count, c2Count int
	for i := 1000; i >= 0; i-- {
		select {
		case <-c1:
			c1Count++
		case <-c2:
			c2Count++
		}
	}
	fmt.Printf("c1Count: %d\nc2Count: %d\n", c1Count, c2Count)
}

func selectTimeoutExample() {
	now := time.Now()
	ch1 := make(chan any)

	select {
	case <-ch1:
	case <-time.After(2 * time.Second):
		fmt.Printf("Timed out after %.2v seconds\n", time.Since(now))
	}
}

func selectDefaultExample() {
    ch := make(chan any)
    counter := 0

    go func() {
        time.Sleep(4 * time.Second)
        close(ch)
    }()

    loop:
    for {
        select {
        case <-ch:
            break loop
        default:
        }
        counter++
        time.Sleep(1 * time.Second)
    }
    fmt.Println("counter:", counter)
}

func main() {
	// simpleSelectExample()

    // selectEqualChanceExample()
	// selectTimeoutExample()
    selectDefaultExample()
}
