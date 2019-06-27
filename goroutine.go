package profiler

import (
	"context"
	"path/filepath"

	"github.com/murphybytes/profiler/config"
)


// Goroutine outputs a goroutine profile
func Goroutine(ctx context.Context, settings *config.Settings) {
	outputFile := filepath.Join(settings.ProfileDirectory, settings.GoroutineProfileName)
	profile(ctx, "goroutine", outputFile, settings.GoroutineProfilerSignal)
}

