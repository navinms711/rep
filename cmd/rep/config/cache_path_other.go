//go:build !linux

package config

// PreferredCachePath is the mount point used for pre-warmed rep cache on Linux.
// On non-Linux platforms this path is not used; ResolveCachePath always returns the fallback.
const PreferredCachePath = "/mnt/rep_cache"

// ResolveCachePath returns fallback on non-Linux (no preferred mount point check).
func ResolveCachePath(fallback string) string {
	return fallback
}
