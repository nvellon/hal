package hal

import (
	"encoding/json"
	"testing"
)

type DummyStruct struct {
	Name string
}

func (ds DummyStruct) Encode() map[string]interface{} {
	return map[string]interface{}{
		"name": ds.Name,
	}
}

func TestNewResource(t *testing.T) {
	ds := DummyStruct{"Dummy"}

	r := NewResource(ds, "uri")

	if len(r.Links) < 1 {
		t.Errorf("No links added to the new resource")
	}

	if r.Links[0].Rel != "self" {
		t.Errorf("No SELF link added to the new resource")
	}

	if len(r.Embedded) > 0 {
		t.Errorf("Embedded list should be empty")
	}
}

func TestLinkMarshal(t *testing.T) {
	l := Link{"self", "http://localhost/"}

	jl, err := json.Marshal(l)
	if err != nil {
		t.Errorf("%s", err)
	}

	if string(jl) != `{"self":{"href":"http://localhost/"}}` {
		t.Errorf("Wrong Link struct: %s", jl)
	}
}

func TestResourceMarshal(t *testing.T) {
	expected := `{"_links":[{"self":{"href":"uri"}}],"name":"Dummy"}`

	ds := DummyStruct{"Dummy"}

	r := NewResource(ds, "uri")

	jr, err := json.Marshal(r)
	if err != nil {
		t.Errorf("%s", err)
	}

	if string(jr) != expected {
		t.Errorf("Wrong Resource struct: %s\n- Given: %s\n- Expected: %s", r, jr, expected)
	}
}
