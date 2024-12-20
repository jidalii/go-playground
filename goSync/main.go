package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
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

//*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*//
//*                         WAITGROUP                          *//
//*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*//

func FindPrimeNumbersInRange(start, end int) {
	// Create a WaitGroup
	var wg sync.WaitGroup

	// Add the number of goroutines to the WaitGroup
	wg.Add(end - start + 1)

	// Create a goroutine for each number
	for i := start; i <= end; i++ {
		go IsPrime(i, &wg)
	}

	// Block here until all goroutines to finish
	wg.Wait()
}

func IsPrime(n int, wg *sync.WaitGroup) bool {
	// Defer Done to notify the WaitGroup when the task done
	defer wg.Done()

	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	fmt.Println("Prime number:", n)
	return true
}

//*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*//
//*                           ONCE                             *//
//*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*//

type DB struct{}

var db *DB
var once sync.Once

func GetDB() *DB {
	once.Do(func() {
		db = &DB{}
		fmt.Println("db instance created.")
	})
	return db
}

func (d *DB) Query() {
	fmt.Println("Querying the db.")
}

func SingletonOnce() {
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			GetDB().Query()
		}()
	}
	wg.Wait()
}

//*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*//
//*                           POOL                             *//
//*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*//

func PoolExample(workNum int) {
	// Counter for new buffer creations
	created := atomic.Int32{}

	bufferPool := sync.Pool{
		// New optionally specifies a function to generate
		// a value when Get would otherwise return nil.
		New: func() interface{} {
			created.Add(1)
			return new(bytes.Buffer)
		},
	}

	worker := func(id int, wg *sync.WaitGroup) {
		defer wg.Done()

		// Get buffer from pool
		buf := bufferPool.Get().(*bytes.Buffer)
		defer bufferPool.Put(buf)

		// Actually use the buffer
		buf.Reset()
		fmt.Fprintf(buf, "Worker %d using buffer\n", id)
		fmt.Print(buf.String())
	}

	wg := new(sync.WaitGroup)
	for i := 0; i < workNum; i++ {
		wg.Add(1)
		// Simulate request pattern: between 0ms - 10ms
		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
		go worker(i, wg)
	}
	wg.Wait()

	fmt.Printf("Created %d buffers\n", created.Load())
}

//*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*//
//*                           Cond                             *//
//*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*//

// ********** Producer-Consumer with Cond **********

func producer(ctx context.Context, id int, c *sync.Cond) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			c.L.Lock()
			// Produce something...
			fmt.Printf("%s- Producer %d: producing...\n", time.Now().Format("15:04:05.999"), id)
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			c.Signal()
			c.L.Unlock()
			time.Sleep(1 * time.Second)
		}
	}
}

func consumer(ctx context.Context, id int, c *sync.Cond) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			c.L.Lock()
			c.Wait()
			// Consume something...
			fmt.Printf("%s- Consumer %d: consuming...\n", time.Now().Format("15:04:05.999"), id)
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			c.L.Unlock()
		}
	}
}

func ProducerSingleConsumer() {
	mutex := sync.Mutex{}
	c := sync.NewCond(&mutex)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	go producer(ctx, 1, c)
	go consumer(ctx, 1, c)
	// go consumer(ctx, 2, c)
	fmt.Println("Producer and consumer started.")
	time.Sleep(5 * time.Second)
	fmt.Println("Producer and consumer ended.")
}

func ProducerMultiConsumers() {
	mutex := sync.Mutex{}
	c := sync.NewCond(&mutex)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go producer(ctx, 1, c)
	go consumer(ctx, 1, c)
	go consumer(ctx, 2, c)
	fmt.Println("Producer and consumer started.")
	time.Sleep(5 * time.Second)
	fmt.Println("Producer and consumer ended.")
}

// ********** Pub/Sub with Cond **********

func publisher(ctx context.Context, id int, c *sync.Cond) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			c.L.Lock()
			// Produce something...
			time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
			fmt.Printf("%s- Publisher %d: generating...\n", time.Now().Format("15:04:05.999"), id)
			c.Broadcast()
			c.L.Unlock()
			time.Sleep(1 * time.Second)
		}
	}
}

func subscriber(ctx context.Context, id int, c *sync.Cond) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			c.L.Lock()
			c.Wait()
			// Consume something...
			fmt.Printf("%s- Subscriber %d: reading...\n", time.Now().Format("15:04:05.999"), id)
			c.L.Unlock()
		}
	}
}

func PubMultiSub() {
	mutex := sync.Mutex{}
	c := sync.NewCond(&mutex)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go publisher(ctx, 1, c)
	go subscriber(ctx, 1, c)
	go subscriber(ctx, 2, c)
	go subscriber(ctx, 3, c)
	fmt.Println("Pub/Sub started.")
	time.Sleep(5 * time.Second)
	fmt.Println("Pub/Sub ended.")
}

func CondExampleWithPizzaRace() {
	mutex := sync.Mutex{}
	c := sync.NewCond(&mutex)
	wg := sync.WaitGroup{}
	pizzaRemaining := 10
	person := func(id int, c *sync.Cond, wg *sync.WaitGroup) {
		defer wg.Done()
		c.L.Lock()
		c.Wait() // Waits for broadcast
		c.L.Unlock()
		for {
			c.L.Lock()
			if pizzaRemaining <= 0 {
				c.L.Unlock()
				return
			}
			pizzaRemaining--
			fmt.Printf("Person %d got the pizza. Remaining: %d - %s\n",
				id, pizzaRemaining, time.Now().Format("15:04:05.999"))
			c.L.Unlock()

			// Simulate eating time
			time.Sleep(100 * time.Millisecond)
		}
	}

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go person(i, c, &wg)
	}
	time.Sleep(100 * time.Millisecond)

	c.Broadcast()

	wg.Wait()
}

func main() {
	// Mutex
	// EatPizzaMutex(10000)
	// EatPizzaRace(10000)
	// BankMutex()

	// RWMutex
	// SimpleRWMutex()
	// BankBalanceBenchmark()

	// WaitGroup
	// FindPrimeNumbersInRange(0, 100)

	// Once
	// SingletonOnce()

	// Pool
	// PoolExample(1000)

	// Cond
	ProducerSingleConsumer()
	ProducerMultiConsumers()
	PubMultiSub()
	CondExampleWithPizzaRace()
}
