package integration_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

const (
	baseDir    = "/var/vcap/data/tmp"
	dateRegexp = `\w{3} \w{3}  ?\d{1,2}.*\d{4}.*`
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
		Eventually(session, time.Second*50).Should(gexec.Exit())
	})

	It("produces a report correctly", func() {
		reportDir := getReportDir(session.Out.Contents())
		tarPath := reportDir + ".tar.gz"

		By("succeeding")
		Expect(session.ExitCode()).To(Equal(0))

		By("showing an initial message")
		Expect(session).To(gbytes.Say("<Useful information below, please copy-paste from here>"))

		By("writing a logfile containing all steps")
		tarballShouldContainFile(tarPath, "dontpanic.log")
		Expect(string(tarballFileContents(tarPath, "dontpanic.log"))).
			To(ContainSubstring("## Date"))

		By("collecting the date")
		tarballShouldContainFile(tarPath, "date.log")
		Expect(string(tarballFileContents(tarPath, "date.log"))).
			To(MatchRegexp(dateRegexp))

		By("collecting the uptime")
		tarballShouldContainFile(tarPath, "uptime.log")
		Expect(string(tarballFileContents(tarPath, "uptime.log"))).
			To(ContainSubstring("load average"))

		By("collecting the garden version")
		Expect(session).To(gbytes.Say("## Garden Version"))

		By("collecting the hostname")
		tarballShouldContainFile(tarPath, "hostname.log")

		By("collecting the free memory")
		tarballShouldContainFile(tarPath, "free.log")
		Expect(string(tarballFileContents(tarPath, "free.log"))).
			To(ContainSubstring("Mem:"))

		By("collecting the kernel details")
		tarballShouldContainFile(tarPath, "uname.log")
		Expect(string(tarballFileContents(tarPath, "uname.log"))).
			To(ContainSubstring("Linux"))

		By("collecting monit summary")
		Expect(session).To(gbytes.Say("## Monit Summary"))

		By("collecting the number of containers")
		tarballShouldContainFile(tarPath, "num-containers.log")
		Expect(string(tarballFileContents(tarPath, "num-containers.log"))).
			To(ContainSubstring("1"))

		By("collecting the number of open files")
		tarballShouldContainFile(tarPath, "num-open-files.log")
		Expect(tarballFileContents(tarPath, "num-open-files.log")).ToNot(BeEmpty())

		By("collecting the max number of open files")
		tarballShouldContainFile(tarPath, "file-max.log")
		Expect(tarballFileContents(tarPath, "file-max.log")).ToNot(BeEmpty())

		By("collecting the disk usage")
		tarballShouldContainFile(tarPath, "df.log")
		Expect(tarballFileContents(tarPath, "df.log")).To(ContainSubstring("Filesystem"))

		By("collecting the open files")
		tarballShouldContainFile(tarPath, "lsof.log")
		Expect(tarballFileContents(tarPath, "lsof.log")).To(ContainSubstring("COMMAND"))

		By("collecting the process information")
		tarballShouldContainFile(tarPath, "ps-info.log")
		Expect(tarballFileContents(tarPath, "ps-info.log")).To(ContainSubstring("PID"))

		By("collecting the process forest information")
		tarballShouldContainFile(tarPath, "ps-forest.log")
		Expect(tarballFileContents(tarPath, "ps-forest.log")).To(ContainSubstring("USER"))

		By("collecting the dmesg")
		tarballShouldContainFile(tarPath, "dmesg.log")
		Expect(tarballFileContents(tarPath, "dmesg.log")).
			To(MatchRegexp(dateRegexp))

		By("collecting the network interfaces")
		tarballShouldContainFile(tarPath, "ifconfig.log")
		Expect(tarballFileContents(tarPath, "ifconfig.log")).To(ContainSubstring("Link"))

		By("collecting the firewall configuration")
		tarballShouldContainFile(tarPath, "iptables-L.log")
		Expect(tarballFileContents(tarPath, "iptables-L.log")).To(ContainSubstring("Chain"))

		By("collecting the NAT info")
		tarballShouldContainFile(tarPath, "iptables-tnat.log")
		Expect(tarballFileContents(tarPath, "iptables-tnat.log")).To(ContainSubstring("Chain"))

		By("collecting the mount table")
		Expect(session).To(gbytes.Say("## Mount Table"))

		By("collecting Garden depot contents")
		tarballShouldContainFile(tarPath, "depot-contents.log")
		Expect(tarballFileContents(tarPath, "depot-contents.log")).To(ContainSubstring("depot"))

		By("collecting XFS fragmentation info")
		Expect(session).To(gbytes.Say("## XFS Fragmentation"))

		By("collecting XFS info")
		Expect(session).To(gbytes.Say("## XFS Info"))

		By("collecting Slabinfo")
		tarballShouldContainFile(tarPath, "slabinfo.log")
		Expect(tarballFileContents(tarPath, "slabinfo.log")).To(ContainSubstring("active_objs"))

		By("collecting Meminfo")
		tarballShouldContainFile(tarPath, "meminfo.log")
		Expect(tarballFileContents(tarPath, "meminfo.log")).To(ContainSubstring("MemTotal"))

		By("collecting iostat")
		tarballShouldContainFile(tarPath, "iostat.log")
		Expect(tarballFileContents(tarPath, "iostat.log")).To(ContainSubstring("Linux"))

		By("collecting vm statistics")
		tarballShouldContainFile(tarPath, "vmstat-s.log")
		Expect(tarballFileContents(tarPath, "vmstat-s.log")).To(ContainSubstring("memory"))

		By("collecting disk statistics")
		tarballShouldContainFile(tarPath, "vmstat-d.log")
		Expect(tarballFileContents(tarPath, "vmstat-d.log")).To(ContainSubstring("disk"))

		By("collecting active and inactive memory statistics")
		tarballShouldContainFile(tarPath, "vmstat-a.log")
		Expect(tarballFileContents(tarPath, "vmstat-a.log")).To(ContainSubstring("memory"))

		By("collecting mass process data")
		currentPid := os.Getpid()
		tarballShouldContainFile(tarPath, filepath.Join("process-data", strconv.Itoa(currentPid), "fd"))
		tarballShouldContainFile(tarPath, filepath.Join("process-data", strconv.Itoa(currentPid), "ns"))
		tarballShouldContainFile(tarPath, filepath.Join("process-data", strconv.Itoa(currentPid), "cgroup"))
		tarballShouldContainFile(tarPath, filepath.Join("process-data", strconv.Itoa(currentPid), "stack"))
		tarballShouldContainFile(tarPath, filepath.Join("process-data", strconv.Itoa(currentPid), "status"))

		By("collecting the kernel logs")
		tarballShouldContainFile(tarPath, "kernel-logs/kern.log")
		tarballShouldContainFile(tarPath, "kernel-logs/kern.log.1")
		tarballShouldContainFile(tarPath, "kernel-logs/kern.log.2.gz")

		By("collecting monit log")
		tarballShouldContainFile(tarPath, "monit.log")
		Expect(tarballFileContents(tarPath, "monit.log")).To(ContainSubstring("monit"))

		By("collecting the syslogs")
		tarballShouldContainFile(tarPath, "syslogs/syslog")
		tarballShouldContainFile(tarPath, "syslogs/syslog.1")
		tarballShouldContainFile(tarPath, "syslogs/syslog.2.gz")

		By("collecting all garden config")
		tarballShouldContainFile(tarPath, "config/config.ini")
		Expect(string(tarballFileContents(tarPath, "config/config.ini"))).To(Equal("hi"))
		tarballShouldContainFile(tarPath, "config/grootfs_config.yml")
		tarballShouldContainFile(tarPath, "config/privileged_grootfs_config.yml")
		tarballShouldContainFile(tarPath, "config/bpm.yml")
		tarballShouldContainFile(tarPath, "config/containerd.toml")

		By("collecting the garden logs")
		tarballShouldContainFile(tarPath, "garden/garden.log")
		tarballShouldContainFile(tarPath, "garden/garden.log.1")
		tarballShouldContainFile(tarPath, "garden/garden.log.2.gz")

		By("deleting the report dir at the end")
		Expect(reportDir).ToNot(BeADirectory())
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

func tarballShouldContainFile(tarballPath, filePath string) {
	ExpectWithOffset(1, tarballPath).ToNot(BeEmpty(), "tarball not found in "+baseDir)

	extractedOsReportPath := strings.TrimRight(filepath.Base(tarballPath), ".tar.gz")
	logFilePath := filepath.Join(extractedOsReportPath, filePath)
	ExpectWithOffset(1, listTarball(tarballPath)).To(ContainSubstring(logFilePath))
}

func tarballFileContents(tarballPath, filePath string) []byte {
	ExpectWithOffset(1, tarballPath).ToNot(BeEmpty(), "tarball not found in "+baseDir)

	extractedOsReportPath := strings.TrimRight(filepath.Base(tarballPath), ".tar.gz")
	osDir := filepath.Base(extractedOsReportPath)

	cmd := exec.Command("tar", "xf", tarballPath, filepath.Join(osDir, filePath), "-O")
	out, err := cmd.Output()
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return out
}

func listTarball(tarball string) string {
	cmd := exec.Command("tar", "tf", tarball)
	files, err := cmd.Output()
	ExpectWithOffset(2, err).NotTo(HaveOccurred())
	return string(files)
}

func getReportDir(output []byte) string {
	re := regexp.MustCompile(`(\/.*os-report-.*)\.tar\.gz`)
	matches := re.FindStringSubmatch(string(output))
	Expect(matches).To(HaveLen(2))
	return matches[1]
}
