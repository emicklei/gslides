package main

import (
	"encoding/json"
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
	pageToken := ""
	docs := []Document{}

	for {
		// https://developers.google.com/drive/api/v3/search-files
		get, _ := http.NewRequest("GET", fmt.Sprint("https://www.googleapis.com/drive/v3/files"), nil)

		// add query params, https://developers.google.com/drive/api/v3/reference/files/list
		q := get.URL.Query()
		q.Add("spaces", "drive")
		q.Add("orderBy", "name")
		presentationsOnly := "mimeType='application/vnd.google-apps.presentation'"
		if owner := c.String("owner"); len(owner) > 0 {
			presentationsOnly = fmt.Sprintf("%s and '%s' in owners", presentationsOnly, owner)
		}
		q.Add("q", presentationsOnly)
		q.Add("fields", "nextPageToken,files(id,name)")
		if len(pageToken) > 0 {
			q.Add("pageToken", pageToken)
		}
		get.URL.RawQuery = q.Encode()

		// send it
		resp, err := client.Do(get)
		if err != nil {
			if resp != nil && resp.Body != nil {
				io.Copy(os.Stdout, resp.Body)
			}
			return fmt.Errorf("unable to list documents: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			io.Copy(os.Stdout, resp.Body)
			return fmt.Errorf("unable to list documents: %v", resp.Status)
		}
		result := new(DocumentList)
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return err
		}
		docs = append(docs, result.Files...)
		pageToken = result.NextPageToken
		if len(pageToken) == 0 {
			break
		}
	}
	for _, each := range docs {
		fmt.Println(each.ID, " : ", each.Name)
	}
	return nil
}

type DocumentList struct {
	NextPageToken string     `json:"nextPageToken"`
	Files         []Document `json:"files"`
}

type Document struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
