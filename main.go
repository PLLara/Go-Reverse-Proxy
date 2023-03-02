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

			//  :authority: rule34.xxx
			// :method: GET
			// :path: /public/autocomplete.php?q=la
			// :scheme: https
			// accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7
			// accept-encoding: gzip, deflate, br
			// accept-language: pt-BR,pt;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6
			// cache-control: max-age=0
			// cookie: __cf_bm=jvXeAywV6HJ5Fp.0MYjfuhk_nW8lxfREcoYlWqUUz7M-1677774008-0-AdiNgUu4eOgCMBov9uF59utDEsjAVId5kRO5zJWNUeuDCl7/I6y5JWTgEzW14mZA2Qv1a4krnDMUuRaVFojz7iRSBKa9Q5ecfdNiMYRgqVDRMSx6Dimf1bC79gpmhvsQhxSgZKZc9WzOMD2Mw+eB2Js=
			// sec-ch-ua: "Chromium";v="110", "Not A(Brand";v="24", "Microsoft Edge";v="110"
			// sec-ch-ua-mobile: ?0
			// sec-ch-ua-platform: "Windows"
			// sec-fetch-dest: document
			// sec-fetch-mode: navigate
			// sec-fetch-site: none
			// sec-fetch-user: ?1
			// upgrade-insecure-requests: 1
			// user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36 Edg/110.0.1587.57
			req.Header.Set("authority", "rule34.xxx")
			req.Header.Set("method", "GET")
			req.Header.Set("path", "/public/autocomplete.php?q=la")
			req.Header.Set("scheme", "https")
			req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
			req.Header.Set("accept-encoding", "gzip, deflate, br")
			req.Header.Set("accept-language", "pt-BR,pt;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
			req.Header.Set("cache-control", "max-age=0")
			req.Header.Set("cookie", "__cf_bm=jvXeAywV6HJ5Fp.0MYjfuhk_nW8lxfREcoYlWqUUz7M-1677774008-0-AdiNgUu4eOgCMBov9uF59utDEsjAVId5kRO5zJWNUeuDCl7/I6y5JWTgEzW14mZA2Qv1a4krnDMUuRaVFojz7iRSBKa9Q5ecfdNiMYRgqVDRMSx6Dimf1bC79gpmhvsQhxSgZKZc9WzOMD2Mw+eB2Js=")
			req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Microsoft Edge\";v=\"110\"")
			req.Header.Set("sec-ch-ua-mobile", "?0")
			req.Header.Set("sec-ch-ua-platform", "\"Windows\"")
			req.Header.Set("sec-fetch-dest", "document")
			req.Header.Set("sec-fetch-mode", "navigate")
			req.Header.Set("sec-fetch-site", "none")
			req.Header.Set("sec-fetch-user", "?1")
			req.Header.Set("upgrade-insecure-requests", "1")
			req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36 Edg/110.0.1587.57")

			// req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
			// req.Header.Set("accept-language", "pt-BR,pt;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
			// req.Header.Set("accept-encoding", "gzip, deflate, br")
			// req.Header.Set("referrer", "rule34.xxx")
			// req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36 Edg/110.0.1587.57")
			// fmt.Println(targetUrl)
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
