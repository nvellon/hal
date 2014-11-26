package hal

import (
	"testing"
	"encoding/json"
)

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

func TestLinkUnmarshal(t *testing.T) {
	rel := "self"
	href := "http://localhost/"

	j := []byte(`{"` + rel + `":{"href":"` + href + `"}}`)
	l := Link{}

	err := json.Unmarshal(j, &l)
	if err != nil {
		t.Errorf("%s", err)
	}

	if l.Rel != rel || l.Href != href {
		t.Errorf("Wrong json string: %s. Given struct: %s", j, l)
	}
}