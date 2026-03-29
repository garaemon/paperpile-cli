package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestGetSession_emptyByDefault(t *testing.T) {
	viper.Reset()
	session := GetSession()
	if session != "" {
		t.Errorf("GetSession() = %q, want empty string", session)
	}
}

func TestSaveAndGetSession(t *testing.T) {
	viper.Reset()

	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, "paperpile-cli")

	// Override HOME so resolveConfigDir points to our temp dir.
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Reconfigure viper to use the temp config directory.
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileType)
	viper.AddConfigPath(configDir)

	err := SaveSession("test-session-123")
	if err != nil {
		t.Fatalf("SaveSession() error: %v", err)
	}

	got := GetSession()
	if got != "test-session-123" {
		t.Errorf("GetSession() = %q, want %q", got, "test-session-123")
	}
}

func TestLoad_noConfigFile(t *testing.T) {
	viper.Reset()

	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
}
