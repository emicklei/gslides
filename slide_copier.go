package main

import (
	"log"

	"github.com/google/uuid"
	"google.golang.org/api/slides/v1"
)

type SlideCopier struct {
	batch              *slides.BatchUpdatePresentationRequest
	sourceSlide        *slides.Page
	targetPresentation *slides.Presentation
	newSlideId         string
}

func (s *SlideCopier) add(req *slides.Request) {
	s.batch.Requests = append(s.batch.Requests, req)
}

func (s *SlideCopier) copyTableOfElement(elem *slides.PageElement) {
	props := new(slides.PageElementProperties) // all props set
	props.PageObjectId = s.newSlideId
	props.Size = elem.Size
	props.Transform = elem.Transform
	shapeId := uuid.NewString()
	req := &slides.CreateTableRequest{
		ObjectId:          shapeId,
		ElementProperties: props,
		Columns:           elem.Table.Columns,
		Rows:              elem.Table.Rows,
	}
	if Verbose {
		log.Println("create table:", shapeId)
	}
	s.add(&slides.Request{CreateTable: req})
	// modifiers
	// TODO how to read BorderPosition
	// {
	// 	req := &slides.UpdateTableBorderPropertiesRequest{
	// 		BorderPosition: elem.Table.
	// 	}
	// 	batch.Requests = append(batch.Requests, &slides.Request{UpdateTableBorderProperties: req})
	// }
}

func (s *SlideCopier) copyImageOfElement(elem *slides.PageElement) {
	props := new(slides.PageElementProperties) // all props set
	props.PageObjectId = s.newSlideId
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
	s.add(&slides.Request{CreateImage: req})
}

func (s *SlideCopier) copyLineOfElement(elem *slides.PageElement) {
	props := new(slides.PageElementProperties) // all props set
	props.PageObjectId = s.newSlideId
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
	s.add(&slides.Request{CreateLine: req})

	// modifiers
	{
		req := &slides.UpdateLinePropertiesRequest{
			ObjectId: shapeId,
			// direct assign lineprops?
			LineProperties: &slides.LineProperties{
				DashStyle:       elem.Line.LineProperties.DashStyle,
				StartArrow:      elem.Line.LineProperties.StartArrow,
				EndArrow:        elem.Line.LineProperties.EndArrow,
				LineFill:        elem.Line.LineProperties.LineFill,
				Weight:          elem.Line.LineProperties.Weight,
				StartConnection: elem.Line.LineProperties.StartConnection,
				EndConnection:   elem.Line.LineProperties.EndConnection,
			},
			Fields: "*",
		}
		s.add(&slides.Request{UpdateLineProperties: req})
	}
}

func (s *SlideCopier) copyShapeOfElement(elem *slides.PageElement, newShapeId string, shapeIdIsMapped bool) {
	// if shapeIdIsMapped {
	// 	log.Println("mapped with transform:", elem.Transform)
	// }

	// if the shape is mapped then it is already created by the layout else we create the extra shape
	if !shapeIdIsMapped {
		props := new(slides.PageElementProperties) // all props set
		props.PageObjectId = s.newSlideId
		props.Size = elem.Size
		props.Transform = elem.Transform
		req := &slides.CreateShapeRequest{ // all props set
			ObjectId:          newShapeId,
			ElementProperties: props,
			ShapeType:         elem.Shape.ShapeType,
		}
		if Verbose {
			log.Println("create shape:", newShapeId, " type:", elem.Shape.ShapeType)
		}
		s.add(&slides.Request{CreateShape: req})
		{
			// modifiers
			req := &slides.UpdateShapePropertiesRequest{
				ObjectId:        newShapeId,
				ShapeProperties: elem.Shape.ShapeProperties,
				// cannot use * or shapeProperties, autofit
				Fields: "shapeBackgroundFill,outline,shadow,contentAlignment,link",
			}
			s.add(&slides.Request{UpdateShapeProperties: req})
		}

	}
	if elem.Shape.ShapeType == "TEXT_BOX" {
		s.copyTextBox(elem.Shape, newShapeId, shapeIdIsMapped)
		return
	}
}

func (s *SlideCopier) copyTextBox(src *slides.Shape, shapeId string, shapeIdIsMapped bool) {
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
			s.add(&slides.Request{InsertText: insertText})

			// if the shape is mapped then it is already styled by the layout else we update the styling
			if !shapeIdIsMapped {
				if te.TextRun.Style != nil {
					updateStyle := &slides.UpdateTextStyleRequest{
						ObjectId: shapeId,
						Style:    te.TextRun.Style,
						Fields:   "*",
					}
					s.add(&slides.Request{UpdateTextStyle: updateStyle})
				}

				// TODO find example
				if te.ParagraphMarker != nil {
					updateStyle := &slides.UpdateParagraphStyleRequest{
						ObjectId: shapeId,
						Style:    te.ParagraphMarker.Style,
						Fields:   "*",
					}
					s.add(&slides.Request{UpdateParagraphStyle: updateStyle})
				}
			}
		}
	}
}
