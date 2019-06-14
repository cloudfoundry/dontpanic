package osreporter_test

import (
	"fmt"
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

	Context("setting paths", func() {
		It("sets the os-report path", func() {
			date := time.Unix(12345, 0)
			r := osreporter.New("/var/vcap/data/tmp", "my-name", date, GinkgoWriter)
			Expect(r.ReportPath).To(Equal("/var/vcap/data/tmp/os-report-my-name-" + date.Format("2006-01-02-15-04-05")))
		})

		It("sets the tarball path", func() {
			date := time.Unix(12345, 0)
			r := osreporter.New("/var/vcap/data/tmp", "my-name", date, GinkgoWriter)
			Expect(r.TarballPath).To(Equal("/var/vcap/data/tmp/os-report-my-name-" + date.Format("2006-01-02-15-04-05") + ".tar.gz"))
		})
	})

	Context("generating tarball from plugin output", func() {
		var (
			baseReportDir string
			extractDir    string
			hostname      string
			timestamp     time.Time
		)

		getTarballPath := func() string {
			date := timestamp.Format("2006-01-02-15-04-05")
			file := fmt.Sprintf("os-report-%s-%s.tar.gz", hostname, date)
			return filepath.Join(baseReportDir, file)
		}

		getExtractOsReportPath := func() string {
			date := timestamp.Format("2006-01-02-15-04-05")
			file := fmt.Sprintf("os-report-%s-%s", hostname, date)
			return filepath.Join(extractDir, file)
		}

		extractTarball := func() {
			tarball := getTarballPath()
			Expect(tarball).To(BeAnExistingFile())
			cmd := exec.Command("tar", "xf", tarball, "-C", extractDir)
			Expect(cmd.Run()).To(Succeed())
		}

		BeforeEach(func() {
			hostname = "hedge"

			var err error
			baseReportDir, err = ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())

			extractDir, err = ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())

			randSecs := rand.Int63n(10000000)
			timestamp = time.Unix(randSecs, 0)
		})

		AfterEach(func() {
			os.RemoveAll(baseReportDir)
			os.RemoveAll(extractDir)
		})

		It("writes a log file into the zipped os-report dir", func() {
			runner := osreporter.New(baseReportDir, hostname, timestamp, GinkgoWriter)
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

		It("notifies when running a plugin", func() {
			out := gbytes.NewBuffer()
			runner := osreporter.New(baseReportDir, hostname, timestamp, out)
			plugin := func() ([]byte, error) {
				return []byte("hello world"), nil
			}
			runner.RegisterStream("hello", "hello.log", plugin)
			err := runner.Run()
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(gbytes.Say("## hello\n"))
		})

		It("notifies failure when a plugin returns an error", func() {
			out := gbytes.NewBuffer()
			runner := osreporter.New(baseReportDir, hostname, timestamp, out)
			plugin := func() ([]byte, error) {
				return nil, fmt.Errorf("foo")
			}
			runner.RegisterStream("hello", "hello.log", plugin)
			err := runner.Run()
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(gbytes.Say("Failure: foo\n"))
		})
	})
})
