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
	l := Link{}
	j := []byte(`{"self":{"href":"http://localhost/"}}`)

	err := json.Unmarshal(j, &l)
	if err != nil {
		t.Errorf("%s", err)
	}

	if l.Rel != "self" || l.Href != "http://localhost/" {
		t.Errorf("Wrong json string: %s. Given struct: %s", j, l)
	}
}