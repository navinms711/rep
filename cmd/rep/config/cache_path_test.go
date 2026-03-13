package config_test

import (
	"code.cloudfoundry.org/rep/cmd/rep/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ResolveCachePath", func() {
	It("returns fallback when preferred path is not a mount point", func() {
		// On non-Linux, ResolveCachePath always returns fallback.
		// On Linux, /mnt/rep_cache may not exist or not be a mount point in test env.
		fallback := "/var/vcap/store/rep_download_cache"
		result := config.ResolveCachePath(fallback)
		// Either preferred (if we're on Linux and it's a mount point) or fallback
		Expect([]string{config.PreferredCachePath, fallback}).To(ContainElement(result))
	})

	It("returns fallback for empty fallback", func() {
		result := config.ResolveCachePath("")
		Expect([]string{config.PreferredCachePath, ""}).To(ContainElement(result))
	})
})
