package main

import (
	"fmt"
	"goconvert/internal/handlers"
	"goconvert/internal/utils"
	"net/http"
)

func main() {
	port := utils.GetPort()
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Test server is running on port %s\n", port)
	})
	http.HandleFunc("/convert", corsMiddleware(handlers.UploadHandler))
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(port, nil)
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
