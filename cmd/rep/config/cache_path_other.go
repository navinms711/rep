//go:build !linux

package config

// PreferredCachePath is the mount point for pre-warmed rep cache on Linux.
const PreferredCachePath = "/var/vcap/store/rep_download_cache"

// DefaultEphemeralCachePath is the fallback cache path when preferred is not mounted (Linux only).
const DefaultEphemeralCachePath = "/var/vcap/data/rep/shared/garden/download_cache"

// ResolveCachePath returns fallback on non-Linux (no preferred mount point check).
func ResolveCachePath(fallback string) string {
	return fallback
}
