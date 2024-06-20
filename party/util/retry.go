package util

import (
	"math/rand"
	"time"
)

// Retry ...
// NOTE: sleep can not be type int
func Retry(attempts int, sleep time.Duration, f func() error) error {
	err := f()
	if err == nil {
		return nil
	}
	if attempts--; attempts > 0 {
		jitter := time.Duration(rand.Int63n(int64(sleep)))
		sleep = sleep + jitter/2
		time.Sleep(sleep)
		return Retry(attempts, 2*sleep, f)
	}
	return err
}
