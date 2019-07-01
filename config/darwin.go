// +build darwin freebsd netbsd openbsd

package config

import "fmt"

const (
	SIGUSR1 = 30
	SIGUSR2 = 31
)

func validateSignal(sig int) error {
	switch sig {
	case SIGUSR1:
	case SIGUSR2:
	default:
		return &ErrSignal{msg: fmt.Sprintf("signal %d is not usable for bsd", sig) }
	}
	return nil
}


