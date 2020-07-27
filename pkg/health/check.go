package health

import (
	"context"
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

func (c CheckClient) runChecks(ctx context.Context) {
	for _, check := range c.checks {
		b, err := check.Check(ctx)
		if err != nil {
			c.logger.Error(err.Error())
		} else if !b {
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
