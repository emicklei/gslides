package main

import "google.golang.org/api/slides/v1"

// return id of new master
func appendMaster(
	source, target *slides.Presentation,
	sourceLayoutID, sourceMasterID string,
	batch *slides.BatchUpdatePresentationRequest) string {
	return ""
}
