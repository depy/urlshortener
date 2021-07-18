package main

import (
	"fmt"
	"net/http"
	"net/url"
)

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
	entries := make(map[string]string)
	entries["hello"] = "https://google.com"
	entries["hi"] = "https://cloudflare.com"
	entries["heyho"] = "sfg3gGgj;"

	mux := http.NewServeMux()
	makeHandlersFromMap(mux, entries)

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	s.ListenAndServe()
}
