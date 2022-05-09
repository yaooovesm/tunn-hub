package timer

import (
	"errors"
	"time"
)

var (
	Timeout = errors.New("timeout")
)

//
// TimeoutTask
// @Description:
// @param f
// @param duration
// @return error
//
func TimeoutTask(f func() error, duration time.Duration) error {
	ch := make(chan error)
	go func() {
		ch <- f()
	}()
	select {
	case <-time.After(duration):
		return Timeout
	case err := <-ch:
		return err
	}
}
