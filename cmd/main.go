package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var keys []string
var mu sync.Mutex

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	if r.Method == "POST" {
		var key struct {
			Key    string `json:"key"`
			ticker *time.Ticker
		}
		if err := json.NewDecoder(r.Body).Decode(&key); err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		key.ticker = time.NewTicker(1 * time.Minute)

		go func() {
			<-key.ticker.C
			mu.Lock()
			defer mu.Unlock()
			for i, k := range keys {
				if k == key.Key {
					keys = append(keys[:i], keys[i+1:]...)
					break
				}
			}
		}()

		keys = append(keys, key.Key)
		w.WriteHeader(http.StatusOK)
	} else {
		data, err := json.Marshal(keys)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(data))
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ok")
	fmt.Fprint(w, http.StatusText(http.StatusOK))
}

func main() {
	fmt.Println("Server started")
	http.Handle("/", http.FileServer(http.Dir("./static/")))
	http.Handle("/healthCheck", http.HandlerFunc(healthCheck))
	http.Handle("/api/", withCORS(http.HandlerFunc(handler)))
	fmt.Println(http.ListenAndServe(":8011", nil))
}

func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}
