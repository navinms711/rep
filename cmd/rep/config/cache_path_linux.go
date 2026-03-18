//go:build linux

package config

import (
	"os"
	"path/filepath"
	"syscall"
)

const (
	// PreferredCachePath is the mount point for pre-warmed rep cache (hot-attached disk mounted by pre-start).
	// If this path exists, is a directory, and is a mount point, rep uses it; otherwise rep uses DefaultEphemeralCachePath.
	PreferredCachePath = "/var/vcap/store/rep_download_cache"
	// DefaultEphemeralCachePath is used when PreferredCachePath is not a mount point (no pre-warmed disk).
	DefaultEphemeralCachePath = "/var/vcap/data/rep/shared/garden/download_cache"
)

// ResolveCachePath returns PreferredCachePath if it is a mount point; otherwise DefaultEphemeralCachePath.
// The fallback argument is ignored on Linux so rep always uses ephemeral when the store path is not mounted.
func ResolveCachePath(fallback string) string {
	if !isMountPoint(PreferredCachePath) {
		return DefaultEphemeralCachePath
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
