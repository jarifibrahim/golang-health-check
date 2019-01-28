package main

import (
	"fmt"
	"log"
	"net/http"
)

var isOnline = true

func main() {
	http.Handle("/api/greeting", loggingMiddleware(http.HandlerFunc(greeting)))
	http.Handle("/api/stop", loggingMiddleware(http.HandlerFunc(stop)))
	http.Handle("/api/health", loggingMiddleware(http.HandlerFunc(health)))
	http.Handle("/", loggingMiddleware(http.FileServer(assetFS())))
	fmt.Println("Web server running on port 8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func greeting(w http.ResponseWriter, r *http.Request) {
	if isOnline {
		message := "World"
		if m := r.FormValue("name"); m != "" {
			message = m
		}
		fmt.Fprintf(w, "Hello %s!", message)
		return
	}
	w.WriteHeader(503)
	w.Write([]byte("Not Online"))
}

func stop(w http.ResponseWriter, r *http.Request) {
	isOnline = false
	w.Write([]byte("Stopping HTTP Server"))
}

func health(w http.ResponseWriter, r *http.Request) {
	if isOnline {
		w.WriteHeader(200)
		return
	}
	w.WriteHeader(500)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("HEADER:%s - METHOD:%s PATH:%s", r.Header, r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
