package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"golang.org/x/exp/slog"
)

// https://gist.github.com/JalfResi/6287706

func main() {
	remote, err := url.Parse(os.Getenv("REMOTE_URL"))
	if err != nil {
		panic(err)
	}
	slog.Info("Proxying to" + remote.String())

	handler := func(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			slog.Info(r.URL.String())
			r.Host = remote.Host
			p.ServeHTTP(w, r)
		}
	}

	slog.Info("Server started")

	proxy := httputil.NewSingleHostReverseProxy(remote)
	http.HandleFunc("/", handler(proxy))
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error("Server failed to start", err)
	}
}
