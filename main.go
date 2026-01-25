package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
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

func handlePostTimestamp(reqCh chan<- request) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("Content-Type") != "text/plain" {
			http.Error(w, "unsupported content type", http.StatusUnsupportedMediaType)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusBadRequest)
			return
		}

		unix, err := strconv.ParseInt(string(body), 10, 64)
		if err != nil {
			http.Error(w, "invalid unix timestamp", http.StatusBadRequest)
			return
		}

		t := time.Unix(unix, 0)

		reqCh <- request{
			op:    opSet,
			value: t,
		}

		w.WriteHeader(http.StatusOK)
	}
}

func handleGetTimestamp(reqCh chan<- request) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		replyCh := make(chan time.Time)

		reqCh <- request{
			op:      opGet,
			replyCh: replyCh,
		}

		t := <-replyCh

		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, t.Unix())
	}
}

func main() {
	fmt.Println("Hello, World!")
	reqCh := make(chan request)
	go timestampOwner(reqCh)

	http.HandleFunc("/timestamp", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlePostTimestamp(reqCh)(w, r)
			return
		}
		if r.Method == http.MethodGet {
			handleGetTimestamp(reqCh)(w, r)
			return
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})
	fmt.Println("listening on :8080")
	http.ListenAndServe(":8080", nil)

}
