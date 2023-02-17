package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	reverseProxy := httputil.NewSingleHostReverseProxy(&url.URL{})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// !
		target := r.URL.Query().Get("target")
		targetUrl, err := url.Parse(target)

		// ?
		if target == "" || targetUrl == nil || err != nil {
			http.Error(w, "Target URL not specified or invalid", http.StatusBadRequest)
			return
		}

		// !
		reverseProxy.Director = func(req *http.Request) {
			req.URL.Scheme = targetUrl.Scheme
			req.URL.Host = targetUrl.Host
			req.URL.Path = targetUrl.Path
			req.Host = targetUrl.Host
			for header := range req.Header {
				if header != "Range" {
					if header != "range" {
						delete(req.Header, header)
					}
				}
			}
			req.Header.Set("Accept-Encoding", "gzip, deflate, br")
			req.Header.Set("accept", "*/*")
			req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		}

		reverseProxy.ServeHTTP(w, r)
	})
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
