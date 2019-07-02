// +build linux

package config

import "fmt"

const (
	SIGUSR1 = 10
	SIGUSR2 = 12
	SIGRTMIN = 34
	SIGRTMIN1 = SIGRTMIN + 1
	SIGRTMIN2 = SIGRTMIN + 2
	SIGRTMIN15 = SIGRTMIN + 15
	SIGRTMAX15 = SIGRTMAX - 14
	SIGRTMAX1 = SIGRTMAX - 1 
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


var sigMap = map[string]int{
	"SIGUSR1": SIGUSR1,
	"SIGUSR2": SIGUSR2,
	"SIGRTMIN+1": SIGRTMIN1,
}

func convertSignal(sigName string)(int, error) {
	if v, ok := sigMap[sigName]; ok {
		return v, nil
	}
	return 0, &ErrSignal{msg: fmt.Sprintf("unknown signal %q", sigName) }
}