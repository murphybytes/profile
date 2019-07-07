// Package config contains project wide configuration settings that are read from environment variables.
package config

import (
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func defaultSignalValues() *Settings {
	return &Settings{
		ProfileDirectory: func() string {
			dir, _ := os.Getwd()
			return dir
		}(),
		HeapProfilerSignal:         Signal(SIGUSR1),
		HeapProfileName:            "heap.profile",
		GoroutineProfilerSignal:    Signal(SIGUSR1),
		GoroutineProfileName:       "goroutine.profile",
		AllocsProfilerSignal:       Signal(SIGUSR1),
		AllocsProfileName:          "allocs.profile",
		ThreadCreateProfilerSignal: Signal(SIGUSR1),
		ThreadCreateProfileName:    "threadcreate.profile",
		BlockProfilerSignal:        Signal(SIGUSR1),
		BlockProfileName:           "block.profile",
		MutexProfilerSignal:        Signal(SIGUSR1),
		MutexProfileName:           "mutex.profile",
		CPUProfilerSignal:          Signal(SIGUSR1),
		CPUProfileName:             "cpu.profile",
		CPUProfileDuration:         30 * time.Second,
		TraceProfileName:           "trace.profile",
		TraceProfilerSignal:        Signal(SIGUSR1),
		TraceProfileDuration:       time.Second,
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		want    *Settings
		wantErr bool
		setup   func()
	}{
		{
			name:    "defaults",
			want:    defaultSignalValues(),
			wantErr: false,
			setup:   nil,
		},
		{
			name:    "nonnumeric signal",
			want:    nil,
			wantErr: true,
			setup: func() {
				_ = os.Setenv("HEAP_PROFILER_SIGNAL", "fifty")
			},
		},
		{
			name: "valid string signal",
			want: func() *Settings {
				s := defaultSignalValues()
				s.HeapProfilerSignal = Signal(SIGUSR2)
				return s
			}(),
			wantErr: false,
			setup: func() {
				_ = os.Setenv("HEAP_PROFILER_SIGNAL", "SIGUSR2")
			},
		},
		{
			name: "valid numeric signal",
			want: func() *Settings {
				s := defaultSignalValues()
				s.HeapProfilerSignal = Signal(SIGUSR2)
				return s
			}(),
			wantErr: false,
			setup: func() {
				_ = os.Setenv("HEAP_PROFILER_SIGNAL", strconv.Itoa(SIGUSR2))
			},
		},
		{
			name:    "invalid signal",
			want:    nil,
			wantErr: true,
			setup: func() {
				_ = os.Setenv("HEAP_PROFILER_SIGNAL", "90")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			got, err := New()
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
