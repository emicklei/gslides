package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/google/uuid"
	"github.com/urfave/cli"
	"google.golang.org/api/slides/v1"
)

// Known issue: https://issuetracker.google.com/issues/36761705?pli=1

// target << source[index]
func cmdAppendSlide(c *cli.Context) error {
	srv, _ := getSlidesClient()
	presentationTarget, err := srv.Presentations.Get(c.Args()[0]).Do()
	if err != nil {
		return fmt.Errorf("Unable to retrieve data from target presentation: %v", err)
	}
	presentationSource, err := srv.Presentations.Get(c.Args()[1]).Do()
	if err != nil {
		return fmt.Errorf("Unable to retrieve data from source presentation: %v", err)
	}
	sourceSlideIndex, err := strconv.Atoi(c.Args()[2])
	if err != nil {
		return fmt.Errorf("Invalid slide presentation index: %v", err)
	}
	if sourceSlideIndex >= len(presentationSource.Slides) {
		return fmt.Errorf("No such slide index: %v", sourceSlideIndex)
	}
	sourceSlideIndex-- // zero indexed

	// collect all changes
	batchReq := new(slides.BatchUpdatePresentationRequest)

	// new slide
	newSlideId := uuid.New().String()
	sourceSlide := presentationSource.Slides[sourceSlideIndex]
	addSlide := &slides.CreateSlideRequest{
		ObjectId:             newSlideId,
		SlideLayoutReference: &slides.LayoutReference{PredefinedLayout: (layoutNameWithID(presentationSource, sourceSlide.SlideProperties.LayoutObjectId))},
	}
	batchReq.Requests = append(batchReq.Requests, &slides.Request{CreateSlide: addSlide})

	if c.GlobalBool("v") {
		for i, each := range presentationTarget.Masters {
			log.Println("master", i, each.ObjectId, each.PageType)
		}
	}

	// updateSlideProps := &slides.UpdateSlidePropertiesRequest{
	// 	ObjectId: newSlideId,
	// 	SlideProperties: &slides.SlideProperties{
	// 		MasterObjectId: presentationTarget.Masters[0].ObjectId,
	// 	},
	// 	Fields: "masterObjectId", // cannot be updated :-(
	// }
	// batchReq.Requests = append(batchReq.Requests, &slides.Request{UpdateSlideProperties: updateSlideProps})

	// copy all elements
	for _, each := range sourceSlide.PageElements {
		if each.Shape != nil {
			copyShapeOfElement(c, each, newSlideId, batchReq)
		}
		if each.ElementGroup != nil {
			todo("slide.pagelement.ElementGroup")
		}
		if each.Image != nil {
			todo("slide.pagelement.Image")
		}
		if each.Line != nil {
			todo("slide.pagelement.Line")
		}
		if each.Table != nil {
			todo("slide.pagelement.Table")
		}
		if each.SheetsChart != nil {
			todo("slide.pagelement.SheetsChart")
		}
		if each.Video != nil {
			todo("slide.pagelement.Video")
		}
		if each.WordArt != nil {
			todo("slide.pagelement.WordArt")
		}
	}
	if c.GlobalBool("v") {
		log.Println("batch requests:", len(batchReq.Requests))
	}
	_, err = srv.Presentations.BatchUpdate(presentationTarget.PresentationId, batchReq).Do()
	if err != nil {
		return fmt.Errorf("Unable to send batch update to presentation: %v", err)
	}
	return nil
}

func todo(path string) {
	log.Println("TODO:", path)
}

func layoutNameWithID(p *slides.Presentation, id string) string {
	for _, each := range p.Layouts {
		if each.ObjectId == id {
			return each.LayoutProperties.Name
		}
	}
	return "?"
}

func copyShapeOfElement(c *cli.Context, elem *slides.PageElement, newSlideId string, batch *slides.BatchUpdatePresentationRequest) {
	props := new(slides.PageElementProperties)
	props.PageObjectId = newSlideId
	props.Size = elem.Size
	props.Transform = elem.Transform
	shapeId := uuid.New().String()
	req := &slides.CreateShapeRequest{
		ObjectId:          shapeId,
		ElementProperties: props,
		ShapeType:         elem.Shape.ShapeType,
	}
	if c.GlobalBool("v") {
		log.Println("add create shape:", elem.Shape.ShapeType)
	}
	batch.Requests = append(batch.Requests, &slides.Request{CreateShape: req})

	if elem.Shape.ShapeType == "TEXT_BOX" {
		copyTextBox(c, elem.Shape, shapeId, batch)
		return
	}
	todo(elem.Shape.ShapeType)
}

func copyTextBox(c *cli.Context, src *slides.Shape, shapeId string, batch *slides.BatchUpdatePresentationRequest) {
	for _, te := range src.Text.TextElements {
		if te.TextRun != nil {
			//add insert text
			insertText := &slides.InsertTextRequest{
				ObjectId: shapeId,
				Text:     te.TextRun.Content,
			}
			if c.GlobalBool("v") {
				log.Println("add insert text:", te.TextRun.Content)
			}
			batch.Requests = append(batch.Requests, &slides.Request{InsertText: insertText})

			// style
			updateStyle := &slides.UpdateTextStyleRequest{
				ObjectId: shapeId,
				Style:    te.TextRun.Style,
				// TextRange: nil,
				Fields: "*",
			}
			if c.GlobalBool("v") {
				log.Printf("set text style:%+v\n", te.TextRun.Style)
			}
			batch.Requests = append(batch.Requests, &slides.Request{UpdateTextStyle: updateStyle})
		}
	}
}
