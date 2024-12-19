package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

//*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*//
//*                           MUTEX                            *//
//*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*//

func EatPizzaMutex(pizzaSlices int) {
	fmt.Println("********** EatPizzaMutex **********")
	total := 0
	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}
	ch := make(chan int, 3)
	defer close(ch)

	fmt.Printf("Total pizza slices: %d\n", pizzaSlices)

	person := func(name string, mutex *sync.Mutex, ch chan int) {
		defer wg.Done()
		cnt := 0
		for {
			mutex.Lock()
			if pizzaSlices > 0 {
				pizzaSlices--
				cnt++
				mutex.Unlock()
			} else {
				mutex.Unlock()
				fmt.Printf("%s took %d slices.\n", name, cnt)
				ch <- cnt
				break
			}
		}
	}

	wg.Add(3)
	go person("Alice", &mutex, ch)
	go person("Bob", &mutex, ch)
	go person("Charlie", &mutex, ch)

	for i := 0; i < 3; i++ {
		total += <-ch
	}

	wg.Wait()

	fmt.Printf("Total slices eaten: %d\n", total)
	fmt.Println("All pizza slices are gone!")
}

func EatPizzaRace(pizzaSlices int) {
	fmt.Println("********** EatPizzaRace **********")
	total := 0
	wg := sync.WaitGroup{}
	ch := make(chan int, 3)
	defer close(ch)

	fmt.Printf("Total pizza slices: %d\n", pizzaSlices)

	person := func(name string, ch chan int) {
		defer wg.Done()
		cnt := 0
		for {
			if pizzaSlices > 0 {
				pizzaSlices--
				cnt++
			} else {
				fmt.Printf("%s took %d slices.\n", name, cnt)
				ch <- cnt
				break
			}
		}
	}

	wg.Add(3)
	go person("Alice", ch)
	go person("Bob", ch)
	go person("Charlie", ch)

	for i := 0; i < 3; i++ {
		total += <-ch
	}

	wg.Wait()

	fmt.Printf("Total slices eaten: %d\n", total)
	fmt.Println("All pizza slices are gone!")
}

func BankMutex() {
	fmt.Println("********** BankMutex **********")
	balance := 1000
	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}
	ch := make(chan int, 3)
	defer close(ch)

	fmt.Printf("Initial balance: %d\n", balance)

	deposit := func(amount int, mutex *sync.Mutex, ch chan int) {
		fmt.Printf("Depositing %d\n", amount)
		defer wg.Done()
		mutex.Lock()
		balance += amount
		mutex.Unlock()
		ch <- amount
	}

	withdraw := func(amount int, mutex *sync.Mutex, ch chan int) {
		fmt.Printf("Withdrawing %d\n", amount)
		defer wg.Done()
		mutex.Lock()
		balance -= amount
		mutex.Unlock()
		ch <- -amount
	}

	wg.Add(3)
	go deposit(500, &mutex, ch)
	go withdraw(200, &mutex, ch)
	go withdraw(600, &mutex, ch)

	for i := 0; i < 3; i++ {
		fmt.Printf("Transaction amount: %d\n", <-ch)
	}

	wg.Wait()

	fmt.Printf("Final balance: %d\n", balance)
}

//*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*//
//*                          RWMUTEX                           *//
//*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*//

func SimpleRWMutex() {
	balance := 1000
	rwMutex := sync.RWMutex{}

	write := func(amount int, rwMutex *sync.RWMutex) {
		rwMutex.Lock()
		defer rwMutex.Unlock()
		balance += amount
	}
	read := func(rwMutex *sync.RWMutex) int {
		rwMutex.RLock()
		defer rwMutex.RUnlock()
		return balance
	}

	go write(100, &rwMutex)
	go read(&rwMutex)
	go write(200, &rwMutex)
	go read(&rwMutex)

	time.Sleep(2 * time.Second)
}

func BankBalanceBenchmark() {
	rwTotal := time.Duration(0)
	mTotal := time.Duration(0)
	for i := 0; i < 20; i++ {
		rwTotal += BankBalance(8000, 1000, true)
		mTotal += BankBalance(8000, 1000, false)
	}
	fmt.Printf("Average time using RWMutex: %s\n", rwTotal/10)
	fmt.Printf("Average time using Mutex:   %s\n", mTotal/10)
}

func BankBalance(rOp int, wOp int, useRWMutex bool) time.Duration {
    // fmt.Println("********** BankBalance **********")
    balance := 1000
    var mutex sync.Locker
    var rwMutex sync.RWMutex
    var simpleMutex sync.Mutex

    if useRWMutex {
        mutex = &rwMutex
    } else {
        mutex = &simpleMutex
    }

    wg := sync.WaitGroup{}
    ch := make(chan int, rOp+wOp)
    defer close(ch)

    updateBalance := func(mutex sync.Locker, ch chan int) {
        randOp := func() int {
            if rand.Intn(2) == 0 {
                return 1
            } else {
                return -1
            }
        }
        defer wg.Done()
        mutex.Lock()
        amount := randOp() * rand.Intn(balance)
        balance += amount
        mutex.Unlock()
        ch <- amount
    }

    readBalance := func(mutex sync.Locker, ch chan int) {
        defer wg.Done()
        if useRWMutex {
            rwMutex.RLock()
        } else {
            mutex.Lock()
        }
        ch <- balance
        if useRWMutex {
            rwMutex.RUnlock()
        } else {
            mutex.Unlock()
        }
    }

    startTime := time.Now()

    for i := 0; i < rOp; i++ {
        wg.Add(1)
        go readBalance(mutex, ch)
    }

    for i := 0; i < wOp; i++ {
        wg.Add(1)
        go updateBalance(mutex, ch)
    }

    wg.Wait()
	totalTime := time.Since(startTime)
	return totalTime
}

func main() {
	// Mutex
	// EatPizzaMutex(10000)
	// EatPizzaRace(10000)
	// BankMutex()

	// RWMutex
	SimpleRWMutex()

	BankBalanceBenchmark()
}
