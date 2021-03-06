package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	base, err := url.Parse(os.Getenv("PP_URL"))
	if err != nil {
		log.Fatalf("Failed to parse URL, exiting: \"%s\"", base)
		os.Exit(1)
	}
	if base.Host == "" {
		log.Fatal("URL is blank, exiting")
		os.Exit(2)
	}

	user := os.Getenv("PP_USER")
	pass := os.Getenv("PP_PASS")
	if user == "" || pass == "" {
		log.Fatal("Missing username or password")
		os.Exit(3)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		reqUrl := url.URL{
			Scheme:   base.Scheme,
			Host:     base.Host,
			Path:     r.URL.Path,
			RawQuery: r.URL.RawQuery,
		}
		req, _ := http.NewRequest(r.Method, reqUrl.String(), nil)
		req.SetBasicAuth(user, pass)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Failed to %s %s: %s", req.Method, req.URL.String(), err.Error())
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		defer res.Body.Close()

		body, _ := ioutil.ReadAll(res.Body)

		fmt.Fprintln(w, string(body))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
