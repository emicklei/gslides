package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/kortschak/utter"
	"github.com/urfave/cli"
	"google.golang.org/api/slides/v1"
)

// Known issue: https://issuetracker.google.com/issues/36761705?pli=1

// target << source[index]
func cmdAppendSlide(c *cli.Context) error {
	srv, _ := getSlidesClient()
	presentationTarget, err := srv.Presentations.Get(c.Args()[0]).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve data from target presentation: %v", err)
	}
	_ = presentationTarget
	presentationSource, err := srv.Presentations.Get(c.Args()[1]).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve data from source presentation: %v", err)
	}
	// accept all or comma list of 1-based indices
	var indices []int
	if len(c.Args()) > 1 {
		if c.Args()[2] == "all" {
			indices = makeIndices(len(presentationSource.Slides))
		} else {
			for _, each := range strings.Split(c.Args()[2], ",") {
				sourceSlideIndex, err := strconv.Atoi(each)
				if err != nil {
					return fmt.Errorf("invalid slide presentation index: %v", err)
				}
				indices = append(indices, sourceSlideIndex)
			}
		}
	}
	if len(indices) == 0 {
		return fmt.Errorf("missing comma separated list of slide indices or [all]")
	}
	// collect all changes
	batchReq := new(slides.BatchUpdatePresentationRequest)
	for _, each := range indices {
		// zero indexed
		if err := appendSlide(each-1, presentationSource, presentationTarget, batchReq); err != nil {
			return fmt.Errorf("unable to append slide presentation index: %d error:%v", each-1, err)
		}
	}
	_, err = srv.Presentations.BatchUpdate(presentationTarget.PresentationId, batchReq).Do()
	if err != nil {
		return fmt.Errorf("unable to send batch update to presentation: %v", err)
	}
	return nil
}

