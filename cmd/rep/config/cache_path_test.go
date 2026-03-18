package config_test

import (
	"code.cloudfoundry.org/rep/cmd/rep/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ResolveCachePath", func() {
	It("returns preferred or default ephemeral path", func() {
		// On Linux: PreferredCachePath if mount point, else DefaultEphemeralCachePath.
		// On non-Linux: returns fallback argument.
		fallback := "/var/vcap/data/rep/shared/garden/download_cache"
		result := config.ResolveCachePath(fallback)
		Expect([]string{config.PreferredCachePath, config.DefaultEphemeralCachePath, fallback}).To(ContainElement(result))
	})

	It("on non-Linux returns empty fallback when given", func() {
		result := config.ResolveCachePath("")
		// On Linux we get DefaultEphemeralCachePath; on non-Linux we get "".
		Expect([]string{config.PreferredCachePath, config.DefaultEphemeralCachePath, ""}).To(ContainElement(result))
	})
})
