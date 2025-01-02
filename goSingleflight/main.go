package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/singleflight"
)

var callCount atomic.Int32
var wg sync.WaitGroup

// Simulate a function that fetches data from a database
func fetchData() (interface{}, error) {
	callCount.Add(1)
	time.Sleep(100 * time.Millisecond)
	return rand.Intn(100), nil
}

// Wrap the fetchData function with singleflight
func fetchDataWrapper(g *singleflight.Group, id int) error {
	defer wg.Done()

	time.Sleep(time.Duration(id) * 40 * time.Millisecond)
	// Assign a unique key to track these requests
	v, err, shared := g.Do("key-fetch-data", fetchData)
	if err != nil {
		return err
	}

	fmt.Printf("Goroutine %d: result: %v, shared: %v\n", id, v, shared)
	return nil
}

func SingleflightExample() {
    fmt.Println("\nSingleflight Example")
	var g singleflight.Group

	// 5 goroutines to fetch the same data
	const numGoroutines = 5
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go fetchDataWrapper(&g, i)
	}

	wg.Wait()
	fmt.Printf("Function was called %d times\n", callCount.Load())
}

func goChanExample() {
    fmt.Println("\nGoChan Example")
	fetchData := func() (interface{}, error) {
		// infinite blocking
		select {}
		return rand.Intn(100), nil
	}
	wg := sync.WaitGroup{}
	g := singleflight.Group{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ch := g.DoChan("key-fetch-data", fetchData)
	select {
	case res := <-ch:
		if res.Err != nil {
			fmt.Printf("error: %v\n", res.Err)
			return
		}
		fmt.Printf("result: %v, shared: %v\n", res.Val, res.Shared)
	case <-ctx.Done():
		fmt.Println("timeout")
		return
	}

	wg.Wait()
}

func forgetExample() {
    fmt.Println("\nForget Example")
	var wg sync.WaitGroup
	var g singleflight.Group

	// Function to simulate a long-running process
	fetchData := func(key string) (interface{}, error) {
		fmt.Printf("Fetching data for key: %s\n", key)
		time.Sleep(2 * time.Second) // Simulate a delay
		return fmt.Sprintf("Result for %s", key), nil
	}

    wg.Add(3)

	// Goroutine 1: Call fetchData with key "key1"
	go func() {
		result, err, shared := g.Do("key1", func() (interface{}, error) {
			return fetchData("key1")
		})
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Printf("G1 received: %v (shared: %v)\n", result, shared)
		}
        wg.Done()
	}()

	// Goroutine 2: Forget the key before the result is available
	go func() {
		time.Sleep(1 * time.Second) // Wait a bit to simulate concurrency
		g.Forget("key1")
		fmt.Println("G2 called Forget for key1")
        wg.Done()
	}()

	// Goroutine 3: Call fetchData with key "key1" after Forget
	go func() {
		time.Sleep(1 * time.Second) // Wait to ensure Forget is called first
		result, err, shared := g.Do("key1", func() (interface{}, error) {
			return fetchData("key1")
		})
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Printf("G3 received: %v (shared: %v)\n", result, shared)
		}
        wg.Done()
	}()

	wg.Wait()
}

func main() {
	SingleflightExample()
	goChanExample()
	forgetExample()
}