func appendSlide(sourceSlideIndex int, presentationSource, presentationTarget *slides.Presentation, batchReq *slides.BatchUpdatePresentationRequest) error {
	if sourceSlideIndex >= len(presentationSource.Slides) {
		return fmt.Errorf("no such slide index: %v", sourceSlideIndex)
	}
	sourceSlide := presentationSource.Slides[sourceSlideIndex]
	sourceLayoutName := layoutNameWithID(presentationSource, sourceSlide.SlideProperties.LayoutObjectId)
	log.Println("src layout name:", sourceLayoutName)
	layoutMappings := []*slides.LayoutPlaceholderIdMapping{}
	ids := new(IDProvider)
	for _, each := range sourceSlide.PageElements {
		if each.Shape != nil && each.Shape.Placeholder != nil {
			newID := ids.create()
			layoutMappings = append(layoutMappings, &slides.LayoutPlaceholderIdMapping{
				ObjectId: newID,
				LayoutPlaceholder: &slides.Placeholder{
					Index: int64(each.Shape.Placeholder.Index),
					Type:  each.Shape.Placeholder.Type,
				},
			})
			log.Println("new mapping", "index:", each.Shape.Placeholder.Index, "type", each.Shape.Placeholder.Type, "->", "id:", newID)
		}
	}

	newSlideID := uuid.NewString()

	addSlide := &slides.CreateSlideRequest{
		ObjectId: newSlideID,
		SlideLayoutReference: &slides.LayoutReference{
			PredefinedLayout: sourceLayoutName,
		},
		PlaceholderIdMappings: layoutMappings,
	}
	batchReq.Requests = append(batchReq.Requests, &slides.Request{CreateSlide: addSlide})

	// elements
	for _, each := range sourceSlide.PageElements {
		if Verbose {
			log.Println("title:", each.Title, " description:", each.Description, "element group", each.ElementGroup)
		}
		if each.Shape != nil {
			id, isMapped := ids.take()
			copyShapeOfElement(each, newSlideID, id, isMapped, batchReq)
		}
		if each.ElementGroup != nil {
			todo("slide.pagelement.ElementGroup")
		}
		if each.Image != nil {
			copyImageOfElement(each, newSlideID, batchReq)
		}
		if each.Line != nil {
			copyLineOfElement(each, newSlideID, batchReq)
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

	// Send the batch
	if Verbose && false {
		log.Println("target batch requests:", len(batchReq.Requests))
		for _, each := range batchReq.Requests {
			utter.Config.OmitZero = true
			fmt.Println(utter.Sdump(each))
		}
	}
	return nil
}

func copyImageOfElement(elem *slides.PageElement, newSlideId string, batch *slides.BatchUpdatePresentationRequest) {
	props := new(slides.PageElementProperties) // all props set
	props.PageObjectId = newSlideId
	props.Size = elem.Size
	props.Transform = elem.Transform
	shapeId := uuid.NewString()
	url := elem.Image.SourceUrl
	if len(url) == 0 {
		url = elem.Image.ContentUrl
	}
	req := &slides.CreateImageRequest{ // all props set
		ObjectId:          shapeId,
		ElementProperties: props,
		Url:               url,
	}
	if Verbose {
		log.Println("create image:", shapeId, " url:", url)
	}
	batch.Requests = append(batch.Requests, &slides.Request{CreateImage: req})
}

func copyLineOfElement(elem *slides.PageElement, newSlideId string, batch *slides.BatchUpdatePresentationRequest) {
	props := new(slides.PageElementProperties) // all props set
	props.PageObjectId = newSlideId
	props.Size = elem.Size
	props.Transform = elem.Transform
	shapeId := uuid.NewString()

	req := &slides.CreateLineRequest{ // all props set
		ObjectId:          shapeId,
		ElementProperties: props,
		Category:          elem.Line.LineCategory,
	}
	if Verbose {
		log.Println("create line:", shapeId, " category:", elem.Line.LineCategory)
	}
	batch.Requests = append(batch.Requests, &slides.Request{CreateLine: req})

	// modifiers
	{
		req := &slides.UpdateLineCategoryRequest{}
		batch.Requests = append(batch.Requests, &slides.Request{CreateLine: req})
	}
}

func copyShapeOfElement(elem *slides.PageElement, newSlideId, shapeId string, shapeIdIsMapped bool, batch *slides.BatchUpdatePresentationRequest) {
	// if the shape is mapped then it is already created by the layout else we create the extra shape
	if !shapeIdIsMapped {
		props := new(slides.PageElementProperties) // all props set
		props.PageObjectId = newSlideId
		props.Size = elem.Size
		props.Transform = elem.Transform
		req := &slides.CreateShapeRequest{ // all props set
			ObjectId:          shapeId,
			ElementProperties: props,
			ShapeType:         elem.Shape.ShapeType,
		}
		if Verbose {
			log.Println("create shape:", shapeId, " type:", elem.Shape.ShapeType)
		}
		batch.Requests = append(batch.Requests, &slides.Request{CreateShape: req})
	}
	if elem.Shape.ShapeType == "TEXT_BOX" {
		copyTextBox(elem.Shape, shapeId, shapeIdIsMapped, batch)
		return
	}
	todo(elem.Shape.ShapeType)
}

func copyTextBox(src *slides.Shape, shapeId string, shapeIdIsMapped bool, batch *slides.BatchUpdatePresentationRequest) {
	if src.Text == nil {
		if Verbose {
			log.Println("skip TEXT_BOX shape without Text (nil)")
		}
		return
	}
	for _, te := range src.Text.TextElements {
		if te.AutoText != nil {
			todo("text box.auto text")
		}
		if te.TextRun != nil {
			insertText := &slides.InsertTextRequest{
				ObjectId: shapeId,
				Text:     te.TextRun.Content,
			}
			if Verbose {
				log.Println("textbox:", shapeId, " gets text:", te.TextRun.Content)
			}
			batch.Requests = append(batch.Requests, &slides.Request{InsertText: insertText})

			// if the shape is mapped then it is already styled by the layout else we update the styling
			if !shapeIdIsMapped {
				if te.TextRun.Style != nil {
					updateStyle := &slides.UpdateTextStyleRequest{
						ObjectId: shapeId,
						Style:    te.TextRun.Style,
						Fields:   "*",
					}
					batch.Requests = append(batch.Requests, &slides.Request{UpdateTextStyle: updateStyle})
				}

				// TODO find example
				if te.ParagraphMarker != nil {
					updateStyle := &slides.UpdateParagraphStyleRequest{
						ObjectId: shapeId,
						Style:    te.ParagraphMarker.Style,
						Fields:   "*",
					}
					batch.Requests = append(batch.Requests, &slides.Request{UpdateParagraphStyle: updateStyle})
				}
			}
		}
	}
}
