package types

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type V1HttpHandlerFunc func(http.ResponseWriter, *http.Request, *ApiConfig)

func BodyDecoder[T interface{}](r *http.Request) (*T, error) {
	var body T
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %s", err)
	}
	return &body, nil
}