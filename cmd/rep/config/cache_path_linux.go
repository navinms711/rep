//go:build linux

package config

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

const (
	// PreferredCachePath is the mount point for pre-warmed rep cache (hot-attached disk mounted by pre-start).
	PreferredCachePath = "/var/vcap/store/rep_download_cache"
	// DefaultEphemeralCachePath is used when PreferredCachePath is not a mount point (no pre-warmed disk).
	DefaultEphemeralCachePath = "/var/vcap/data/rep/shared/garden/download_cache"
	// repCachePathFile is written by mount_rep_cache (bpm-pre-start) only after a successful
	// disk mount.  rep reads it from inside the BPM container where /var/vcap/jobs/rep/config
	// is bind-mounted read-only.  This avoids relying solely on isMountPoint(), which cannot
	// distinguish a freshly attached disk from the same empty mkdir that was created on the
	// root filesystem before the disk was available.
	repCachePathFile = "/var/vcap/jobs/rep/config/rep_cache_path"
)

// ResolveCachePath returns the cache path to use.
//
// Decision logic (in order):
//  1. Read repCachePathFile written by bpm-pre-start on successful disk mount.
//  2. Validate the path exists, is a directory, and is a real mount point (device differs
//     from its parent) — guards against a stale flag file left over from a prior boot.
//  3. Fall back to DefaultEphemeralCachePath when either check fails.
//
// The fallback argument is ignored on Linux.
func ResolveCachePath(fallback string) string {
	data, err := os.ReadFile(repCachePathFile)
	if err != nil {
		// Flag file absent: disk was not mounted (no callback, attach failed, etc.).
		return DefaultEphemeralCachePath
	}
	path := strings.TrimSpace(string(data))
	if path == "" {
		return DefaultEphemeralCachePath
	}
	// Secondary guard: confirm the path is still an actual mount point inside the BPM
	// container.  After /var/vcap/store/rep_download_cache is added to additional_volumes
	// in bpm.yml the bind-mount propagates sdc into the container, so isMountPoint() will
	// return true (device 8:32) vs its parent /var/vcap/store (device 8:18).
	if !isMountPoint(path) {
		return DefaultEphemeralCachePath
	}
	return path
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
