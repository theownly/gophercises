package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type Arc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type Story struct {
	Arcs map[string]Arc `json:"arcs"`
}

type storyHandler struct {
	S Story
}

func (h storyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := strings.TrimPrefix(r.URL.Path, "/")

	if url == "" {

		url = "intro"

	}

	tmpl := template.Must(template.ParseFiles("template.html"))
	tmpl.Execute(w, h.S.Arcs[url])
}

func main() {
	jsonFile, err := os.Open("gopher.json")

	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	byteValues, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		log.Fatal(err)
	}

	var St Story

	json.Unmarshal(byteValues, &St.Arcs)

	handler := storyHandler{S: St}

	fmt.Println("Starting Server at Port 8080...")

	http.ListenAndServe(":8080", handler)
}
