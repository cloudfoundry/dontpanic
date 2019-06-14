package osreporter_test

import (
	"io/ioutil"
	"os"
	"time"

	"code.cloudfoundry.org/dontpanic/osreporter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Plugins", func() {
	var (
		tmpDir string
	)

	BeforeEach(func() {
		var err error
		tmpDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(os.RemoveAll(tmpDir)).To(Succeed())
	})

	It("will run a registered stream plugin", func() {
		run := false
		plug := func() ([]byte, error) {
			run = true
			return nil, nil
		}

		r := osreporter.New(tmpDir, "hostname", time.Now(), GinkgoWriter)
		r.RegisterStream("some name", "out.file", plug)
		r.Run()

		Expect(run).To(BeTrue())
	})
})
