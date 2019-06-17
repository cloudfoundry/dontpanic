package date_test

import (
	"context"

	"code.cloudfoundry.org/dontpanic/plugins/date"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Date", func() {
	It("returns date bytes", func() {
		out, err := date.Run(context.TODO())
		Expect(err).NotTo(HaveOccurred())
		Expect(string(out)).To(MatchRegexp(`\w{3} \w{3} \d{1,2}.*\d{4}.*`))
	})
})
