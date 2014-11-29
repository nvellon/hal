package gohal

import (
	"encoding/json"
)

type (
	LinkTranslator interface {
		ToLink(b []byte) error
		ToMap() map[string]interface{}
	}

	ResourceTranslator interface {
		ToResource(b []byte) error
		ToMap() map[string]interface{}
	}

	Link struct {
		Rel  string
		Href string
	}

	Resource struct {
		Links    []LinkTranslator     `json:"_links"`
		Embedded []ResourceTranslator `json:"_embedded"`
	}
)

// Link object methods
func NewLink(rel, href string) Link {
	return Link{rel, href}
}

func (l Link) ToMap() map[string]interface{} {
	return map[string]interface{}{
		l.Rel: map[string]interface{}{
			"href": l.Href,
		},
	}
}

func (l *Link) ToLink(b []byte) error {
	var m map[string]json.RawMessage

	err := json.Unmarshal(b, &m)

	if err != nil {
		return err
	}

	for rel, raw := range m {
		l.Rel = rel

		var ml map[string]string

		err := json.Unmarshal(raw, &ml)
		if err != nil {
			return err
		}

		l.Href = string(ml["href"])
	}

	return nil
}

func (l Link) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.ToMap())
}

func (l *Link) UnmarshalJSON(b []byte) error {
	return l.ToLink(b)
}

// Resource object methods
func NewResource(l LinkTranslator) Resource {
	var r Resource

	r.AddLink(l)

	return r
}

func (r *Resource) AddLink(l LinkTranslator) {
	r.Links = append(r.Links, l)
}

func (r *Resource) Embed(nr ResourceTranslator) {
	r.Embedded = append(r.Embedded, nr)
}

/*
func (r Resource) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"_links": r.Rel(),
		"_embedded":
	}
}

func (r *Resource) ToResource(b []byte) error {
	var m map[string]json.RawMessage

	err := json.Unmarshal(b, &m)

	if err != nil {
		return err
	}

	for rel, raw := range m {
		l.rel = rel

		var ml map[string]string

		err := json.Unmarshal(raw, &ml)
		if err != nil {
			return err
		}

		l.href = string(ml["href"])
	}

	return nil
}

func (r Resource) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.ToMap())
}

func (r *Resource) UnmarshalJSON(b []byte) error {
	return r.ToResource(b)
}
*/
