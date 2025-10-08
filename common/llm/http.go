package llm

import (
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

// RetryConfig holds configuration for retry behavior
type RetryConfig struct {
	MaxRetries  int
	WaitTime    time.Duration
	MaxWaitTime time.Duration
}

// DefaultRetryConfig returns a default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:  3,
		WaitTime:    1 * time.Second,
		MaxWaitTime: 10 * time.Second,
	}
}

// WrapClient wraps an existing HTTP client with retry functionality
func WrapClient(client *http.Client) *retryablehttp.Client {
	return WrapClientWithConfig(client, DefaultRetryConfig())
}

// WrapClientWithConfig wraps an existing HTTP client with custom retry configuration
func WrapClientWithConfig(client *http.Client, config *RetryConfig) *retryablehttp.Client {
	retryClient := retryablehttp.NewClient()

	// Set retry configuration
	retryClient.RetryMax = config.MaxRetries
	retryClient.RetryWaitMin = config.WaitTime
	retryClient.RetryWaitMax = config.MaxWaitTime

	// Use the provided HTTP client
	retryClient.HTTPClient = client
	return retryClient
}
