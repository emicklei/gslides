package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/slides/v1"
)

// TODO find a way to use Google Default Application Credentials

func getSlidesClient() (*slides.Service, *http.Client) {
	b, err := ioutil.ReadFile("gslides.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/presentations")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client, _ := getClient(config)

	srv, err := slides.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Slides client: %v", err)
	}
	return srv, client
}
