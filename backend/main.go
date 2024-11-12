package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/go-github/github"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World!")
}

func someCode(w http.ResponseWriter, r *http.Request) {

	/*
			client := &http.Client{}

			// https://docs.github.com/en/rest/repos/contents?apiVersion=2022-11-28
			req, err := http.NewRequest("GET", "https://api.github.com/repos/davidthuman/recreation-Gov2Vec/contents/data.py", nil)
			if err != nil {
				panic(err)
			}

			req.Header.Set("Accept", "application/vnd.github+json")
		    //req.Header.Set("Authorization", "Bearer ")
			req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}

		    resp.Body.
	*/

	client := github.NewClient(nil)

	fileContent, _, _, err := client.Repositories.GetContents(context.TODO(), "davidthuman", "recreation-Gov2Vec", "data.py", nil)
	if err != nil {
		panic(err)
	}

	content, err := fileContent.GetContent()
	if err != nil {
		panic(err)
	}

	lines := strings.Split(content, "\n")
	someCode := strings.Join(lines[53:60], "\n")

	w.Header().Add("Access-Control-Allow-Origins", "*")

	fmt.Fprint(w, someCode)

}

func main() {

	http.HandleFunc("/", index)
	http.HandleFunc("/test", someCode)

	log.Printf("Serving on per 3333 ...")
	log.Fatal(http.ListenAndServe(":3333", nil))

}
