package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"gopkg.in/yaml.v2"
)

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
	var filename string
	flag.StringVar(&filename, "filename", "entries.yaml", "Name of the YAML file to read entries from.")
	flag.Parse()

	bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Println(err)
		return
	}

	entries := readEntriesFromYaml(bytes)
	mux := http.NewServeMux()
	makeHandlersFromMap(mux, entries)

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	s.ListenAndServe()
}
