package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
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
		w.Header().Set("Connection", "keep-alive")

		target := strings.Replace(r.URL.RawQuery, "q=", "", 1)
		targetUrl, err := url.Parse(target)

		if err != nil {
			http.Error(w, "Target URL invalid", http.StatusBadRequest)
			return
		}

		// Set the director function to modify the request
		reverseProxy.Director = func(req *http.Request) {
			req.URL = targetUrl
			for header := range req.Header {
				if header != "Range" && header != "range" {
					delete(req.Header, header)
				}
			}

			req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
			req.Header.Set("accept-language", "pt-BR,pt;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
			req.Header.Set("accept-encoding", "gzip, deflate, br")
			req.Header.Set("referrer", "rule34.xxx")
			req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36 Edg/110.0.1587.57")
			fmt.Println(targetUrl)
		}

		// CORS
		reverseProxy.ModifyResponse = func(res *http.Response) error {
			res.Header.Set("Access-Control-Allow-Origin", "*")
			res.Header.Set("Access-Control-Allow-Methods", "*")
			res.Header.Set("Access-Control-Allow-Headers", "*")
			res.Header.Set("Access-Control-Allow-Credentials", "true")
			res.Header.Set("Access-Control-Expose-Headers", "*")
			res.Header.Set("Connection", "keep-alive")
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
