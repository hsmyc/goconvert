package main

import (
	"fmt"
	"goconvert/internal/handlers"
	"goconvert/internal/utils"
	"net/http"
)

func main() {
	port := utils.GetPort()
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
