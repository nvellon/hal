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
	"reflect"
	"strings"
)

type (
	Entry map[string]interface{}

	// Mapper is the interface implemented by the objects
	// that can be converted into HAL format.
	Mapper interface {
		GetMap() Entry
	}

	Relation string

	CurieHandle struct {
		Name string
		*Resource
	}

	// Link types that store hyperlinks and its attributes.
	LinkAttr       map[string]interface{}
	Link           LinkAttr
	LinkCollection []Link
	LinkRelations  map[Relation]interface{}

	// Resource is a struct that stores a resource data.
	// It represents a converted object in the HAL spec by
	// containing all its fields and also a set of related links
	// and a sub-set of recursively related resources.
	Resource struct {
		Payload  interface{}
		Links    LinkRelations
		Embedded Embedded
		Curies   map[string]*CurieHandle
	}
	ResourceCollection []*Resource

	Embedded map[Relation]interface{}
)

// AddNewLink adds a link to the resources_link collection
// prepended with the curie Name
func (c CurieHandle) AddNewLink(rel Relation, href string) {
	rel = Relation(c.Name) + ":" + rel
	c.Resource.AddLink(rel, NewLink(href, nil))
}

// AddCollection appends the resource into the list of embedded
// resources with the specified relation.
// r should be  a ResourceCollection
func (e Embedded) AddCollection(rel Relation, r ResourceCollection) {
	n := e[rel]
	if n == nil {
		//new embed
		e[rel] = r
		return
	}

	if nc, ok := n.([]*Resource); ok {
		e[rel] = append(nc, r...)
		return
	}

	if nr, ok := n.(*Resource); ok {
		e[rel] = append([]*Resource{nr}, r...)
	}

}

// Add appends the resource into the list of embedded
// resources with the specified relation.
// r should be a Resource
func (e Embedded) Add(rel Relation, r *Resource) {
	n := e[rel]
	if n == nil {
		//new embed
		e[rel] = r
		return
	}

	if nec, ok := n.([]*Resource); ok {
		e[rel] = append(nec, r)
		return
	}

	if nee, ok := n.(*Resource); ok {
		e[rel] = append([]*Resource{nee}, r)
		return
	}

	//something went wrong.. replace what is there with what is new
	e[rel] = []*Resource{r}
}

// Set sets the resource into the list of embedded
// resources with the specified relation. It replaces
// any existing resources associated with the relation.
// r should be a pointer to a Resource
func (e Embedded) Set(rel Relation, r *Resource) {
	e[rel] = r
}

// Set sets the resource into the list of embedded
// resources with the specified relation. It replaces
// any existing resources associated with the relation.
// r should be a ResourceCollection
func (e Embedded) SetCollection(rel Relation, r ResourceCollection) {
	e[rel] = r
}

// Get gets the resources associated with the
// given relation.
//func (e Embedded) Get(rel Relation) []*Resource {
//	return e[rel]
//}

// Del deletes the resources associated with the
// given relation.
func (e Embedded) Del(rel Relation) {
	delete(e, rel)
}

// NewResource creates a Resource object for a given struct
// and its link.
func NewResource(p interface{}, selfUri string) *Resource {
	var r Resource

	r.Payload = p

	r.Links = make(LinkRelations)
	r.AddNewLink("self", selfUri)

	r.Embedded = make(Embedded)
	r.Curies = make(map[string]*CurieHandle)

	return &r
}

// AddLinkCollection appends a LinkCollection to the resource.
// l should be a LinkCollection
func (r *Resource) AddLinkCollection(rel Relation, l LinkCollection) {
	n := r.Links[rel]
	if n == nil {
		//new link
		r.Links[rel] = l
		return
	}

	if nc, ok := n.(LinkCollection); ok {
		r.Links[rel] = append(nc, l...)
		return
	}

	if nl, ok := n.(Link); ok {
		//prepend existing link to collection
		r.Links[rel] = append(LinkCollection{nl}, l...)
	}
}

// AddLink appends a Link to the resource.
// l should be a Link
func (r *Resource) AddLink(rel Relation, l Link) {
	n := r.Links[rel]
	if n == nil {
		//new link
		r.Links[rel] = l
		return
	}

	if nc, ok := n.(LinkCollection); ok {
		r.Links[rel] = append(nc, l)
		return
	}

	if nl, ok := n.(Link); ok {
		r.Links[rel] = append(LinkCollection{nl}, l)
		return
	}

	//something went wrong.. replace what is there with what is new
	r.Links[rel] = LinkCollection{l}
}

// AddNewLink appends a new Link object based on
// the rel and href params.
func (r *Resource) AddNewLink(rel Relation, href string) {
	r.AddLink(rel, NewLink(href, nil))
}

// RegisterCurie adds a Link relation of type 'curies' and returns a CurieHandle
// to allow users to fluently add new links that have this curie relation definition
func (r *Resource) RegisterCurie(name, href string, templated bool) *CurieHandle {
	l := LinkCollection{
		NewLink(href, LinkAttr{"name": name}, LinkAttr{"templated": templated}),
	}
	r.AddLinkCollection("curies", l)

	handle := &CurieHandle{Name: name, Resource: r}

	r.Curies[name] = handle
	return handle
}

// Embed appends a Resource to the array of
// embedded resources.
// re should be a pointer to a Resource
func (r *Resource) Embed(rel Relation, re *Resource) {
	r.Embedded.Add(rel, re)
}

// EmbedCollection appends a ResourceCollection to the array of
// embedded resources.
// re should be a ResourceCollection
func (r *Resource) EmbedCollection(rel Relation, re ResourceCollection) {
	r.Embedded.AddCollection(rel, re)
}

// Map implements the interface Mapper.
func (r Resource) GetMap() Entry {
	mapped := make(Entry)

	var mp Entry
	// Check if payload implements Mapper interface
	if mapper, ok := r.Payload.(Mapper); ok {
		mp = mapper.GetMap()
	} else {
		mp = r.getPayloadMap()
	}

	for k, v := range mp {
		mapped[k] = v
	}

	mapped["_links"] = r.Links

	if len(r.Embedded) > 0 {
		mapped["_embedded"] = r.Embedded
	}

	return mapped
}

func (r *Resource) getPayloadMap() Entry {

	val := reflect.ValueOf(r.Payload)
	payloadMap := Entry{}

	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		tagValue := tag.Get("json")
		if strings.Contains(tagValue, "omitempty") {
			l := strings.Split(tagValue, ",")
			for i, el := range l {
				if el == "omitempty" {
					l = append(l[:i], l[i+1:]...)
					break
				}
			}
			tagValue = strings.Join(l, ",")

		}
		if tagValue != "-" {
			valueField := val.Field(i)

			if tagValue == "" {
				tagValue = typeField.Name
			}

			payloadMap[tagValue] = valueField.Interface()
		}
	}

	return payloadMap
}

// MarshalJSON is a Marshaler interface implementation
// for Resource struct
func (r Resource) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.GetMap())
}

// NewLink returns a new Link object.
func NewLink(href string, attrs ...LinkAttr) Link {
	l := make(Link)

	l["href"] = href

	for _, attr := range attrs {
		for k, v := range attr {
			l[k] = v
		}
	}

	return l
}
