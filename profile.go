// Package profile contains function that will perform pprof profile operations.
package profile

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"syscall"
	"time"

	"github.com/murphybytes/profile/config"
)

// Mutex stack traces of holders of contended mutexes
func Mutex(ctx context.Context, settings *config.Settings) {
	profile(ctx, settings.MutexProfilerSignal, func() error {
		outputFile := filepath.Join(settings.ProfileDirectory, settings.MutexProfileName)
		return generate("mutex", outputFile)
	})
}

// Block stack traces that led to blocking on synchronization primitives
func Block(ctx context.Context, settings *config.Settings) {
	profile(ctx, settings.BlockProfilerSignal, func() error {
		outputFile := filepath.Join(settings.ProfileDirectory, settings.BlockProfileName)
		return generate("block", outputFile)
	})
}

// ThreadCreate outputs stack traces that led to the creation of new OS threads
func ThreadCreate(ctx context.Context, settings *config.Settings) {
	profile(ctx, settings.ThreadCreateProfilerSignal, func() error {
		outputFile := filepath.Join(settings.ProfileDirectory, settings.ThreadCreateProfileName)
		return generate("threadcreate", outputFile)
	})
}

// Allocs output a profile containing memory allocations
func Allocs(ctx context.Context, settings *config.Settings) {
	profile(ctx, settings.AllocsProfilerSignal, func() error {
		outputFile := filepath.Join(settings.ProfileDirectory, settings.AllocsProfileName)
		return generate("allocs", outputFile)
	})
}

// Goroutine outputs a goroutine profile.
func Goroutine(ctx context.Context, settings *config.Settings) {
	profile(ctx, settings.GoroutineProfilerSignal, func() error {
		outputFile := filepath.Join(settings.ProfileDirectory, settings.GoroutineProfileName)
		return generate("goroutine", outputFile)
	})
}

// Heap outputs a heap profile.
func Heap(ctx context.Context, settings *config.Settings) {
	profile(ctx, settings.HeapProfilerSignal, func() error {
		outputFile := filepath.Join(settings.ProfileDirectory, settings.HeapProfileName)
		return generate("heap", outputFile)
	})
}

// CPU generates a profile when the designated signal is received.
func CPU(ctx context.Context, settings *config.Settings) {
	profile(ctx, settings.CPUProfilerSignal, func() error {
		outputFile := filepath.Join(settings.ProfileDirectory, settings.CPUProfileName)
		return generateCPUProfile(ctx, outputFile, settings.CPUProfileDuration)
	})
}

func generateCPUProfile(ctx context.Context, fileName string, pause time.Duration) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	if err = pprof.StartCPUProfile(f); err != nil {
		return err
	}
	sleep(ctx, pause)
	pprof.StopCPUProfile()
	return nil
}

func sleep(ctx context.Context, duration time.Duration) {
	select {
	case <-time.After(duration):
	case <-ctx.Done():
	}
}

// Trace generates output for the Go execution tracer
func Trace(ctx context.Context, settings *config.Settings) {
	profile(ctx, settings.TraceProfilerSignal, func() error {
		outputFile := filepath.Join(settings.ProfileDirectory, settings.TraceProfileName)
		return generateTrace(ctx, outputFile, settings.TraceProfileDuration)
	})
}

func generateTrace(ctx context.Context, fileName string, pause time.Duration) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	if err = trace.Start(f); err != nil {
		return nil
	}
	sleep(ctx, pause)
	trace.Stop()
	return nil
}

// Run asynchronously executes a profile function.
func Run(profileFunc func(ctx context.Context, cfg *config.Settings)) {
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
		defer func() {
			signal.Stop(ch)
			close(ch)
		}()
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			<-ch
			cancel()
		}()

		cfg, err := config.New()
		if err != nil {
			log.Printf("unable to fetch profiler configuration %q", err)
		}
		profileFunc(ctx, cfg)
		<-ch
	}()
}

func profile(ctx context.Context, sig config.Signal, fn func() error) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.Signal(sig))

	go func() {
		defer func() {
			signal.Stop(ch)
			close(ch)
		}()

		for {
			select {
			case <-ch:
				if err := fn(); err != nil {
					log.Println("profile failed with", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func generate(profileName, fileName string) error {
	if profileName == "heap" {
		runtime.GC()
	}
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	prof := pprof.Lookup(profileName)
	if prof == nil {
		return fmt.Errorf("could not create %q profile, no such profiler exists", profileName)
	}
	if err = prof.WriteTo(f, 0); err != nil {
		return err
	}
	return nil
}
