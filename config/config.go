// Package config contains project wide configuration settings that are read from environment variables.
package config

import (
	"os"
	"strconv"

	"github.com/joeshaw/envdecode"
)

type Signal int

func (s Signal) Signal() {}

func (s Signal) String() string {
	return "signal " + strconv.Itoa(int(s))
}

func(s *Signal) Decode(sval string) error {
	i, err := strconv.Atoi(sval)
	if err != nil {
		return err
	}
	*s = Signal(i)
	return nil
}

// Settings contain information either read from environment variables or default values.
type Settings struct {
	// ProfileDirectory the output directory for profile files.
	ProfileDirectory string `env:"PROFILE_DIRECTORY"`
	// HeapProfilerSignal the signal that triggers a heap profile dump. Defaults to SIGRTMAX-14
	HeapProfilerSignal int `env:"HEAP_PROFILER_SIGNAL,strict,default=50"`
	// HeapProfileFileName is the name of the output file for the heap profile
	HeapProfileFileName string `env:"HEAP_PROFILE_FILE_NAME,default=heap.profile"`
}

// New returns a structure containing the configuration for the application.
func New() (*Settings, error) {
	var cfg Settings
	if err := envdecode.Decode(&cfg); err != nil {
		return nil, err
	}
	if cfg.ProfileDirectory == "" {
		if dir, err := os.Getwd(); err != nil {
			return nil, err
		} else {
			cfg.ProfileDirectory = dir
		}
	}

	return &cfg, nil
}
