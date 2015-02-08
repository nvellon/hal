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

	Relation string

	// Link types that store hyperlinks and its attributes.
	LinkAttr       map[string]string
	Link           LinkAttr
	LinkCollection []Link
	LinkRelations  map[Relation]LinkCollection

	// Resource is a struct that stores a resource data.
	// It represents a converted object in the HAL spec by
	// containing all its fields and also a set of related links
	// and a sub-set of recursively related resources.
	Resource struct {
		Payload  Mapper
		Links    LinkRelations
		Embedded Embedded
	}

	Embedded map[Relation][]*Resource
)

// NewResource creates a Resource object for a given struct
// and its link.
func NewResource(p Mapper, selfUri string) *Resource {
	var r Resource

	r.Payload = p

	r.Links = make(LinkRelations)
	r.AddNewLink("self", selfUri)

	r.Embedded = make(Embedded)

	return &r
}

// AddLink appends a Link to the resource.
func (r *Resource) AddLink(rel Relation, l Link) {
	r.Links[rel] = append(r.Links[rel], l)
}

// AddNewLink appends a new Link object based on
// the rel and href params.
func (r *Resource) AddNewLink(rel Relation, href string) {
	r.AddLink(rel, NewLink(href, nil))
}

// Add appends the resource into the list of embedded
// resources with the specified relation.
func (e Embedded) Add(rel Relation, r *Resource) {
	e[rel] = append(e[rel], r)
}

// Set sets the resource into the list of embedded
// resources with the specified relation. It replaces
// any existing resources associated with the relation.
func (e Embedded) Set(rel Relation, r *Resource) {
	e[rel] = []*Resource{r}
}

// Get gets the resources associated with the
// given relation.
func (e Embedded) Get(rel Relation) []*Resource {
	return e[rel]
}

// Del deletes the resources associated with the
// given relation.
func (e Embedded) Del(rel Relation) {
	delete(e, rel)
}

// Embed appends a Resource to the array of
// embedded resources.
func (r *Resource) Embed(rel Relation, er *Resource) {
	r.Embedded.Add(rel, er)
}

// Map implements the interface Mapper.
func (r Resource) GetMap() Entry {
	mapped := make(Entry)

	mp := r.Payload.GetMap()

	for k, v := range mp {
		mapped[k] = v
	}

	mapped["_links"] = r.Links

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
func NewLink(href string, attr LinkAttr) Link {
	l := make(Link)

	l["href"] = href

	for k, v := range attr {
		l[k] = v
	}

	return l
}
