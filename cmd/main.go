package main

import (
	"fmt"
	"goconvert/internal/handlers"
	"net/http"
)

func main() {
	http.HandleFunc("/convert", corsMiddleware(handlers.UploadHandler))
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func corsMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		h(w, r)
	}
}
