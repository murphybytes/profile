// Package config contains project wide configuration settings that are read from environment variables.
package config

import (
	"os"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		want    *Settings
		wantErr bool
		setup func()
	}{
		{
			name:    "defaults",
			want:    &Settings{
				ProfileDirectory:    func() string {
					dir, _ := os.Getwd()
					return dir
				}(),
				HeapProfilerSignal:  50,
				HeapProfileFileName: "heap.profile",
			},
			wantErr: false,
			setup:   nil,
		},
		{
			name:    "invalid signal",
			want:    nil,
			wantErr: true,
			setup:   func() {
				_ = os.Setenv("HEAP_PROFILER_SIGNAL", "fifty")
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
