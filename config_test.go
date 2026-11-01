package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigUsesWriteTimeoutEnvVar(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(workingDir); err != nil {
			t.Fatalf("restore cwd: %v", err)
		}
	})

	configPath := filepath.Join(tempDir, ConfigFile)
	config := `{
		"TLS": {"CERTPATH": "", "KEYPATH": "", "KEYPASSWORD": "", "CAPATH": ""},
		"UIMANAGER": {"MODULЕID": "", "SPACEID": "", "DATAMANAGERURL": []},
		"HTTP": {"HOST": "0.0.0.0", "PORT": "7007", "READTIMEOUT": 10, "WRITETIMEOUT": 20, "RL": 5, "TLS": false, "SKIPVERIFY": false, "TRUSTPROXY": false},
		"JWT": {"SECRET": "12345678901234567890123456789012", "ALLOWINSECUREHTTP": false},
		"DB": {"HOSTS": ["127.0.0.1:9042"], "USERNAME": "test", "PASSWORD": "secret", "TLS": false, "SKIPVERIFY": false, "KEYSPACE": "test", "CONSISTENCY": "Quorum", "CONSISTENSYREAD": "One", "TIMEOUT": 1000}
	}`
	if err := os.WriteFile(configPath, []byte(config), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	t.Setenv("MONITOR_HTTP_READTIMEOUT", "12")
	t.Setenv("MONITOR_HTTP_WRITETIMEOUT", "45")

	var appCfg TAppConfig
	if err := appCfg.LoadConfig(); err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if appCfg.HTTP.READTIMEOUT != 12 {
		t.Fatalf("HTTP.READTIMEOUT = %d, want 12", appCfg.HTTP.READTIMEOUT)
	}
	if appCfg.HTTP.WRITETIMEOUT != 45 {
		t.Fatalf("HTTP.WRITETIMEOUT = %d, want 45", appCfg.HTTP.WRITETIMEOUT)
	}
}
