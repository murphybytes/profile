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

var sigMap = map[string]int{
	"SIGUSR1": SIGUSR1,
	"SIGUSR2": SIGUSR2,
}

func convertSignal(sigName string)(int, error) {
	if v, ok := sigMap[sigName]; ok {
		return v, nil
	}
	return 0, &ErrSignal{msg: fmt.Sprintf("unknown signal %q", sigName) }
}


