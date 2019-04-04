package hal

import (
	"testing"
)

type DummyStruct struct {
	Name string `json:"name"`
}

func TestNewResource(t *testing.T) {
	ds := DummyStruct{"Dummy"}

	r := NewResource(ds, "uri")

	if len(r.Links) < 1 {
		t.Errorf("No links added to the new resource")
	}

	if r.Links["self"] == nil {
		t.Errorf("No SELF link added to the new resource")
	}

	if len(r.Embedded) > 0 {
		t.Errorf("Embedded list should be empty")
	}
}

func TestLinkMarshal(t *testing.T) {
	l := make(Link)
	l["href"] = "http://localhost/"

	jl, err := json.Marshal(l)
	if err != nil {
		t.Errorf("%s", err)
	}

	if string(jl) != `{"href":"http://localhost/"}` {
		t.Errorf("Wrong Link struct: %s", jl)
	}
}

func TestResourceMarshal(t *testing.T) {
	expected := `{"_links":{"self":{"href":"uri"}},"name":"Dummy"}`

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

type DummyStructWithMapper struct {
	Name string
}

func (dswm DummyStructWithMapper) GetMap() Entry {
	return Entry{
		"customName": dswm.Name,
	}
}

func TestResourceMarshallWithMapper(t *testing.T) {
	expected := `{"_links":{"self":{"href":"uri"}},"customName":"Dummy"}`

	ds := DummyStructWithMapper{"Dummy"}

	r := NewResource(ds, "uri")

	jr, err := json.Marshal(r)
	if err != nil {
		t.Errorf("%s", err)
	}

	if string(jr) != expected {
		t.Errorf("Wrong Resource struct: %s\n- Given: %s\n- Expected: %s", r, jr, expected)
	}
}

/* Test Links */
func TestNewLink(t *testing.T) {
	expected := `{"href":"bar","templated":true}`

	l := NewLink("bar", LinkAttr{"templated": true})

	jr, err := json.Marshal(l)
	if err != nil {
		t.Errorf("%s", err)
	}

	if string(jr) != expected {
		t.Errorf("Wrong link struct: %s\n- Given: %s\n- Expected: %s", l, jr, expected)
	}
}

func TestNewLinkMultipleAttributes(t *testing.T) {
	expected := `{"href":"http://haltalk.herokuapp.com/docs/{rel}","name":"doc","templated":true}`

	l := NewLink("http://haltalk.herokuapp.com/docs/{rel}", LinkAttr{"name": "doc"}, LinkAttr{"templated": true})

	jr, err := json.Marshal(l)
	if err != nil {
		t.Errorf("%s", err)
	}

	if string(jr) != expected {
		t.Errorf("Wrong link struct: %s\n- Given: %s\n- Expected: %s", l, jr, expected)
	}
}

func TestRegisterCurie(t *testing.T) {
	expected := `{"_links":{"curies":[{"href":"http://haltalk.herokuapp.com/docs/{rel}","name":"doc","templated":true}],"doc:foo":{"href":"bar"},"self":{"href":"uri"}},"name":"Dummy"}`

	ds := DummyStruct{"Dummy"}

	r := NewResource(ds, "uri")
	r.RegisterCurie("doc", "http://haltalk.herokuapp.com/docs/{rel}", true).AddNewLink("foo", "bar")

	jr, err := json.Marshal(r)
	if err != nil {
		t.Errorf("%s", err)
	}

	if string(jr) != expected {
		t.Errorf("Wrong Resource struct: %s\n- Given: %s\n- Expected: %s", r, jr, expected)
	}
}

func TestRegisterMultipleCuries(t *testing.T) {
	expected := `{"_links":{"curies":[{"href":"http://haltalk.herokuapp.com/docs/{rel}","name":"doc","templated":true},{"href":"http://haltalk.herokuapp.com/abc/{rel}","name":"abc","templated":true}],"doc:foo":{"href":"bar"},"self":{"href":"uri"}},"name":"Dummy"}`

	ds := DummyStruct{"Dummy"}

	r := NewResource(ds, "uri")
	r.RegisterCurie("doc", "http://haltalk.herokuapp.com/docs/{rel}", true).AddNewLink("foo", "bar")
	r.RegisterCurie("abc", "http://haltalk.herokuapp.com/abc/{rel}", true)

	jr, err := json.Marshal(r)
	if err != nil {
		t.Errorf("%s", err)
	}

	if string(jr) != expected {
		t.Errorf("Wrong Resource struct: %s\n- Given: %s\n- Expected: %s", r, jr, expected)
	}
}

func TestResourceCuries(t *testing.T) {
	ds := DummyStruct{"Dummy"}
	curieName := "doc"

	r := NewResource(ds, "uri")
	curie := r.RegisterCurie(curieName, "http://haltalk.herokuapp.com/docs/{rel}", true)

	curie.AddNewLink("foo", "bar")
	curies := r.Curies

	if len(curies) != 1 {
		t.Errorf("Wrong number of CurieHandles returned from resource:\n - Given: %v\n- Expected: %v\n", len(curies), 1)
	}

	if curies[curieName].Resource != r {
		t.Errorf("CurieHandle.Resource does not reference owning resource")
	}

	if curie != curies[curieName] {
		t.Errorf("curieHandle returned by RegisterCurie() is not the same reference")
	}

}

func TestAddNewLink(t *testing.T) {
	expected := `{"_links":{"foo":{"href":"bar"},"self":{"href":"uri"}},"name":"Dummy"}`

	ds := DummyStruct{"Dummy"}

	r := NewResource(ds, "uri")
	r.AddNewLink("foo", "bar")

	jr, err := json.Marshal(r)
	if err != nil {
		t.Errorf("%s", err)
	}

	if string(jr) != expected {
		t.Errorf("Wrong Resource struct: %s\n- Given: %s\n- Expected: %s", r, jr, expected)
	}
}

func TestAddNewLinkTwice(t *testing.T) {
	expected := `{"_links":{"foo":[{"href":"bar"},{"href":"bar2"}],"self":{"href":"uri"}},"name":"Dummy"}`

	ds := DummyStruct{"Dummy"}

	r := NewResource(ds, "uri")
	r.AddNewLink("foo", "bar")
	r.AddNewLink("foo", "bar2")

	jr, err := json.Marshal(r)
	if err != nil {
		t.Errorf("%s", err)
	}

	if string(jr) != expected {
		t.Errorf("Wrong Resource struct: %s\n- Given:    %s\n- Expected: %s", r, jr, expected)
	}
}

func TestAddLinkCollection(t *testing.T) {
	expected := `{"_links":{"foo":[{"href":"bar"},{"href":"bar2"}],"self":{"href":"uri"}},"name":"Dummy"}`

	ds := DummyStruct{"Dummy"}

	r := NewResource(ds, "uri")
	r.AddLinkCollection("foo", LinkCollection{NewLink("bar", nil), NewLink("bar2", nil)})

	jr, err := json.Marshal(r)
	if err != nil {
		t.Errorf("%s", err)
	}

	if string(jr) != expected {
		t.Errorf("Wrong Resource struct: %s\n- Given:    %s\n- Expected: %s", r, jr, expected)
	}
}

func TestAddLinkCollectionToLink(t *testing.T) {
	expected := `{"_links":{"foo":[{"href":"baz"},{"href":"bar"},{"href":"bar2"}],"self":{"href":"uri"}},"name":"Dummy"}`

	ds := DummyStruct{"Dummy"}

	r := NewResource(ds, "uri")
	r.AddNewLink("foo", "baz")
	r.AddLinkCollection("foo", LinkCollection{NewLink("bar", nil), NewLink("bar2", nil)})

	jr, err := json.Marshal(r)
	if err != nil {
		t.Errorf("%s", err)
	}

	if string(jr) != expected {
		t.Errorf("Wrong Resource struct: %s\n- Given:    %s\n- Expected: %s", r, jr, expected)
	}
}

/* Test Embedded */
func TestEmbed(t *testing.T) {
	expected := `{"_embedded":{"foo":{"_links":{"self":{"href":"uri2"}},"name":"DummyEmbed"}},"_links":{"self":{"href":"uri"}},"name":"Dummy"}`

	ds := DummyStruct{"Dummy"}
	ds2 := DummyStruct{"DummyEmbed"}

	r := NewResource(ds, "uri")
	r2 := NewResource(ds2, "uri2")
	r.Embed("foo", r2)

	jr, err := json.Marshal(r)
	if err != nil {
		t.Errorf("%s", err)
	}

	if string(jr) != expected {
		t.Errorf("Wrong Resource struct: %s\n- Given: %s\n- Expected: %s", r, jr, expected)
	}
}

func TestEmbedTwice(t *testing.T) {
	expected := `{"_embedded":{"foo":[{"_links":{"self":{"href":"uri2"}},"name":"DummyEmbed"},{"_links":{"self":{"href":"uri3"}},"name":"DummyEmbed2"}]},"_links":{"self":{"href":"uri"}},"name":"Dummy"}`

	ds := DummyStruct{"Dummy"}
	ds2 := DummyStruct{"DummyEmbed"}
	ds3 := DummyStruct{"DummyEmbed2"}

	r := NewResource(ds, "uri")
	r2 := NewResource(ds2, "uri2")
	r3 := NewResource(ds3, "uri3")
	r.Embed("foo", r2)
	r.Embed("foo", r3)

	jr, err := json.Marshal(r)
	if err != nil {
		t.Errorf("%s", err)
	}

	if string(jr) != expected {
		t.Errorf("Wrong Resource struct: %s\n- Given:    %s\n- Expected: %s", r, jr, expected)
	}
}

func TestAddResourceCollection(t *testing.T) {
	expected := `{"_embedded":{"foo":[{"_links":{"self":{"href":"uri2"}},"name":"DummyEmbed"},{"_links":{"self":{"href":"uri3"}},"name":"DummyEmbed2"}]},"_links":{"self":{"href":"uri"}},"name":"Dummy"}`

	ds := DummyStruct{"Dummy"}
	ds2 := DummyStruct{"DummyEmbed"}
	ds3 := DummyStruct{"DummyEmbed2"}

	r := NewResource(ds, "uri")
	r2 := NewResource(ds2, "uri2")
	r3 := NewResource(ds3, "uri3")
	r.EmbedCollection("foo", ResourceCollection{r2, r3})

	jr, err := json.Marshal(r)
	if err != nil {
		t.Errorf("%s", err)
	}

	if string(jr) != expected {
		t.Errorf("Wrong Resource struct: %s\n- Given:    %s\n- Expected: %s", r, jr, expected)
	}
}

func TestAddResourceCollectionToResource(t *testing.T) {
	expected := `{"_embedded":{"foo":[{"_links":{"self":{"href":"uri2"}},"name":"DummyEmbed"},{"_links":{"self":{"href":"uri3"}},"name":"DummyEmbed2"},{"_links":{"self":{"href":"uri4"}},"name":"DummyEmbed3"}]},"_links":{"self":{"href":"uri"}},"name":"Dummy"}`

	ds := DummyStruct{"Dummy"}
	ds2 := DummyStruct{"DummyEmbed"}
	ds3 := DummyStruct{"DummyEmbed2"}
	ds4 := DummyStruct{"DummyEmbed3"}

	r := NewResource(ds, "uri")
	r2 := NewResource(ds2, "uri2")
	r3 := NewResource(ds3, "uri3")
	r4 := NewResource(ds4, "uri4")
	r.Embed("foo", r2)
	r.EmbedCollection("foo", ResourceCollection{r3, r4})

	jr, err := json.Marshal(r)
	if err != nil {
		t.Errorf("%s", err)
	}

	if string(jr) != expected {
		t.Errorf("Wrong Resource struct: %s\n- Given:    %s\n- Expected: %s", r, jr, expected)
	}
}

func TestOmitEmptyReflection(t *testing.T) {
	expected := `{"_links":{"self":{"href":"test"}},"id":null}`
	dummyStruct := struct {
		ID *int `json:"id,omitempty"`
	}{}
	r := NewResource(dummyStruct,"test")
	res, err := json.Marshal(r)
	if err != nil {
		t.Errorf("%s", err)
	}
	if string(res) != expected {
		t.Errorf("Wrong Resource struct: %s\n- Given:    %s\n- Expected: %s", r, res, expected)
	}
}
