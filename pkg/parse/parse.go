package parseJSON

import (
	"encoding/json"
	"io"
	"net/http"
)

func Decode2(v any, r io.Reader) error {
	err := json.NewDecoder(r).Decode(&v)
	if err != nil {
		return err
	}
	return nil
}

func Decode(v any, r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		return err
	}
	return nil
}
