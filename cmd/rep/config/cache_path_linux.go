//go:build linux

package config

import (
	"os"
	"path/filepath"
	"syscall"
)

const (
	// PreferredCachePath is the mount point used for pre-warmed rep cache (e.g. persistent disk mounted by pre-start).
	// If this path exists, is a directory, and is a mount point, rep uses it; otherwise rep uses the configured default.
	PreferredCachePath = "/mnt/rep_cache"
)

// ResolveCachePath returns the cache path to use: PreferredCachePath if it exists, is a directory, and is a
// mount point (i.e. the pre-warmed disk is mounted); otherwise returns fallback (the configured default).
func ResolveCachePath(fallback string) string {
	if !isMountPoint(PreferredCachePath) {
		return fallback
	}
	return PreferredCachePath
}

func isMountPoint(path string) bool {
	fi, err := os.Stat(path)
	if err != nil || !fi.IsDir() {
		return false
	}
	parent := filepath.Dir(path)
	pathStat, err := statDev(path)
	if err != nil {
		return false
	}
	parentStat, err := statDev(parent)
	if err != nil {
		return false
	}
	// Different device => path is a mount point
	return pathStat != parentStat
}

func statDev(path string) (uint64, error) {
	var st syscall.Stat_t
	if err := syscall.Stat(path, &st); err != nil {
		return 0, err
	}
	return uint64(st.Dev), nil
}
