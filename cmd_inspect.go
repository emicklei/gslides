package main

import (
	"fmt"
	"strconv"

	"github.com/skinass/go-spew/spew"
	"github.com/urfave/cli"
)

// go run *.go -v inspect 1q9VqtPPwyGre9-o3uzlu_u7AEJh-jVFlkB02wJfi4EA 13
func cmdInspect(c *cli.Context) error {
	spew.Config.DisableNilValues = true
	spew.Config.DisableZeroValues = true
	srv, _ := getSlidesClient()
	presentationTarget, err := srv.Presentations.Get(c.Args()[0]).Do()
	if err != nil {
		return fmt.Errorf("Unable to retrieve data from target presentation: %v", err)
	}
	if len(c.Args()) == 2 {
		i, err := strconv.Atoi(c.Args()[1])
		if err != nil {
			return fmt.Errorf("Unable to convert slide index: %v", err)
		}
		if i <= 0 || i > len(presentationTarget.Slides) {
			return fmt.Errorf("No such slide index: %v", err)
		}
		spew.Dump(presentationTarget.Slides[i-1])
	} else {
		spew.Dump(presentationTarget)
	}
	return nil
}
