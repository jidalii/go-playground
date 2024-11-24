//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"log"
	"os"
	"io"
	_"bytes"
)

type User struct {
	Name string
	Age  int
}

var logFileName = "./output1.log"
var newLogger *log.Logger

func init() {
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}

	log.SetPrefix("[app]: ")
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	log.SetOutput(logFile)

	newLogger = log.New(logFile, "[new]: ", log.Lshortfile|log.Ldate|log.Lmicroseconds)
}

func main() {
	u := User{
		Name: "Alice",
		Age:  18,
	}

	log.Printf("%s login, age:%d", u.Name, u.Age)
	log.Printf("Hello, %s", u.Name)
	// log.Panicf("Oh, system error when %s login", u.Name)
	// log.Fatalf("Danger! hacker %s login", u.Name)
	
	newLogger.Printf("Hello, %s", u.Name)

	// writer1 := &bytes.Buffer{}
	writer1 := os.Stdout
	writer2, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("create file log.txt failed: %v", err)
	}

	logger := log.New(io.MultiWriter(writer1, writer2), "[logger]: ", log.Lshortfile|log.LstdFlags)
	logger.Printf("%s login, age:%d", u.Name, u.Age)
}
