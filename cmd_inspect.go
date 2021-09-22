package main

import (
	"fmt"
	"strconv"

	"github.com/urfave/cli/v2"
	"google.golang.org/api/slides/v1"
)

// go run *.go -v inspect 1q9VqtPPwyGre9-o3uzlu_u7AEJh-jVFlkB02wJfi4EA 13
func cmdInspect(c *cli.Context) error {
	srv, _ := getSlidesClient()
	presentationTarget, err := srv.Presentations.Get(c.Args().First()).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve data from target presentation: %v", err)
	}
	if c.Args().Len() == 2 {
		i, err := strconv.Atoi(c.Args().Get(1))
		if err != nil {
			return fmt.Errorf("unable to convert slide index: %v", err)
		}
		if i <= 0 || i > len(presentationTarget.Slides) {
			return fmt.Errorf("no such slide index: %v", err)
		}
		reportSlide(presentationTarget.Slides[i-1])
	} else {
		reportPresentation(presentationTarget)
	}
	return nil
}

func reportSlide(s *slides.Page) {

}
func reportPresentation(p *slides.Presentation) {
	fmt.Println("masters:", len(p.Masters))
	for _, each := range p.Masters {
		fmt.Println("master:", each.ObjectId, "name:", each.MasterProperties.DisplayName, "pages:", len(each.PageElements))
	}
	fmt.Println("layouts:", len(p.Layouts))
	for _, each := range p.Layouts {
		fmt.Println("layout:", each.ObjectId, "name:", each.LayoutProperties.Name, "pages:", len(each.PageElements))
	}
}
