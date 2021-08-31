package main

import "github.com/google/uuid"

// stack of string (generics here?)
type IDProvider struct {
	next int
	ids  []string
}

func (p *IDProvider) create() string {
	new := uuid.NewString()
	p.ids = append(p.ids, new)
	return new
}

// returns UUID and whether it was created with mapping.
func (p *IDProvider) take() (string, bool) {
	if p.next == len(p.ids) {
		return uuid.NewString(), false
	}
	id := p.ids[p.next]
	p.next++
	return id, true
}
