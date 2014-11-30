// Copyright 2014 Nicolas Vellon.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package gohal implements encoding of structs into HAL as defined in
// http://stateless.co/hal_specification.html.
//
// See the basic example for an introduction to this package:
// https://github.com/nvellon/gohal/blob/master/examples/basic.go

package gohal

import (
	"encoding/json"
)

type (
	// HalEncoder is the interface implemented by the objects
	// that can be converted into HAL format.
	HalEncoder interface {
		Encode() map[string]interface{}
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
		Payload  HalEncoder
		Links    []Link
		Embedded []Resource
	}
)

// NewResource creates a Resource object for a given struct
// and its link.
func NewResource(p HalEncoder, selfUri string) Resource {
	var r Resource

	r.Payload = p

	sl := NewLink("self", selfUri)

	r.AddLink(sl)

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

// Encode implements the interface HalEncoder.
func (r Resource) Encode() map[string]interface{} {
	mapped := make(map[string]interface{})

	mp := r.Payload.Encode()

	for k, v := range mp {
		mapped[k] = v
	}

	mapped["_links"] = r.Links

	if len(r.Embedded) > 0 {
		mapped["_embedded"] = r.Embedded
	}

	return mapped
}

func (r Resource) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Encode())
}

// NewLink returns a new Link object.
func NewLink(rel, href string) Link {
	return Link{rel, href}
}

// Encode implements the interface HalEncoder.
func (l Link) Encode() map[string]interface{} {
	return map[string]interface{}{
		l.Rel: map[string]interface{}{
			"href": l.Href,
		},
	}
}

func (l Link) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.Encode())
}
