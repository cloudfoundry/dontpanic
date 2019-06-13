package bpmchecker_test

import (
	"fmt"
	"io/ioutil"
	"os"

	"code.cloudfoundry.org/dontpanic/bpmchecker"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bpm Checker", func() {
	var (
		bpmChecker bpmchecker.Checker
	)

	When("the default constructor is used", func() {
		BeforeEach(func() {
			bpmChecker = bpmchecker.New()
		})

		It("looks for the magic text in /proc/1/cmdline", func() {
			Expect(bpmChecker.File).To(Equal("/proc/1/cmdline"))
		})
	})

	Describe("HasGardenPid1", func() {
		var (
			testFile string
		)

		JustBeforeEach(func() {
			bpmChecker = bpmchecker.New(testFile)
		})

		AfterEach(func() {
			Expect(os.RemoveAll(testFile)).To(Succeed())
		})

		When("file contains 'garden_start'", func() {
			BeforeEach(func() {
				testFile = createAndWriteTmpFile("garden_start")
			})

			It("returns true", func() {
				ok, err := bpmChecker.HasGardenPid1()
				Expect(err).ToNot(HaveOccurred())
				Expect(ok).To(BeTrue())
			})
		})

		When("file does not contain 'garden_start'", func() {
			BeforeEach(func() {
				testFile = createAndWriteTmpFile("init")
			})

			It("returns false", func() {
				ok, err := bpmChecker.HasGardenPid1()
				Expect(err).ToNot(HaveOccurred())
				Expect(ok).To(BeFalse())
			})
		})

		When("file does not exist", func() {
			BeforeEach(func() {
				testFile = "/does/not/exist"
			})

			It("returns an error", func() {
				_, err := bpmChecker.HasGardenPid1()
				Expect(err).To(MatchError(ContainSubstring("no such file")))
			})
		})
	})
})

func createAndWriteTmpFile(text string) string {
	testFile, err := ioutil.TempFile("", "")
	Expect(err).NotTo(HaveOccurred())
	defer testFile.Close()

	_, err = fmt.Fprint(testFile, text)
	Expect(err).NotTo(HaveOccurred())
	return testFile.Name()
}
