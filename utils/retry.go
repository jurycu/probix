package utils

import (
	"errors"
	"time"

	"k8s.io/klog"
)

func Retry(operation func() error, name string, attempts int, retryWaitSeconds int) (err error) {
	return RetryInc(operation, name, attempts, retryWaitSeconds, 0)
}

func RetryInc(operation func() error, name string, attempts int, retryWaitSeconds int, retryWaitIncSeconds int) (err error) {
	for i := 0; ; i++ {
		err = operation()
		if err == nil {
			if i > 0 {
				klog.Infof("retry #%d %v finally succeed", i, name)
			}
			return nil
		}
		klog.Errorf("retry #%d %v, error: %s", i, name, err)

		if i >= (attempts - 1) {
			break
		}

		time.Sleep(time.Second * time.Duration(retryWaitSeconds))
		retryWaitSeconds = retryWaitSeconds + retryWaitIncSeconds
	}
	return err
}

// MaxRetries is the maximum number of retries before bailing.
var MaxRetries = 10

var errMaxRetriesReached = errors.New("exceeded retry limit")

// Func represents functions that can be retried.
type Func func(attempt int) (retry bool, err error)

// Do keeps trying the function until the second argument
// returns false, or no error is returned.
func TryFunc(fn Func) error {
	var err error
	var cont bool
	attempt := 1
	for {
		cont, err = fn(attempt)
		if !cont || err == nil {
			break
		}
		attempt++
		if attempt > MaxRetries {
			return errMaxRetriesReached
		}
	}
	return err
}
