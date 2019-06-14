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
)

var _ = Describe("osreporter", func() {

	Context("setting paths", func() {
		It("sets the os-report path", func() {
			date := time.Unix(12345, 0)
			r := osreporter.New("/var/vcap/data/tmp", "my-name", date)
			Expect(r.ReportPath).To(Equal("/var/vcap/data/tmp/os-report-my-name-" + date.Format("2006-01-02-15-04-05")))
		})

		It("sets the tarball path", func() {
			date := time.Unix(12345, 0)
			r := osreporter.New("/var/vcap/data/tmp", "my-name", date)
			Expect(r.TarballPath).To(Equal("/var/vcap/data/tmp/os-report-my-name-" + date.Format("2006-01-02-15-04-05") + ".tar.gz"))
		})
	})

	Context("log files", func() {
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

		Context("date.log", func() {
			It("writes the date into a date.log file in the os-report dir", func() {
				runner := osreporter.New(baseReportDir, hostname, timestamp)
				err := runner.Run()
				Expect(err).NotTo(HaveOccurred())

				extractTarball()

				dateLog := filepath.Join(getExtractOsReportPath(), "date.log")
				Expect(dateLog).To(BeAnExistingFile())
				contents, err := ioutil.ReadFile(dateLog)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(contents)).To(Equal(timestamp.Format(time.UnixDate)))
			})
		})
	})
})
