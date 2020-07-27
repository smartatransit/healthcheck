package health

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Response struct {
	Service  string      `json:"service"`
	Name     string      `json:"name"`
	Metadata interface{} `json:"metadata"`
}

type EndpointCheck struct {
	endpoint string
	enabled  bool
	name     string
	err      string
}

func (e *EndpointCheck) Check(ctx context.Context) (bool, error) {
	r, err := http.Get(e.endpoint)
	if err != nil {
		return false, err
	}
	if r.Status == "200" {
		resp := Response{}
		derr := json.NewDecoder(r.Body).Decode(resp)
		switch {
		case derr == io.EOF:
			return true, nil
		case derr != nil:
			return false, derr
		}
		return true, err
	}
	e.err = fmt.Sprintf("Service %s returned status code %s at endpoint %s", e.name, r.Status, e.endpoint)
	return false, err
}

func (e *EndpointCheck) ErrorMessage() string {
	return e.err
}
