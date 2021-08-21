package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/urfave/cli"
	"google.golang.org/api/slides/v1"
)

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

	var pageObjectId string
	{
		batchReq := new(slides.BatchUpdatePresentationRequest)
		addSlide := &slides.Request{CreateSlide: new(slides.CreateSlideRequest)}
		batchReq.Requests = append(batchReq.Requests, addSlide)
		batchResp, err := srv.Presentations.BatchUpdate(presentationTarget.PresentationId, batchReq).Do()
		if err != nil {
			return fmt.Errorf("Unable to send batch update to presentation: %v", err)
		}
		log.Println(batchResp)
		pageObjectId = batchResp.Replies[0].CreateSlide.ObjectId
	}

	{
		batchReq := new(slides.BatchUpdatePresentationRequest)
		sourceSlide := presentationSource.Slides[sourceSlideIndex]
		for _, each := range sourceSlide.PageElements {
			fmt.Println("----------------------------------------")
			//spew.Dump(each)
			if each.Shape != nil {
				spew.Dump(each.Shape)
				props := new(slides.PageElementProperties)
				props.PageObjectId = pageObjectId
				props.Size = each.Size
				props.Transform = each.Transform
				shapeId := uuid.New().String()
				req := &slides.CreateShapeRequest{
					ObjectId:          shapeId,
					ElementProperties: props,
					ShapeType:         each.Shape.ShapeType,
				}
				addShape := &slides.Request{CreateShape: req}
				batchReq.Requests = append(batchReq.Requests, addShape)

				if each.Shape.ShapeType == "TEXT_BOX" {
					for _, te := range each.Shape.Text.TextElements {
						if te.TextRun != nil {
							//add insert text
							insertText := &slides.InsertTextRequest{
								ObjectId: shapeId,
								Text:     te.TextRun.Content,
							}
							log.Println("add insert text")
							batchReq.Requests = append(batchReq.Requests, &slides.Request{InsertText: insertText})
						}
					}
				}

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
			// if each.Transform != nil {
			// 	todo("slide.pagelement.Transform")
			// }
			if each.Video != nil {
				todo("slide.pagelement.Video")
			}
			if each.WordArt != nil {
				todo("slide.pagelement.WordArt")
			}
		}
		log.Println("batch requests:", len(batchReq.Requests))
		batchResp, err := srv.Presentations.BatchUpdate(presentationTarget.PresentationId, batchReq).Do()
		if err != nil {
			return fmt.Errorf("Unable to send batch update to presentation: %v", err)
		}
		log.Println(batchResp)

	}

	return nil
}

func todo(path string) {
	log.Println("TODO:", path)
}
