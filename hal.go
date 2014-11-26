package hal

import (
	"encoding/json"
)

type (
	Link struct {
		Rel		string
		Href	string
	}

	Resource struct {
		Links		[]Link		`json:"_links"`
		Embedded	[]Resource	`json:"_embedded"`
	}
)

func (l Link) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		map[string]interface{}{
			l.Rel: map[string]interface{}{
				"href": l.Href,
			},
	})
}

func (l *Link) UnmarshalJSON(b []byte) error {
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