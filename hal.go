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
		Links		[]Link
		Embedded	[]Resource
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
	var m map[string]*json.RawMessage

	err := json.Unmarshal(b, &m)

	if err != nil {
		return err
	}

	for key, value := range m {
		l.Rel = string(key)

		var ml map[string]*json.RawMessage

		err := json.Unmarshal(*value, &ml)
		if err != nil {
			return err
		}

		l.Href = string(*ml["href"])
	}

	return nil
}