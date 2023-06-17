package main

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"golang.org/x/exp/slog"
	"google.golang.org/api/idtoken"
)

func initLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func main() {

	initLogger()

	rURL := os.Getenv("REMOTE_URL")
	remote, err := url.Parse(rURL)
	if err != nil {
		slog.Error("Failed to parse REMOTE_URL", "REMOTE_URL", rURL)
		panic(err)
	}

	slog.Info("Proxying to" + remote.String())
	handler := func(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := context.Background()
			ts, err := idtoken.NewTokenSource(ctx, remote.String())
			if err != nil {
				slog.Error("Failed to get token", err)
				return
			}
			token, err := ts.Token()
			if err != nil {
				slog.Error("Failed to get token", err)
				return
			}
			r.Header.Set("Authorization", "Bearer "+token.AccessToken)
			slog.Info(r.URL.String())
			r.Host = remote.Host
			p.ServeHTTP(w, r)
		}
	}

	slog.Info("Server started")

	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetXForwarded()
			r.SetURL(remote)
		},
	}

	http.HandleFunc("/", handler(proxy))
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error("Server failed to start", err)
	}
}
