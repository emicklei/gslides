package main

import "github.com/skinass/go-spew/spew"

func dump(what interface{}) {

	spew.Config.DisableNilValues = true
	spew.Config.DisableZeroValues = true
	spew.Dump(what)
}
