package uptime_test

import (
	"context"

	"code.cloudfoundry.org/dontpanic/plugins/uptime"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Uptime", func() {
	It("returns `uptime` output", func() {
		bytes, err := uptime.Run(context.TODO())
		Expect(err).NotTo(HaveOccurred())
		Expect(string(bytes)).To(MatchRegexp(`\d{2}:\d{2}:\d{2} up \d+ days, \d{2}:\d{2},\s+\d users?,\s+load average: \d\.\d\d, \d\.\d\d, \d\.\d\d`))
	})
})
