package main

import (
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

type yamlConfig struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	mapHandler := MapHandler(pathsToUrls, mux)

	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`

	yamlHandler, err := YAMLHandler([]byte(yaml), mapHandler)

	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		url, ok := pathsToUrls[r.URL.Path]

		if !ok {
			fallback.ServeHTTP(w, r)

			return
		}

		http.Redirect(w, r, url, http.StatusPermanentRedirect)
	}
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	yamlMappings := []yamlConfig{}

	err := yaml.Unmarshal(yml, &yamlMappings)

	return func(w http.ResponseWriter, r *http.Request) {

		for _, m := range yamlMappings {

			if m.Path == r.URL.Path {
				http.Redirect(w, r, m.URL, http.StatusPermanentRedirect)

				return
			}
		}

		fallback.ServeHTTP(w, r)
	}, err
}
