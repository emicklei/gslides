package main

import (
	"log"

	"github.com/kortschak/utter"
	"google.golang.org/api/slides/v1"
)

func dump(what interface{}) {
	utter.Config.OmitZero = false
	utter.Dump(what)
}

func todo(path string) {
	log.Println("TODO:", path)
}

func layoutIDWithName(p *slides.Presentation, name string) string {
	for _, each := range p.Layouts {
		if each.LayoutProperties.Name == name {
			return each.ObjectId
		}
	}
	return "?"
}

func layoutNameWithID(p *slides.Presentation, id string) string {
	for _, each := range p.Layouts {
		if each.ObjectId == id {
			return each.LayoutProperties.Name
		}
	}
	return "?"
}
