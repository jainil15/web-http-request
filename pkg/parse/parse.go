package parseJSON

import (
	"encoding/json"
	"net/http"
)

func Decode(v any, r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		return err
	}
	return nil
}
