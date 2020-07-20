package health

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type CheckClient struct {
	config Config
	checks []Check
	logger *zap.Logger
}

type Check interface {
	Check(context.Context) (bool, error)
	ErrorMessage() string
}

func NewCheckClient(c Config, logger *zap.Logger) *CheckClient {
	return &CheckClient{config: c, logger: logger}
}

type EndpointCheck struct {
	endpoint string
	enabled  bool
	name     string
	err      string
}

func (e *EndpointCheck) Check(ctx context.Context) (bool, error) {
	resp, err := http.Get(e.endpoint)
	if err != nil {
		return false, err
	}
	if resp.Status == "200" {
		return true, err
	}
	e.err = fmt.Sprintf("Service %s returned status code %s at endpoint %s", e.name, resp.Status, e.endpoint)
	return false, err
}

func (e *EndpointCheck) ErrorMessage() string {
	return e.err
}

func (c CheckClient) runChecks(ctx context.Context) {
	for _, check := range c.checks {
		b, err := check.Check(ctx)
		if err != nil {
			c.logger.Error(err.Error())
		}
		if !b {
			c.logger.Error(check.ErrorMessage())
		}
	}
}

func (c CheckClient) buildChecks() []Check {
	var checks []Check
	for _, service := range c.config.Services {
		checks = append(checks, &EndpointCheck{service.Endpoint, service.Enabled, service.Name, ""})
	}
	return checks
}

func (c *CheckClient) Start(ctx context.Context) {
	c.checks = c.buildChecks()
	go func() {
		for {
			select {
			case <-ctx.Done():
				c.logger.Info("exiting poll")
				return
			default:
			}

			c.runChecks(ctx)
			time.Sleep(time.Duration(c.config.Options.PollTimeSeconds) * time.Second)
		}
	}()
}
