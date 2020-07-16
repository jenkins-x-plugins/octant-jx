package files

import (
	"os"
	"path/filepath"
)

func HomeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	h := os.Getenv("USERPROFILE") // windows
	if h == "" {
		h = "."
	}
	return h
}

func JXOPSHomeDir() string {
	homeDir := os.Getenv("JX_OPS_HOME")
	if homeDir == "" {
		homeDir = filepath.Join(HomeDir(), ".jx-ops")
	}
	return homeDir
}
