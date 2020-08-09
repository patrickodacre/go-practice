package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"
)

// A map containing ALL redirects regardless of their source config, eg: yaml, json, etc.
var redirects *map[string]string

func main() {

	var (
		yamlFile = flag.String("yaml-redirects", "./yaml-redirects.yaml", "Specify a local yaml file with the desired redirects.")
		jsonFile = flag.String("json-redirects", "./json-redirects.json", "Specify a local json file with the desired redirects.")
	)

	flag.Parse()

	// default mux will be our final fallback
	mux := defaultMux()

	redirects = &map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	// yaml file
	yamlRedirects, err := ioutil.ReadFile(*yamlFile)

	if err != nil {
		panic(err)
	}

	jsonRedirects, err := ioutil.ReadFile(*jsonFile)

	if err != nil {
		panic(err)
	}

	addYAMLRedirects([]byte(yamlRedirects), redirects)
	addJSONRedirects([]byte(jsonRedirects), redirects)

	handler := createHandler(redirects, mux)

	fmt.Println("Starting the server on :8080")
	fmt.Println("with the following redirects configured:")

	for path, url := range *redirects {
		fmt.Println("Path", path)
		fmt.Println("to URL", url)
	}

	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

// createHandler replaces MapHandler.
// There is no need for multiple handlers
// when all we need to do is parse the various
// configurations (yaml, json, etc) and add those
// redirects to our ultimate redirect map.
func createHandler(redirects *map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		url, ok := (*redirects)[r.URL.Path]

		if !ok {
			fallback.ServeHTTP(w, r)

			return
		}

		http.Redirect(w, r, url, http.StatusPermanentRedirect)
	}
}

// addYAMLRedirects appends the YAML-configured redirects to the master redirect map.
func addYAMLRedirects(yml []byte, redirects *map[string]string) {

	// Just inline the struct.
	// No need for a named type
	// before we need this struct elsewhere.
	yamlMappings := []struct {
		Path string `yaml:"path"`
		URL  string `yaml:"url"`
	}{}

	err := yaml.Unmarshal(yml, &yamlMappings)

	// if we can't parse the redirects,
	// we shouldn't continue.
	if err != nil {
		panic(err)
	}

	for _, m := range yamlMappings {
		(*redirects)[m.Path] = m.URL
	}
}

// addJSONRedirects appends the YAML-configured redirects to the master redirect map.
func addJSONRedirects(jsonRedirects []byte, redirects *map[string]string) {

	// Just inline the struct.
	// No need for a named type
	// before we need this struct elsewhere.
	jsonMappings := []struct {
		Path string `json:"path"`
		URL  string `json:"url"`
	}{}

	err := json.Unmarshal(jsonRedirects, &jsonMappings)

	// if we can't parse the redirects,
	// we shouldn't continue.
	if err != nil {
		panic(err)
	}

	for _, m := range jsonMappings {
		(*redirects)[m.Path] = m.URL
	}
}
