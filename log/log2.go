//go:build ignore
// +build ignore

package main

import (
	"log"
	"os"
	"time"
)

var logFileName = "./output2.log"
var (
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	DebugLogger *log.Logger
	ErrorLogger *log.Logger
	FatalLogger *log.Logger
)

func init() {
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Unable to create Logger file:", err.Error())
		return
	}
	log.SetOutput(logFile)
	logFlag := log.Ldate | log.Ltime | log.Lshortfile
	InfoLogger = log.New(logFile, "Info:", logFlag)
	WarnLogger = log.New(logFile, "Warn:", logFlag)
	DebugLogger = log.New(logFile, "Debug:", logFlag)
	ErrorLogger = log.New(logFile, "Error:", logFlag)
	FatalLogger = log.New(logFile, "Fatal:", logFlag)
}

func main() {
	InfoLogger.Println("Starting Service!")
	t := time.Now()
	InfoLogger.Printf("Time taken: %s \n", time.Since(t))
	DebugLogger.Println("Debug: Service is running!")
	WarnLogger.Println("Warning: Service is about to shutdown!")
	ErrorLogger.Println("Error: Service is shutting down!")
	FatalLogger.Println("Fatal: Service shut down!")
}
