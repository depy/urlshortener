package main

import (
	"fmt"
	"net/http"
	"net/url"

	"gopkg.in/yaml.v2"
)

var data = `
hello: https://google.com
hi: https://cloudflare.com
heya: https://microsoft.com
yoho: https://cloud.google.com
bogus: notAurlHaha!
`

func readEntriesFromYaml(data []byte) map[string]string {
	t := make(map[string]string)
	yaml.Unmarshal(data, &t)
	return t
}

func makeHandlersFromMap(mux *http.ServeMux, entries map[string]string) {
	for shortcode, redirectUrl := range entries {
		_, err := url.ParseRequestURI(redirectUrl)

		if err != nil {
			fmt.Println("Skipping", redirectUrl, "=> Error: ", err)
			continue
		}

		path := "/" + shortcode
		mux.Handle(path, http.RedirectHandler(redirectUrl, http.StatusMovedPermanently))
	}
}

func main() {
	entries := readEntriesFromYaml([]byte(data))
	mux := http.NewServeMux()
	makeHandlersFromMap(mux, entries)

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	s.ListenAndServe()
}
