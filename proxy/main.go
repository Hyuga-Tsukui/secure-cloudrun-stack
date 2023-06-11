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

func main() {
	remote, err := url.Parse(os.Getenv("REMOTE_URL"))
	if err != nil {
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
