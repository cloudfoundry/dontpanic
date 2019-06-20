package integration_test

import (
	"io/ioutil"
	"os"
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
	uptimeRegexp = `day.*user.*load average`
)

var _ = Describe("Integration", func() {
	var (
		session *gexec.Session
		cmd     *exec.Cmd
	)

	BeforeEach(func() {
		Expect(os.MkdirAll("/var/vcap/jobs/garden/config", 0755)).To(Succeed())
		Expect(ioutil.WriteFile("/var/vcap/jobs/garden/config/config.ini", []byte("hi"), 0644)).To(Succeed())

		cmd = exec.Command(dontPanicBin)
	})

	JustBeforeEach(func() {
		var err error
		session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit())
	})

	It("produces a report correctly", func() {
		By("succeeding")
		Expect(session.ExitCode()).To(Equal(0))

		By("showing an initial message")
		Expect(session).To(gbytes.Say("<Useful information below, please copy-paste from here>"))

		By("collecting the date")
		tarballShouldContainFile("date.log")
		Expect(string(tarballFileContents("date.log"))).
			To(MatchRegexp(dateRegexp))

		By("collecting the uptime")
		tarballShouldContainFile("uptime.log")
		Expect(string(tarballFileContents("uptime.log"))).
			To(MatchRegexp(uptimeRegexp))

		By("collecting the garden version")
		Expect(session).To(gbytes.Say("## Garden Version"))

		By("collecting the hostname")
		tarballShouldContainFile("hostname.log")

		// tarballShouldContainFile("config.ini")
		// Expect(string(tarballFileContents("config.ini"))).To(Equal("hi"))
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

	When("running in BPM", func() {
		It("fails", func() {
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
