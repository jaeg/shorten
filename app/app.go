package app

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/google/logger"
	"github.com/gorilla/mux"
)

const AppName = "shorten"

var certFile = flag.String("cert-file", "", "location of cert file")
var keyFile = flag.String("key-file", "", "location of key file")
var port = flag.String("port", "8090", "port to host on")
var logPath = flag.String("log-path", "./logs.txt", "Logs location")
var baseUrl = flag.String("base-url", "localhost:8090", "Base url for shortened result")

type Response struct {
	Link  string `json:"link,omitempty"`
	Error string `json:"error,omitempty"`
}

type Request struct {
	Link string `json:"link"`
}

type App struct {
	links map[string]string
	srv   *http.Server
}

func (a *App) Init() {
	a.links = make(map[string]string)
	flag.Parse()
	//Start the logger
	lf, err := os.OpenFile(*logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}

	logger.Init(AppName, true, true, lf)

	logger.Infof("%s Starting", AppName)

	//Setup the HTTP server
	r := mux.NewRouter()
	r.HandleFunc("/shorten", a.handleShorten)
	r.PathPrefix("/").HandlerFunc(a.handleRedirect)
	r.HandleFunc("/liveness", a.handleLiveness)

	a.srv = &http.Server{
		Addr:    ":" + *port,
		Handler: r,
	}
}

func (a *App) Run(ctx context.Context) {
	defer logger.Close()

	//Run the http server
	go func() {
		if *certFile == "" || *keyFile == "" {
			logger.Info("Starting http")
			a.srv.ListenAndServe()
		} else {
			logger.Info("Starting https")
			a.srv.ListenAndServeTLS(*certFile, *keyFile)
		}
	}()

	// Handle shutdowns gracefully
	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := a.srv.Shutdown(ctxShutDown); err != nil {
		logger.Errorf("Failed to shutdown server %s", err.Error())
	} else {
		logger.Info("HTTP Server shutdown")
	}

}

func (a *App) handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			ErrorWithJson(w, "Error reading payload", http.StatusBadRequest)
			return
		}

		var req Request
		err = json.Unmarshal(body, &req)
		if err != nil {
			ErrorWithJson(w, "Error parsing payload", http.StatusBadRequest)
			return
		}

		if req.Link == "" {
			ErrorWithJson(w, "Invalid request", http.StatusBadRequest)

			return
		}

		shortPath := ""
		for shortPath = generateRandomName(5); a.links["/"+shortPath] != ""; shortPath = generateRandomName(5) {
		}
		a.links["/"+shortPath] = req.Link

		RespondWithJson(w, *baseUrl+"/"+shortPath)
	}
}

func (a *App) handleRedirect(w http.ResponseWriter, r *http.Request) {
	longURL := a.links[r.RequestURI]

	if longURL != "" {
		http.Redirect(w, r, longURL, http.StatusMovedPermanently)
	} else {
		ErrorWithJson(w, "Invalid link", http.StatusNotFound)
	}
}

func (a *App) handleLiveness(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "{}")
}

func generateRandomName(length int) (out string) {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	for i := 0; i < length; i++ {
		out += string(chars[rand.Intn(len(chars))])
	}

	return
}

func RespondWithJson(w http.ResponseWriter, link string) error {
	resp := &Response{Link: link}

	respJSON, err := json.Marshal(resp)

	if err != nil {
		logger.Errorf("Error when creating response json: %s", err.Error())
		http.Error(w, "{\"Error\":\"Serious error\"}", http.StatusInternalServerError)
		return err
	}

	w.Write(respJSON)

	return nil
}

func ErrorWithJson(w http.ResponseWriter, errorMsg string, status int) error {
	resp := &Response{Error: errorMsg}

	respJSON, err := json.Marshal(resp)

	if err != nil {
		logger.Errorf("Error when creating response json: %s", err.Error())
		http.Error(w, "{\"Error\":\"Serious error\"}", http.StatusInternalServerError)
		return err
	}

	http.Error(w, string(respJSON), status)

	return nil
}
