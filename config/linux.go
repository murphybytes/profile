// +build linux

package config

import "fmt"

const (
	SIGUSR1 = 10
	SIGUSR2 = 12
	SIGRTMIN = 34
	SIGRTMAX = 64
)

func validateSignal(sig int) error {
	switch sig {
	case SIGUSR1:
	case SIGUSR2:
	case SIGRTMIN <= sig && SIGRTMAX >= sig :
	default:
		return &ErrSignal{msg: fmt.Sprintf("signal %d is not usable for linux", sig)}
	}
	return nil
}


