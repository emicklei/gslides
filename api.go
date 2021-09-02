package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/slides/v1"
)

// TODO find a way to use Google Default Application Credentials

func getSlidesClient() (*slides.Service, *http.Client) {
	b := readClientID()

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/presentations")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client, _ := getClient(config)

	srv, err := slides.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Slides client: %v", err)
	}
	return srv, client
}

var configFilename = "gslides.json"

func readClientID() []byte {
	// local
	b, err := ioutil.ReadFile(configFilename)
	if err != nil {
		// home
		home := os.Getenv("HOME")
		b, err = ioutil.ReadFile(path.Join(home, configFilename))
		if err != nil {
			log.Fatalf("Unable to read client secret file (local or home): %v", err)
		}
	}
	return b
}
