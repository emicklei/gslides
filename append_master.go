package main

import (
	"log"

	"google.golang.org/api/slides/v1"
)

func appendMasterAndLayout(
	source, target *slides.Presentation,
	sourceMasterID, sourceLayoutID string,
	batch *slides.BatchUpdatePresentationRequest) string {

	masterExists := false
	for _, each := range target.Masters {
		if each.ObjectId == sourceMasterID {
			masterExists = true
			break
		}
	}
	layoutExists := false
	for _, each := range target.Layouts {
		if each.ObjectId == sourceLayoutID {
			layoutExists = true
			break
		}
	}
	if !masterExists {
		var master *slides.Page
		for _, each := range source.Masters {
			if each.ObjectId == sourceMasterID {
				master = each
				break
			}
		}
		if master == nil {
			panic("no such master" + sourceMasterID)
		}
		log.Println("copying src master:", sourceMasterID)
		appendMaster(master, batch)
	}
	if !layoutExists {
		var layout *slides.Page
		for _, each := range source.Layouts {
			if each.ObjectId == sourceLayoutID {
				layout = each
				break
			}
		}
		if layout == nil {
			panic("no such layout" + sourceLayoutID)
		}
		log.Println("copying src layout:", sourceLayoutID)
	}

	return ""
}

func appendMaster(sourceMaster *slides.Page, batch *slides.BatchUpdatePresentationRequest) {
	//dump(sourceMaster)
}
