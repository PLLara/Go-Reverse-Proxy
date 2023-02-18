package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	reverseProxy := httputil.NewSingleHostReverseProxy(&url.URL{})
	reverseProxy.Transport = &http.Transport{
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   100,
		MaxConnsPerHost:       100,
		IdleConnTimeout:       10 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	http.HandleFunc("/proxy", func(w http.ResponseWriter, r *http.Request) {
		// Get the target URL from the q parameter
		target := r.URL.Query().Get("q")
		targetUrl, err := url.Parse(target)

		if err != nil {
			http.Error(w, "Target URL invalid", http.StatusBadRequest)
			return
		}

		// Set the director function to modify the request
		reverseProxy.Director = func(req *http.Request) {
			req.URL = targetUrl
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

		reverseProxy.ModifyResponse = func(res *http.Response) error {
			res.Header.Set("Access-Control-Allow-Origin", "*")
			res.Header.Set("Access-Control-Allow-Methods", "*")
			res.Header.Set("Access-Control-Allow-Headers", "*")
			res.Header.Set("Access-Control-Allow-Credentials", "true")
			res.Header.Set("Access-Control-Expose-Headers", "*")
			// set infinite cache header
			res.Header.Set("Cache-Control", "max-age=31536000")
			res.Header.Set("Expires", "Thu, 31 Dec 2037 23:55:55 GMT")
			res.Header.Set("Pragma", "cache")

			return nil
		}

		reverseProxy.ServeHTTP(w, r)
	})
	if port == "" {
		port = "8080"
	}
	fmt.Println("Server started on http://localhost:" + port)
	// Start the server
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
