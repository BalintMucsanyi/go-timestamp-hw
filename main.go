package main

import (
	"fmt"
	"time"
)

type operation int

const (
	opSet operation = iota
	opGet
)

type request struct {
	op      operation
	value   time.Time
	replyCh chan time.Time
}

func timestampOwner(reqCh <-chan request) {
	var stored time.Time
	for req := range reqCh {
		switch req.op {
		case opSet:
			stored = req.value
		case opGet:
			req.replyCh <- stored
		}
	}
}

func main() {
	fmt.Println("Hello, World!")
	reqCh := make(chan request)
	go timestampOwner(reqCh)

	select {}
}
