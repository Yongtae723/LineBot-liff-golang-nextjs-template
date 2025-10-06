package llm

import (
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type RetryConfig struct {
	MaxRetries  int
	WaitTime    time.Duration
	MaxWaitTime time.Duration
}

func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:  3,
		WaitTime:    1 * time.Second,
		MaxWaitTime: 10 * time.Second,
	}
}

func WrapClient(client *http.Client) *retryablehttp.Client {
	return WrapClientWithConfig(client, DefaultRetryConfig())
}

func WrapClientWithConfig(client *http.Client, config *RetryConfig) *retryablehttp.Client {
	retryClient := retryablehttp.NewClient()

	retryClient.RetryMax = config.MaxRetries
	retryClient.RetryWaitMin = config.WaitTime
	retryClient.RetryWaitMax = config.MaxWaitTime

	retryClient.HTTPClient = client
	return retryClient
}
