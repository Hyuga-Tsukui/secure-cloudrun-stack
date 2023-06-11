package main

import (
	"net/http"

	"golang.org/x/exp/slog"
)

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	slog.Info("Server started")
	http.ListenAndServe(":8080", nil)
}
