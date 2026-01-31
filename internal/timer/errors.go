package timer

import "errors"

var (
	ErrTimerAlreadyRunning = errors.New("timer is already running")
	ErrTimerNotRunning     = errors.New("timer is not running")
	ErrTimerNotStopped     = errors.New("timer is not stopped")
)
