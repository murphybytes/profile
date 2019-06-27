package profiler

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

func TestHeapProfiler(t *testing.T) {
	wd, _ := os.Getwd()
	testDir := filepath.Join(wd, "testout")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatal("could not create test dir", testDir)
	}
	defer func() {
		os.RemoveAll(testDir)
	}()
	testapp := filepath.Join(testDir, "app")
	cmd := exec.Command("go", "build", "-o", testapp, "github.com/murphybytes/profiler/example/heap")
	if err := cmd.Run(); err != nil {
		t.Fatal("could not build test app")
	}

	var buff bytes.Buffer
	cmd = exec.Command(testapp)
	cmd.Stdout = &buff
	cmd.Stderr = cmd.Stdout

	cmd.Env = []string{
		fmt.Sprintf("PROFILE_DIRECTORY=%s", testDir),
		"HEAP_PROFILER_SIGNAL=30",
	}
	if err := cmd.Start(); err != nil {
		t.Fatal("could not start test app", err)
	}
	time.Sleep(10*time.Millisecond)
	if err := cmd.Process.Signal(syscall.Signal(30)); err != nil {
		t.Fatal("attempt to send signal failed")
	}

	time.Sleep(time.Millisecond)
	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		t.Fatal("could not send interrupt signal")
	}
	_ = cmd.Wait()

	if _, err := os.Stat(filepath.Join(testDir, "heap.profile")); os.IsNotExist(err) {
		t.Fatal("profile output is not present")
	}
}