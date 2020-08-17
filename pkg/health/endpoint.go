package health

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Response struct {
	Statuses []struct {
		Service  string      `json:"service"`
		Name     string      `json:"name"`
		Metadata interface{} `json:"metadata"`
		Healthy  bool        `json:"healthy"`
	} `json:"statuses"`
}

type EndpointCheck struct {
	endpoint string
	enabled  bool
	name     string
	err      string
}

func (e *EndpointCheck) validateResp(resp Response) bool {
	if len(resp.Statuses) == 0 {
		return true
	}
	var str strings.Builder
	str.WriteString(fmt.Sprintf("Service %s has reported the following subservices are in a degraded state: ", e.name))
	for _, status := range resp.Statuses {
		if !status.Healthy {
			str.WriteString(fmt.Sprintf("Service %s, Name %s, Metadata %v;", status.Service, status.Name, status.Metadata))
		}
	}
	e.err = str.String()
	return false
}

func (e *EndpointCheck) Check(ctx context.Context) (bool, error) {
	r, err := http.Get(e.endpoint)
	if err != nil {
		return false, err
	}
	if r.StatusCode == 200 {
		resp := Response{}
		derr := json.NewDecoder(r.Body).Decode(&resp)
		switch {
		case derr == io.EOF:
			return true, nil
		case derr != nil:
			return false, derr
		}
		return e.validateResp(resp), err
	}
	e.err = fmt.Sprintf("Service %s returned status code %s at endpoint %s", e.name, r.Status, e.endpoint)
	return false, err
}

func (e *EndpointCheck) ErrorMessage() string {
	return e.err
}

func (e *EndpointCheck) Enabled() bool {
	return e.enabled
}
