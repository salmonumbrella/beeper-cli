package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestConfigDir(t *testing.T) {
	dir, err := ConfigDir()
	if err != nil {
		t.Fatalf("ConfigDir() error: %v", err)
	}
	if !strings.Contains(dir, "beeper-cli") {
		t.Errorf("ConfigDir() = %q, want path containing 'beeper-cli'", dir)
	}
}

func TestConfigDirWithXDG(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("XDG not used on Windows")
	}

	tmpDir := t.TempDir()
	_ = os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer func() { _ = os.Unsetenv("XDG_CONFIG_HOME") }()

	dir, err := ConfigDir()
	if err != nil {
		t.Fatalf("ConfigDir() error: %v", err)
	}
	expected := filepath.Join(tmpDir, "beeper-cli")
	if dir != expected {
		t.Errorf("ConfigDir() = %q, want %q", dir, expected)
	}
}
