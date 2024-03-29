// Package config contains project wide configuration settings that are read from environment variables.
package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joeshaw/envdecode"
)

type configError string

func (e configError) Error() string {
	return string(e)
}

const errSignal = configError("signal value not supported")

// Signal handles human friendly signal names in environment variables.
type Signal int

// Signal support for syscall.Signal interface
func (s Signal) Signal() {}

// String support for Stringer
func (s Signal) String() string {
	return "signal " + strconv.Itoa(int(s))
}

// Decode process and validate environment variables. Both signal names such as
// SIGUSR1 and numbers are supported.
func (s *Signal) Decode(sval string) error {
	i, err := strconv.Atoi(sval)
	if err != nil {
		if i, err = convertSignal(sval); err != nil {
			return err
		}
		*s = Signal(i)
		return nil
	}
	if err := validateSignal(i); err != nil {
		return err
	}
	*s = Signal(i)
	return nil
}

// Settings contain information either read from environment variables or default values.
type Settings struct {
	// ProfileDirectory the output directory for profile files.
	ProfileDirectory string `env:"PROFILE_DIRECTORY"`
	// HeapProfilerSignal the signal that triggers a heap profile dump. Defaults to SIGUSR1
	HeapProfilerSignal Signal `env:"HEAP_PROFILER_SIGNAL,strict,default=SIGUSR1"`
	// HeapProfileName is the name of the output file for the heap profile
	HeapProfileName string `env:"HEAP_PROFILE_NAME,default=heap.profile"`
	// GoroutineProfileName is the output file of the goroutine profiler
	GoroutineProfileName string `env:"GOROUTINE_PROFILE_NAME,default=goroutine.profile"`
	// GoroutineProfilerSignal the signal that triggers a heap profile dump. Defaults to SIGUSR1
	GoroutineProfilerSignal Signal `env:"GOROUTINE_PROFILER_SIGNAL,strict,default=SIGUSR1"`
	// AllocsProfileName is the output file of the allocs profiler. Defaults to SIGUSR1
	AllocsProfileName string `env:"ALLOCS_PROFILE_NAME,default=allocs.profile"`
	// AllocsProfilerSignal is the signal that triggers the allocs profiler
	AllocsProfilerSignal Signal `env:"ALLOCS_PROFILER_SIGNAL,strict,default=SIGUSR1"`
	// ThreadCreateProfileName the output file for the thread creation profile
	ThreadCreateProfileName string `env:"THREADCREATE_PROFILE_NAME,default=threadcreate.profile"`
	// ThreadCreateProfilerSignal is the signal the triggers the thread create profiler
	ThreadCreateProfilerSignal Signal `env:"THREADCREATE_PROFILER_SIGNAL,strict,default=SIGUSR1"`
	// BlockProfileName the output file for the block profile
	BlockProfileName string `env:"BLOCK_PROFILE_NAME,default=block.profile"`
	// BlockProfilerSignal is the signal that triggers the generation of a block profile
	BlockProfilerSignal Signal `env:"BLOCK_PROFILER_SIGNAL,strict,default=SIGUSR1"`
	// MutexProfileName name of the output file for the mutex profile output.
	MutexProfileName string `env:"MUTEX_PROFILE_NAME,default=mutex.profile"`
	// MutexProfilerSignal is the signal that triggers the generation of the mutex profile.
	MutexProfilerSignal Signal `env:"MUTEX_PROFILER_SIGNAL,strict,default=SIGUSR1"`
	// CPUProfileName is the name of the output file that will contain CPU profile information
	CPUProfileName string `env:"CPU_PROFILE_NAME,default=cpu.profile"`
	// CPUProfilerSignal is the signal that will trigger the creation of a CPU profile
	CPUProfilerSignal Signal `env:"CPU_PROFILER_SIGNAL,strict,default=SIGUSR1"`
	// CPRProfileDuration the amount of time taken to collect CPU profile
	CPUProfileDuration time.Duration `env:"CPU_PROFILE_DURATION,strict,default=30s"`
	// TraceProfileName is the name of the output file emitted by the trace profiler
	TraceProfileName string `env:"TRACE_PROFILE_NAME,default=trace.profile"`
	// TraceProfilerSignal the signal that will trigger the trace profiler
	TraceProfilerSignal Signal `env:"TRACE_PROFILER_SIGNAL,strict,default=SIGUSR1"`
	// TraceProfileDuration the amount of time taken to collect the trace profile.
	TraceProfileDuration time.Duration `env:"TRACE_PROFILE_DURATION,strict,default=1s"`
}

// New returns a structure containing the configuration for the application.
func New() (*Settings, error) {
	var cfg Settings
	if err := envdecode.Decode(&cfg); err != nil {
		return nil, err
	}
	if cfg.ProfileDirectory == "" {
		dir, err := os.Getwd();
		if err != nil {
			return nil, err
		}
		cfg.ProfileDirectory = dir
	}

	return &cfg, nil
}
