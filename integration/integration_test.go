package integration_test

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

const (
	baseDir      = "/var/vcap/data/tmp"
	dateRegexp   = `\w{3} \w{3} \d{1,2}.*\d{4}.*`
	uptimeRegexp = `\d{2}:\d{2}:\d{2} up \d+ days, .*,\s+\d users?,\s+load average: \d\.\d\d, \d\.\d\d, \d\.\d\d`
)

var _ = Describe("Integration", func() {
	var (
		session *gexec.Session
		cmd     *exec.Cmd
	)

	BeforeEach(func() {
		cmd = exec.Command(dontPanicBin)
	})

	JustBeforeEach(func() {
		var err error
		session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit())
	})

	When("running normally", func() {
		It("runs the binary and shows the initial messages", func() {
			Expect(session.ExitCode()).To(Equal(0))
			Expect(session).To(gbytes.Say("<Useful information below, please copy-paste from here>"))
		})
	})

	When("running as a non-root user", func() {
		BeforeEach(func() {
			cmd.SysProcAttr = &syscall.SysProcAttr{
				Credential: &syscall.Credential{Uid: 5000, Gid: 5000},
			}
		})

		It("warns and exits", func() {
			Expect(session.ExitCode()).ToNot(Equal(0))
			Expect(session.Err).To(gbytes.Say("Keep Calm and Re-run as Root!"))
		})
	})

	It("does not allow execution within a BPM container", func() {
		Skip("return to this")
		// wd, err := os.Getwd()
		// Expect(err).NotTo(HaveOccurred())
		// Expect(os.Symlink(dontPanicBin, wd+"/assets/rootfs/bin/dontPanicBin")).To(Succeed())
		//
		// runcRun := exec.Command("runc", "run", "assets/config.json", "fake-bpm-container")
		// _, err = gexec.Start(runcRun, GinkgoWriter, GinkgoWriter)
		// Expect(err).NotTo(HaveOccurred())
		//
		// runcExec := exec.Command("runc", "exec", "fake-bpm-container", "/bin/dontPanicBin")
		// sess, err := gexec.Start(runcExec, GinkgoWriter, GinkgoWriter)
		// Expect(err).NotTo(HaveOccurred())
		// Eventually(sess).Should(gexec.Exit(1))
	})

	When("running with standard plugins", func() {
		It("produces expected log files", func() {
			Expect(session.ExitCode()).To(Equal(0))

			tarballShouldContainFile("date.log")
			Expect(string(tarballFileContents("date.log"))).
				To(MatchRegexp(dateRegexp))

			tarballShouldContainFile("uptime.log")
			Expect(string(tarballFileContents("uptime.log"))).
				To(MatchRegexp(uptimeRegexp))
		})
	})
})

func tarballShouldContainFile(filePath string) {
	tarball := getTarball()
	ExpectWithOffset(1, tarball).ToNot(BeEmpty(), "tarball not found in "+baseDir)

	extractedOsReportPath := strings.TrimRight(filepath.Base(tarball), ".tar.gz")
	logFilePath := filepath.Join(extractedOsReportPath, filePath)
	ExpectWithOffset(1, listTarball(tarball)).To(ContainSubstring(logFilePath))
}

func tarballFileContents(filePath string) []byte {
	tarball := getTarball()
	ExpectWithOffset(1, tarball).ToNot(BeEmpty(), "tarball not found in "+baseDir)

	extractedOsReportPath := strings.TrimRight(filepath.Base(tarball), ".tar.gz")
	osDir := filepath.Base(extractedOsReportPath)

	cmd := exec.Command("tar", "xf", tarball, filepath.Join(osDir, filePath), "-O")
	out, err := cmd.Output()
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return out
}

func getTarball() string {
	dirEntries, err := ioutil.ReadDir(baseDir)
	ExpectWithOffset(2, err).NotTo(HaveOccurred())

	re := regexp.MustCompile(`os-report-.*\.tar\.gz`)
	for _, info := range dirEntries {
		if info.IsDir() {
			continue
		}
		if re.MatchString(info.Name()) {
			return filepath.Join(baseDir, info.Name())
		}
	}
	return ""
}

func listTarball(tarball string) string {
	cmd := exec.Command("tar", "tf", tarball)
	files, err := cmd.Output()
	ExpectWithOffset(2, err).NotTo(HaveOccurred())
	return string(files)
}
