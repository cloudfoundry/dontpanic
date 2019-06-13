package integration_test

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"syscall"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Integration", func() {
	It("can run the binary", func() {
		cmd := exec.Command(dontPanicBin)
		sess, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess).Should(gexec.Exit(0))
		Expect(sess).To(gbytes.Say("<Useful information below, please copy-paste from here>"))
	})

	It("does not run as non-root user", func() {
		cmd := exec.Command(dontPanicBin)
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Credential: &syscall.Credential{Uid: 5000, Gid: 5000},
		}
		sess, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess).Should(gexec.Exit())
		Expect(sess.ExitCode()).ToNot(Equal(0))
		Expect(sess.Err).To(gbytes.Say("Keep Calm and Re-run as Root!"))
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

	Context("date.log", func() {

		AfterEach(func() {
			//clean up extract dir
		})

		FIt("write the date into a date.log file in the os-report dir", func() {
			Expect(getTarballPath()).To(BeEmpty())
			cmd := exec.Command(dontPanicBin)
			sess, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(sess).Should(gexec.Exit(0))

			tarball := getTarballPath()
			Expect(tarball).ToNot(BeEmpty())

			tarballContainsValidDateLog(tarball)
		})
	})
})

func tarballContainsValidDateLog(tarball string) {
	extractDir := extractTarball(tarball)
	fmt.Printf("extractDir= %+v\n", extractDir)
	osReportDir := getOsReportDir(extractDir)
	fmt.Printf("osReportDir= %+v\n", osReportDir)
	Expect(osReportDir).ToNot(BeEmpty())
	dateLogFile := filepath.Join(osReportDir, "date.log")
	Expect(dateLogFile).To(BeAnExistingFile())
	contents, err := ioutil.ReadFile(dateLogFile)
	Expect(err).NotTo(HaveOccurred())
	Expect(string(contents)).To(MatchRegexp(`foo bar`))
}

func extractTarball(tarball string) string {
	tmpDir, err := ioutil.TempDir("", "")
	Expect(err).NotTo(HaveOccurred())
	cmd := exec.Command("tar", "xf", tarball, "-C", tmpDir)
	err = cmd.Run()
	Expect(err).NotTo(HaveOccurred())
	return tmpDir
}

func getOsReportDir(extractDir string) string {
	entries, err := ioutil.ReadDir(extractDir)
	Expect(err).NotTo(HaveOccurred())
	re := regexp.MustCompile(`^os-report-.*-\d\d\d\d-\d\d-\d\d-\d\d-\d\d-\d\d/`)
	for _, entry := range entries {
		if re.MatchString(entry.Name()) {
			return filepath.Join(extractDir, entry.Name())
		}
	}
	return ""
}

func getTarballPath() string {
	targetDir := "/var/vcap/data/tmp"
	dirEntries, err := ioutil.ReadDir(targetDir)
	if err != nil {
		return ""
	}

	re := regexp.MustCompile(`os-report-.*\d\d-\d\d-\d\d.tar.gz`)
	for _, f := range dirEntries {
		if f.IsDir() {
			continue
		}
		if re.MatchString(f.Name()) {
			return filepath.Join(targetDir, f.Name())
		}
	}
	return ""
}
