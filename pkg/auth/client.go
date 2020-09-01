package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type TokenResponse struct {
	Token string `json:"token"`
}

type Client struct {
	http   *http.Client
	logger *zap.Logger
}

func NewClient(http *http.Client, logger *zap.Logger) *Client {
	return &Client{
		http,
		logger,
	}
}

func (c Client) Get(url string) (*http.Response, error) {
	var tokenResp TokenResponse
	r, err := c.http.Get(url)
	if err != nil {
		return nil, err
	}
	if r.StatusCode == 200 {
		return r, err
	} else if r.StatusCode == 401 {
		derr := json.NewDecoder(r.Body).Decode(&tokenResp)
		switch {
		case derr == io.EOF:
			return nil, errors.New("no token was given")
		case derr != nil:
			return nil, derr
		}
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenResp.Token))
	resp, err := c.http.Do(req)
	return resp, err
}
