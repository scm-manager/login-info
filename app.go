package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
)

type Configuration map[string]LocalizedLoginInfo

type LocalizedLoginInfo struct {
	Plugins  []InfoItem
	Features []InfoItem
}

type InfoItem struct {
	Title   string
	Summary string
	Link    string
}

type Links map[string]Link

type Link struct {
	Href string `json:"href"`
}

type ResponseItem struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Links   Links  `json:"_links"`
}

type Response struct {
	Plugin  ResponseItem `json:"plugin"`
	Feature ResponseItem `json:"feature"`
	Links   Links        `json:"_links"`
}

func main() {
	configuration := readConfiguration()

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/live", NewOkHandler())
	r.HandleFunc("/api/v1/ready", NewOkHandler())
	r.HandleFunc("/api/v1/login-info", NewLoginInfoHandler(configuration))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Println("start login-info on port", port)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal("http server returned err: ", err)
	}
}

func NewOkHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}
}

func NewLoginInfoHandler(configuration Configuration) http.HandlerFunc {
	localeMap := map[language.Tag]LocalizedLoginInfo{}

	var languages []language.Tag
	for langKey, item := range configuration {
		t := language.MustParse(langKey)
		log.Printf("append %v to list of languages", t)
		languages = append(languages, t)
		localeMap[t] = item
	}

	matcher := language.NewMatcher(languages)

	return func(w http.ResponseWriter, r *http.Request) {
		accept := r.Header.Get("Accept-Language")
		_, index := language.MatchStrings(matcher, accept)

		tag := languages[index]

		info, ok := localeMap[tag]
		if !ok {
			http.Error(w, "failed to find response language", 500)
			return
		}

		response := Response{
			Feature: mapItem(pickRandom(info.Features)),
			Plugin:  mapItem(pickRandom(info.Plugins)),
			Links:   createLinksWithSelf(r.RequestURI),
		}

		data, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "failed to marshal response", 500)
		}

		w.WriteHeader(200)
		w.Header().Add("Content-Type", "application/json")
		_, err = w.Write(data)
		if err != nil {
			log.Println("failed to write response", err)
		}
	}
}

func mapItem(item InfoItem) ResponseItem {
	return ResponseItem{
		Title:   item.Title,
		Summary: item.Summary,
		Links:   createLinksWithSelf(item.Link),
	}
}

func createLinksWithSelf(href string) Links {
	links := map[string]Link{}
	links["self"] = Link{
		Href: href,
	}
	return links
}

func pickRandom(items []InfoItem) InfoItem {
	return items[rand.Intn(len(items))]
}

func readConfiguration() Configuration {
	configPath := os.Getenv("CONFIG")
	if configPath == "" {
		configPath = "config.yaml"
	}
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("failed to read configuration %s: %v", configPath, err)
	}
	config := Configuration{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("failed to unmarshal configuration %s: %v", configPath, err)
	}

	return config
}
