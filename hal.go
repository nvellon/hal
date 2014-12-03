// Copyright 2014 Nicolas Vellon.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package hal implements encoding of structs into HAL as defined in
// http://stateless.co/hal_specification.html.
//
// See the basic example for an introduction to this package:
// https://github.com/nvellon/hal/blob/master/examples/basic.go

package hal

import (
	"encoding/json"
)

type (
	Entry map[string]interface{}

	// Mapper is the interface implemented by the objects
	// that can be converted into HAL format.
	Mapper interface {
		GetMap() Entry
	}

	// Link is a struct that stores a hyperlink data.
	Link struct {
		Rel  string
		Href string
	}

	// Resource is a struct that stores a resource data.
	// It represents a converted object in the HAL spec by
	// containing all its fields and also a set of related links
	// and a sub-set of recursively related resources.
	Resource struct {
		Payload  Mapper
		Links    []Link
		Embedded []Resource
	}
)

// NewResource creates a Resource object for a given struct
// and its link.
func NewResource(p Mapper, selfUri string) Resource {
	var r Resource

	r.Payload = p

	r.AddLink(NewLink("self", selfUri))

	return r
}

// AddLink appends a Link to the resource.
func (r *Resource) AddLink(l Link) {
	r.Links = append(r.Links, l)
}

// Embed appends a Resource to the array of
// embedded resources.
func (r *Resource) Embed(er Resource) {
	r.Embedded = append(r.Embedded, er)
}

// Map implements the interface Mapper.
func (r Resource) GetMap() Entry {
	mapped := make(Entry)

	mp := r.Payload.GetMap()

	for k, v := range mp {
		mapped[k] = v
	}

	ml := make(Entry)

	for _, link := range r.Links {
		el := link.GetMap()
		for rel, l := range el {
			ml[rel] = l
		}
	}

	mapped["_links"] = ml

	if len(r.Embedded) > 0 {
		mapped["_embedded"] = r.Embedded
	}

	return mapped
}

// MarshalJSON is a Marshaler interface implementation
// for Resource struct
func (r Resource) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.GetMap())
}

// NewLink returns a new Link object.
func NewLink(rel, href string) Link {
	return Link{rel, href}
}

// Map implements the interface Mapper.
func (l Link) GetMap() Entry {
	return Entry{
		l.Rel: Entry{
			"href": l.Href,
		},
	}
}

// MarshalJSON is a Marshaler interface implementation
// for Link struct
func (l Link) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.GetMap())
}
