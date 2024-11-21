package health

import (
	"context"
	"sync"

	"github.com/rhoat/go-exercise/pkg/health/checks"
)

type Result string

var (
	OK   Result = "OK"
	Fail Result = "Fail"
)

type Checker interface {
	Check(context.Context) error
	Name() string
}

type Status struct {
	Name   string `json:"name"`
	Status Result `json:"status"`
	Error  string `json:"error,omitempty"`
}

type Registry struct {
	mu     sync.RWMutex
	checks []Checker
}

func NewRegistry() *Registry {
	return &Registry{
		checks: NewHealthChecks(),
	}
}

func (r *Registry) RunChecks(ctx context.Context) []Status {
	r.mu.RLock()
	defer r.mu.RUnlock()

	results := []Status{}
	for _, check := range r.checks {
		err := check.Check(ctx)
		status := OK
		var errorMsg string
		if err != nil {
			status = Fail
			errorMsg = err.Error()
		}
		results = append(results, Status{
			Name:   check.Name(),
			Status: status,
			Error:  errorMsg,
		})
	}
	return results
}

// Constructor for Health Checks - Returns a slice of HealthChecker.
func NewHealthChecks() []Checker {
	var timeout = 1000
	return []Checker{
		checks.NewPingCheck("https://www.google.com", "GET", timeout, nil, nil),
		// failing check
		checks.NewPingCheck("https://wwws.google.com", "GET", timeout, nil, nil),
		// Add other health checks here if needed
	}
}
