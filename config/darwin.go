// +build darwin freebsd netbsd openbsd

package config

const (
	SIGUSR1 = 30
	SIGUSR2 = 31
)

func validateSignal(sig int) error {
	switch sig {
	case SIGUSR1:
	case SIGUSR2:
	default:
		return errSignal
	}
	return nil
}

var sigMap = map[string]int{
	"SIGUSR1": SIGUSR1,
	"SIGUSR2": SIGUSR2,
}

func convertSignal(sigName string) (int, error) {
	if v, ok := sigMap[sigName]; ok {
		return v, nil
	}
	return 0, errSignal
}
