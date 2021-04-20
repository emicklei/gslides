package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/emicklei/tre"
	"github.com/urfave/cli"
)

type thumbnail struct {
	ContentURL string `json:"contentUrl"`
}

func cmdExportThumbnails(c *cli.Context) error {
	srv, hc := getSlidesClient()
	presentationId := c.Args()[0]
	presentation, err := srv.Presentations.Get(presentationId).Do()
	if err != nil {
		return fmt.Errorf("Unable to retrieve data from presentation: %v", err)
	}
	for i, each := range presentation.Slides {
		// Get the thumbnail image
		resp, err := hc.Get(fmt.Sprintf("https://slides.googleapis.com/v1/presentations/%s/pages/%s/thumbnail", presentationId, each.ObjectId))
		if err != nil {
			return fmt.Errorf("Unable to retrieve slide thumbnail from presentation: %v", err)
		}
		data, _ := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		img := new(thumbnail)
		_ = json.Unmarshal(data, img)
		if err := exportImage(img.ContentURL, fmt.Sprintf("%s_slide_%d.png", presentationId, i+1)); err != nil {
			fmt.Println("slide export failed:" + err.Error())
		}
	}
	return nil
}

func exportImage(url string, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return tre.New(err, "unable to export image", "url", url)
	}
	defer resp.Body.Close()
	out, err := os.Create(filename)
	if err != nil {
		return tre.New(err, "unable to read image content", "url", url)
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return tre.New(err, "unable to write image content", "url", url)
}
