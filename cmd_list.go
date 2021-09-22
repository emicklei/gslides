package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2/google"
)

func cmdList(c *cli.Context) error {
	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(readClientID(), "https://www.googleapis.com/auth/drive")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client, _ := getClient(config)
	// https://developers.google.com/drive/api/v3/search-files
	get, err := http.NewRequest("GET", fmt.Sprint("https://www.googleapis.com/drive/v3/files?q=mimeType%3D'application/vnd.google-apps.presentation'&fields=files(id,name)"), nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(get)
	if err != nil {
		io.Copy(os.Stdout, resp.Body)
		return fmt.Errorf("unable to list documents: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		io.Copy(os.Stdout, resp.Body)
		return fmt.Errorf("unable to list documents: %v", resp.Status)
	}
	io.Copy(os.Stdout, resp.Body)
	return nil
}
