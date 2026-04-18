package healthcheck

import (
	"fmt"
	"net/http"
	"time"
)

// Status represents the result of a Vault health check.
type Status struct {
	Reachable   bool
	Sealed      bool
	StatusCode  int
	Latency     time.Duration
	Error       string
}

// Checker performs health checks against a Vault instance.
type Checker struct {
	client  *http.Client
	baseURL string
}

// New returns a Checker targeting the given Vault base URL.
func New(baseURL string, timeout time.Duration) *Checker {
	return &Checker{
		client:  &http.Client{Timeout: timeout},
		baseURL: baseURL,
	}
}

// Check calls the Vault /v1/sys/health endpoint and returns a Status.
func (c *Checker) Check() Status {
	url := fmt.Sprintf("%s/v1/sys/health", c.baseURL)
	start := time.Now()

	resp, err := c.client.Get(url)
	latency := time.Since(start)

	if err != nil {
		return Status{
			Reachable: false,
			Latency:   latency,
			Error:     err.Error(),
		}
	}
	defer resp.Body.Close()

	s := Status{
		Reachable:  true,
		StatusCode: resp.StatusCode,
		Latency:    latency,
	}

	// Vault returns 200 = initialized, unsealed, active
	// 429 = unsealed, standby; 472/473 = DR/perf standby
	// 501 = not initialized; 503 = sealed
	switch resp.StatusCode {
	case http.StatusOK, 429, 472, 473:
		s.Sealed = false
	case 503:
		s.Sealed = true
		s.Error = "vault is sealed"
	default:
		s.Error = fmt.Sprintf("unexpected status %d", resp.StatusCode)
	}

	return s
}

// String returns a human-readable summary of the status.
func (s Status) String() string {
	if !s.Reachable {
		return fmt.Sprintf("unreachable: %s", s.Error)
	}
	if s.Sealed {
		return fmt.Sprintf("sealed (HTTP %d, latency %s)", s.StatusCode, s.Latency.Round(time.Millisecond))
	}
	return fmt.Sprintf("healthy (HTTP %d, latency %s)", s.StatusCode, s.Latency.Round(time.Millisecond))
}
