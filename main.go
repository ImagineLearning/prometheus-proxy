package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	url, err := url.Parse(os.Getenv("PRP_URL"))
	if err != nil {
		log.Fatalf("Failed to parse URL, exiting: \"%s\"", url)
		os.Exit(1)
	}

	user := os.Getenv("PRP_USER")
	pass := os.Getenv("PRP_PASS")
	if user == "" || pass == "" {
		log.Fatal("Missing username or password")
		os.Exit(2)
	}

	proxy := &httputil.ReverseProxy{Director: buildAuthHandler(user, pass)}

	http.HandleFunc("/", buildProxyHandler(proxy))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func buildAuthHandler(user, pass string) func(*http.Request) {
	return func(req *http.Request) {
		req.SetBasicAuth(user, pass)
	}
}

func buildProxyHandler(p *httputil.ReverseProxy) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		p.ServeHTTP(w, r)
	}
}
