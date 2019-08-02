package main

import (
	"log"
	"math/rand"
	"net/http"
)

var links map[string]string

const baseURL = "localhost:9090/"

func main() {
	links = make(map[string]string)

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
		for shortPath = generateRandomName(5); links["/"+shortPath] != ""; shortPath = generateRandomName(5) {
		}
		links["/"+shortPath] = longURL
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
	longURL := links[r.RequestURI]
	if longURL != "" {
		http.Redirect(w, r, longURL, 301)
	} else {
		w.Write([]byte("Invalid Link"))
	}
}
