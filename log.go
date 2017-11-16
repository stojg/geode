package main

import (
	"fmt"
	"os"
	"time"
)

func NewLogger(file string) *Logger {
	l := &Logger{}
	f, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	l.file = f
	l.Println("Program started")
	return l
}

type Logger struct {
	file *os.File
}

func (l *Logger) Close() error {
	l.Println("Program stopped")
	return l.file.Close()
}

func (l *Logger) Println(a ...interface{}) {
	args := append([]interface{}{l.ts()}, a...)
	_, err := fmt.Fprintln(l.file, args...)
	if err != nil {
		panic(err)
	}
}

func (l *Logger) Printf(format string, a ...interface{}) {
	args := append([]interface{}{l.ts()}, a...)
	_, err := fmt.Fprintf(l.file, "%s "+format, args...)
	if err != nil {
		panic(err)
	}
}

func (l *Logger) Error(inError error) {
	_, err := fmt.Fprintf(l.file, "%s %v\n", l.ts(), inError)
	fmt.Fprintf(os.Stderr, "%s %v\n", l.ts(), inError)
	if err != nil {
		panic(err)
	}
}

func (l *Logger) ts() string {
	return time.Now().Format("15:04:05.000000000")
}
