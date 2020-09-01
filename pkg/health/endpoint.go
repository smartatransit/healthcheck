package health

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}

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
	client   HTTPClient
}

func (e *EndpointCheck) validateResp(resp Response) bool {
	var unhealthy bool
	var str strings.Builder
	str.WriteString(fmt.Sprintf("Service %s has reported the following subservices are in a degraded state: ", e.name))
	for _, status := range resp.Statuses {
		if !status.Healthy {
			unhealthy = true
			str.WriteString(fmt.Sprintf("Service %s, Name %s, Metadata %v;", status.Service, status.Name, status.Metadata))
		}
	}
	if unhealthy {
		e.err = str.String()
		return false
	}
	return true
}

func (e *EndpointCheck) Check(ctx context.Context) (bool, error) {
	r, err := e.client.Get(e.endpoint)
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
