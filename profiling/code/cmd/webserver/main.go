package main

import _ "net/http/pprof"

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type InMemoryCache struct {
	mu    sync.Mutex
	store map[string][]byte
}

func (c *InMemoryCache) Set(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = value
}

func (c *InMemoryCache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, ok := c.store[key]
	return val, ok
}

func NewInMemoryCache() InMemoryCache {
	return InMemoryCache{store: make(map[string][]byte)}
}

var cache = NewInMemoryCache()

func leakyHandler(w http.ResponseWriter, r *http.Request) {
	// Memory leak: continuously adding to the cache without removing old entries
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key-%d", time.Now().UnixNano())
		value := make([]byte, 1024) // 1KB per entry
		cache.Set(key, value)
	}
	fmt.Fprintf(w, "Leaky handler called")
}

func lockingHandler(w http.ResponseWriter, r *http.Request) {
	cache.mu.Lock()
	fmt.Fprintf(w, "Locking handler called")
}

func unlockingHandler(w http.ResponseWriter, r *http.Request) {
	cache.mu.Unlock()
	fmt.Fprintf(w, "Unlocking handler called")
}

func main() {
	go func() {
		log.Println("Starting pprof server on :6060")
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/leak", leakyHandler)
	serveMux.HandleFunc("/clear", func(w http.ResponseWriter, r *http.Request) {
		cache = NewInMemoryCache()
		fmt.Fprintf(w, "Cache cleared")
	})
	serveMux.HandleFunc("/lock", lockingHandler)
	serveMux.HandleFunc("/unlock", unlockingHandler)
	serveMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", serveMux))
}
