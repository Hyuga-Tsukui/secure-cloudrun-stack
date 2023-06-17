package main

import (
	"net/http"
	"os"

	"golang.org/x/exp/slog"
)

func initLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func main() {
	initLogger()

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	slog.Info("Server started", "port", "8080")
	http.ListenAndServe(":8080", nil)
}
