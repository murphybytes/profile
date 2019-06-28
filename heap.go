package profile

import (
	"context"
	"path/filepath"

	"github.com/murphybytes/profile/config"
)

// Heap outputs a heap profile to a file given a signal. Heap expects a cancel context as an argument. The default
// file is heap.profile in the current working directory.
func Heap(ctx context.Context, settings *config.Settings) {
	outputFile := filepath.Join(settings.ProfileDirectory, settings.HeapProfileName)
	profile(ctx, "heap", outputFile, settings.HeapProfilerSignal)
}
