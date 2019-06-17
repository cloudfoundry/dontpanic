package osreporter_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"code.cloudfoundry.org/dontpanic/osreporter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("runner", func() {
	var (
		date         time.Time
		runner       osreporter.Runner
		baseDir      string
		hostname     string
		outputWriter io.Writer
		dateFormat   string
	)

	getOSReportPath := func() string {
		return fmt.Sprintf("%s/os-report-%s-%s", baseDir, hostname, date.Format(dateFormat))
	}

	getTarballPath := func() string {
		return getOSReportPath() + ".tar.gz"
	}

	BeforeEach(func() {
		// Mon Jun 17 09:12:49 UTC 2019
		date = time.Unix(1560762769, 0)
		baseDir = "/var/vcap/data/tmp"
		hostname = "my-name"
		dateFormat = "2006-01-02-15-04-05"
		outputWriter = GinkgoWriter
	})

	Context("setting paths", func() {
		JustBeforeEach(func() {
			runner = osreporter.New(baseDir, hostname, date, outputWriter)
		})

		It("sets the os-report path", func() {
			expectedReportPath := getOSReportPath()
			Expect(runner.ReportPath).To(Equal(expectedReportPath))
		})

		It("sets the tarball path", func() {
			expectedTarballPath := getTarballPath()
			Expect(runner.TarballPath).To(Equal(expectedTarballPath))
		})
	})

	Context("generating tarball from plugin output", func() {
		var (
			extractDir string
		)

		extractTarball := func() {
			tarball := getTarballPath()
			ExpectWithOffset(1, tarball).To(BeAnExistingFile())
			cmd := exec.Command("tar", "xf", tarball, "-C", extractDir)
			ExpectWithOffset(1, cmd.Run()).To(Succeed())
		}

		getExtractOsReportPath := func() string {
			return filepath.Join(extractDir, fmt.Sprintf("os-report-%s-%s", hostname, date.Format(dateFormat)))
		}

		BeforeEach(func() {
			var err error
			baseDir, err = ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())

			extractDir, err = ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())

			randSecs := rand.Int63n(10000000)
			date = time.Unix(randSecs, 0)
			outputWriter = gbytes.NewBuffer()
		})

		AfterEach(func() {
			os.RemoveAll(baseDir)
			os.RemoveAll(extractDir)
		})

		JustBeforeEach(func() {
			runner = osreporter.New(baseDir, hostname, date, outputWriter)
		})

		It("writes a log file into the zipped os-report dir", func() {
			plugin := func() ([]byte, error) {
				return []byte("hello world"), nil
			}
			runner.RegisterStream("hello", "hello.log", plugin)
			err := runner.Run()
			Expect(err).NotTo(HaveOccurred())

			extractTarball()

			logPath := filepath.Join(getExtractOsReportPath(), "hello.log")
			Expect(logPath).To(BeAnExistingFile())
			contents, err := ioutil.ReadFile(logPath)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(contents)).To(Equal("hello world"))
		})

		When("running a plugin", func() {
			BeforeEach(func() {
				plugin := func() ([]byte, error) {
					return []byte("hello world"), nil
				}
				runner.RegisterStream("hello", "hello.log", plugin)
			})

			It("notifies start of operation", func() {
				Expect(runner.Run()).To(Succeed())
				Expect(outputWriter).To(gbytes.Say("## hello\n"))
			})

		})

		When("a plugin returns an error", func() {
			BeforeEach(func() {
				plugin := func() ([]byte, error) {
					return nil, fmt.Errorf("foo")
				}
				runner.RegisterStream("hello", "hello.log", plugin)
			})

			It("notifies failure", func() {
				Expect(runner.Run()).To(Succeed())
				Expect(outputWriter).To(gbytes.Say("Failure: foo\n"))
			})
		})
	})
})
