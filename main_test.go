package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func TestStoreAndFetchTimestamp(t *testing.T) {
	server := httptest.NewServer(newServer())
	defer server.Close()

	// POST timestamp
	resp, err := http.Post(
		server.URL+"/timestamp",
		"text/plain",
		strings.NewReader("1700000000"),
	)
	if err != nil {
		t.Fatalf("POST failed: %v", err)
	}
	resp.Body.Close()

	// GET timestamp
	resp, err = http.Get(server.URL + "/timestamp")
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "1700000000" {
		t.Fatalf("expected 1700000000, got %s", body)
	}
}

func TestConcurrentAccess(t *testing.T) {
	server := httptest.NewServer(newServer())
	defer server.Close()

	const workers = 50
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		i := i // capture loop variable
		go func() {
			defer wg.Done()

			ts := strconv.Itoa(1700000000 + i)

			// POST
			resp, err := http.Post(
				server.URL+"/timestamp",
				"text/plain",
				strings.NewReader(ts),
			)
			if err != nil {
				t.Errorf("POST failed: %v", err)
				return
			}
			resp.Body.Close()

			// GET
			resp, err = http.Get(server.URL + "/timestamp")
			if err != nil {
				t.Errorf("GET failed: %v", err)
				return
			}
			resp.Body.Close()
		}()
	}

	wg.Wait()
}
