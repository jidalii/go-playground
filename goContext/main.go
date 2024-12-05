package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// MultipleWatch()
	// WithCancelExample()
	// WithDeadlineExample()
	// WithTimeoutExample()
	// WithValueExample()
    MixedExample()
}

func MultipleWatch() {
	ctx, cancel := context.WithCancel(context.Background())
	go watch(ctx, "[Monitor1]")
	go watch(ctx, "[Monitor2]")
	go watch(ctx, "[Monitor3]")

	time.Sleep(3 * time.Second)
	fmt.Println("It is time to terminate all the monitors")
	cancel()
	time.Sleep(2 * time.Second)
}

func watch(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(name, "stopped monitoring at", time.Now().Format("15:04:05"))
			return
		default:
			fmt.Println(name, "monitoring at", time.Now().Format("15:04:05"))
			time.Sleep(1 * time.Second)
		}
	}
}

func WithCancelExample() {
	increment := func(ctx context.Context) <-chan int {
		dst := make(chan int)
		n := 0
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case dst <- n:
					n++
				}
			}
		}()
		return dst
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for n := range increment(ctx) {
		fmt.Println(n)
		if n == 5 {
			break
		}
	}
}

func WithDeadlineExample() {
	ddl := time.Now().Add(3 * time.Second)
	fmt.Println("Monitoring will be stopped at", ddl.Format("15:04:05"))
	ctx, cancel := context.WithDeadline(context.Background(), ddl)
	defer cancel()
	go watch(ctx, "Monitor1")
	time.Sleep(5 * time.Second)
}

func WithTimeoutExample() {
	timeout := 3 * time.Second
	fmt.Println("Monitoring will be stopped after", timeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	go watch(ctx, "Monitor1")
	time.Sleep(5 * time.Second)
}

func monitor(ctx context.Context, mName MonitorName) {
	name := ctx.Value(mName)
	for {
		select {
		case <-ctx.Done():
			fmt.Println(name, "stopped monitoring at", time.Now().Format("15:04:05"))
			return
		default:
			fmt.Println(name, "monitoring at", time.Now().Format("15:04:05"))
			time.Sleep(1 * time.Second)
		}
	}
}

type MonitorName string

func WithValueExample() {
	monitorName1 := MonitorName("MonitorKey1")
	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, monitorName1, "[Monitor1]")
	go monitor(ctx, monitorName1)
	time.Sleep(3 * time.Second)
	cancel()
}

func MixedExample() {
	timeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, os.Interrupt, syscall.SIGINT)

	go watch(ctx, "Monitor1")
	for {
		select {
		case <-cancelChan:
			cancel()
			fmt.Println("Cancel signal received and stop monitoring")
			return
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
