package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/kortschak/utter"
	"github.com/urfave/cli/v2"
	"google.golang.org/api/slides/v1"
)

// Known issue: https://issuetracker.google.com/issues/36761705?pli=1

// target << source[index]
func cmdAppendSlide(c *cli.Context) error {
	srv, _ := getSlidesClient()
	presentationTarget, err := srv.Presentations.Get(c.Args().First()).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve data from target presentation: %v", err)
	}
	_ = presentationTarget
	presentationSource, err := srv.Presentations.Get(c.Args().Get(1)).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve data from source presentation: %v", err)
	}
	// accept all or comma list of 1-based indices
	var indices []int
	if c.Args().Len() > 1 {
		if c.Args().Get(2) == "all" {
			indices = makeIndices(len(presentationSource.Slides))
		} else {
			for _, each := range strings.Split(c.Args().Get(2), ",") {
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

	// background
	// TODO: googleapi: Error 500: Internal error encountered., backendError
	/**
	{
		req := &slides.UpdatePagePropertiesRequest{
			ObjectId: newSlideID,
			PageProperties: &slides.PageProperties{
				PageBackgroundFill: sourceSlide.PageProperties.PageBackgroundFill,
			},
			Fields: "pageBackgroundFill.solidFill.color",
		}
		batchReq.Requests = append(batchReq.Requests, &slides.Request{UpdatePageProperties: req})
	}
	**/

	copier := &SlideCopier{
		sourceSlide:        sourceSlide,
		batch:              batchReq,
		targetPresentation: presentationTarget,
		newSlideId:         newSlideID,
	}

	// elements
	for _, each := range sourceSlide.PageElements {
		if each.Shape != nil {
			id, isMapped := ids.take()
			copier.copyShapeOfElement(each, id, isMapped)
		}
		if each.ElementGroup != nil {
			todo("slide.pagelement.ElementGroup")
		}
		if each.Image != nil {
			copier.copyImageOfElement(each)
		}
		if each.Line != nil {
			copier.copyLineOfElement(each)
		}
		if each.Table != nil {
			copier.copyTableOfElement(each)
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
