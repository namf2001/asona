package pg

import (
	"time"

	"github.com/cenkalti/backoff/v4"
)

// ExponentialBackOff returns the exponential backoff configuration
// Reference: https://docs.microsoft.com/en-us/azure/postgresql/concepts-connectivity#handling-transient-errors
func ExponentialBackOff(maxRetries uint64, maxElapsedTime time.Duration) backoff.BackOff {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = 5 * time.Second
	b.RandomizationFactor = 0
	b.MaxElapsedTime = maxElapsedTime

	return backoff.WithMaxRetries(b, maxRetries)
}
