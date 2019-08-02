package main

import (
	"log"
	"math/rand"
	"net/http"

	"github.com/go-redis/redis"
)

const baseURL = "localhost:9090/"

var client *redis.Client

func main() {
	redisAddress := "localhost:6379"
	client = redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	http.HandleFunc("/", redirectHandler)
	http.HandleFunc("/shorten", shortenHandler)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		longURL := r.FormValue("url")
		if longURL == "" {
			w.Write([]byte("Invalid"))
			return
		}

		shortPath := ""
		for shortPath = generateRandomName(5); client.Exists("/"+shortPath).Val() == 1; shortPath = generateRandomName(5) {
		}
		client.Set("/"+shortPath, longURL, -1)
		w.Write([]byte(baseURL + shortPath))
	}
}

func generateRandomName(length int) (out string) {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	for i := 0; i < length; i++ {
		out += string(chars[rand.Intn(len(chars))])
	}

	return
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	longURL := client.Get(r.RequestURI).String()
	if longURL != "" {
		http.Redirect(w, r, longURL, 301)
	} else {
		w.Write([]byte("Invalid Link"))
	}
}
