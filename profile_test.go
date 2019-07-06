package profile

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/murphybytes/profile/config"
)

func makeTestDir() (string, func(), error) {
	wd, _ := os.Getwd()
	testDir := filepath.Join(wd, "testout")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		return "", nil, fmt.Errorf("could not create test dir %q", testDir)
	}
	deleter := func() {
		os.RemoveAll(testDir)
	}
	return testDir, deleter, nil
}

func TestProfiler(t *testing.T) {
	profiles := []string{
		"heap",
		"goroutine",
		"allocs",
		"threadcreate",
		"block",
		"mutex",
	}

	for _, profile := range profiles {
		t.Run(profile, func(t *testing.T) {
			testDir, deleter, err := makeTestDir()
			if err != nil {
				t.Fatal(err)
			}
			defer deleter()
			testapp := filepath.Join(testDir, "app")
			cmd := exec.Command("go", "build", "-o", testapp, "github.com/murphybytes/profile/example")
			if err := cmd.Run(); err != nil {
				t.Fatal("could not build test app")
			}

			var buff bytes.Buffer
			cmd = exec.Command(testapp)
			cmd.Stdout = &buff
			cmd.Stderr = cmd.Stdout

			cmd.Env = []string{
				fmt.Sprintf("PROFILE_DIRECTORY=%s", testDir),
				fmt.Sprintf("%s_PROFILER_SIGNAL=%d", strings.ToUpper(profile), config.SIGUSR2),
			}
			if err := cmd.Start(); err != nil {
				t.Fatal("could not start test app", err)
			}
			time.Sleep(10 * time.Millisecond)
			if err := cmd.Process.Signal(syscall.Signal(config.SIGUSR2)); err != nil {
				t.Fatal("attempt to send signal failed")
			}

			time.Sleep(time.Millisecond)
			if err := cmd.Process.Signal(os.Interrupt); err != nil {
				t.Fatal("could not send interrupt signal")
			}
			_ = cmd.Wait()

			if _, err := os.Stat(filepath.Join(testDir, fmt.Sprintf("%s.profile", profile))); os.IsNotExist(err) {
				t.Fatal("profile output is not present")
			}
		})
	}
}

func TestProfileHappyPath(t *testing.T) {
	testDir, deleter, err := makeTestDir()
	if err != nil {
		t.Fatal(err)
	}
	defer deleter()
	cfg, _ := config.New()
	ctx, cancel := context.WithCancel(context.Background())
	profile(ctx,
		cfg.HeapProfilerSignal,
		func() error {
			return generate("heap", filepath.Join(testDir, cfg.HeapProfileName))
		})
	time.Sleep(time.Millisecond)
	prc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}
	prc.Signal(syscall.Signal(cfg.HeapProfilerSignal))
	time.Sleep(10 * time.Millisecond)
	cancel()
	if _, err := os.Stat(filepath.Join(testDir, "heap.profile")); os.IsNotExist(err) {
		t.Fatal("profile output is not present")
	}
}
