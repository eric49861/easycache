package main

import (
	"cache"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"HELLO": "1",
	"WORLD": "2",
	"!":     "3",
}

func main() {
	cache.NewGroup("score", 1024, cache.GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exists", key)
	}))
	addr := "localhost:9090"
	pool := cache.NewHTTPPool(addr)
	log.Fatalln(http.ListenAndServe(addr, pool))
}
