package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

type arc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	}
}

func main() {

	// ReadFile is a wrapper around os.Open() and .readAll()
	b, err := ioutil.ReadFile("./story.json")

	if err != nil {
		panic(err)
	}

	storyArcs := make(map[string]arc)

	// since we already have the json in memory, Unmarshal
	// is a valid option over json.NewDecoder().Decode()
	// see: https://stackoverflow.com/questions/21197239/decoding-json-using-json-unmarshal-vs-json-newdecoder-decode
	// see: https://blog.golang.org/json
	err = json.Unmarshal(b, &storyArcs)

	if err != nil {
		panic(err)
	}

	tpls := template.Must(template.New("").ParseGlob("templates/*.gohtml"))

	// ignore the request for the favicon
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})

	http.HandleFunc("/not-found", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		tpls.ExecuteTemplate(w, "404.gohtml", nil)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		arcName := strings.Trim(r.URL.Path, "/")

		// Default to the /intro route
		// since that is the beginning of our story
		if arcName == "" {
			arcName = "intro"
		}

		arcDetails, ok := storyArcs[arcName]

		if !ok {
			http.Redirect(w, r, "/not-found", http.StatusSeeOther)

			return
		}

		tpls.ExecuteTemplate(w, "arc.gohtml", arc{arcDetails.Title, arcDetails.Story, arcDetails.Options})
	})

	http.ListenAndServe(":8080", nil)
}
