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

func (p *IDProvider) take() string {
	if p.next == len(p.ids) {
		panic("no more ids")
	}
	id := p.ids[p.next]
	p.next++
	return id
}
