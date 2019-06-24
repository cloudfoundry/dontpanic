package integration_test

import (
	"io/ioutil"
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
	baseDir      = "/var/vcap/data/tmp"
	dateRegexp   = `\w{3} \w{3} \d{1,2}.*\d{4}.*`
	uptimeRegexp = `day.*user.*load average`
)

var _ = Describe("Integration", func() {
	var (
		session         *gexec.Session
		cmd             *exec.Cmd
		gardenConfigDir = "/var/vcap/jobs/garden/config"
		gardenLogDir    = "/var/vcap/sys/log/garden"
		varLogDir       = "/var/log"
		gardenDepotDir  = "/var/vcap/data/garden/depot"
		monitDir        = "/var/vcap/monit"
	)

	BeforeEach(func() {
		Expect(os.MkdirAll(gardenConfigDir, 0755)).To(Succeed())
		Expect(ioutil.WriteFile(filepath.Join(gardenConfigDir, "config.ini"), []byte("hi"), 0644)).To(Succeed())
		Expect(ioutil.WriteFile(filepath.Join(gardenConfigDir, "grootfs_config.yml"), []byte("groot"), 0644)).To(Succeed())
		Expect(ioutil.WriteFile(filepath.Join(gardenConfigDir, "privileged_grootfs_config.yml"), []byte("groot"), 0644)).To(Succeed())
		Expect(ioutil.WriteFile(filepath.Join(gardenConfigDir, "bpm.yml"), []byte("bpm"), 0644)).To(Succeed())
		Expect(ioutil.WriteFile(filepath.Join(gardenConfigDir, "containerd.toml"), []byte("nerd"), 0644)).To(Succeed())

		Expect(os.MkdirAll(gardenLogDir, 0755)).To(Succeed())
		Expect(ioutil.WriteFile(filepath.Join(gardenLogDir, "garden.log"), []byte("cur"), 0644)).To(Succeed())
		Expect(ioutil.WriteFile(filepath.Join(gardenLogDir, "garden.log.1"), []byte("prev"), 0644)).To(Succeed())
		Expect(ioutil.WriteFile(filepath.Join(gardenLogDir, "garden.log.2.gz"), []byte("Z"), 0644)).To(Succeed())

		Expect(os.MkdirAll(varLogDir, 0755)).To(Succeed())
		Expect(ioutil.WriteFile(filepath.Join(varLogDir, "kern.log"), []byte("cur"), 0644)).To(Succeed())
		Expect(ioutil.WriteFile(filepath.Join(varLogDir, "kern.log.1"), []byte("prev"), 0644)).To(Succeed())
		Expect(ioutil.WriteFile(filepath.Join(varLogDir, "kern.log.2.gz"), []byte("Z"), 0644)).To(Succeed())

		Expect(ioutil.WriteFile(filepath.Join(varLogDir, "syslog"), []byte("cur"), 0644)).To(Succeed())
		Expect(ioutil.WriteFile(filepath.Join(varLogDir, "syslog.1"), []byte("prev"), 0644)).To(Succeed())
		Expect(ioutil.WriteFile(filepath.Join(varLogDir, "syslog.2.gz"), []byte("Z"), 0644)).To(Succeed())

		Expect(os.MkdirAll(monitDir, 0755)).To(Succeed())
		Expect(ioutil.WriteFile(filepath.Join(monitDir, "monit.log"), []byte("monit"), 0644)).To(Succeed())

		Expect(os.MkdirAll(filepath.Join(gardenDepotDir, "container1"), 0755)).To(Succeed())

		cmd = exec.Command(dontPanicBin)
	})

	AfterEach(func() {
		Expect(os.RemoveAll(gardenConfigDir)).To(Succeed())
		Expect(os.RemoveAll(gardenLogDir)).To(Succeed())
	})

	JustBeforeEach(func() {
		var err error
		session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session, time.Second*50).Should(gexec.Exit())
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

		By("collecting the free memory")
		tarballShouldContainFile("free.log")
		Expect(string(tarballFileContents("free.log"))).
			To(ContainSubstring("Mem:"))

		By("collecting the kernel details")
		tarballShouldContainFile("uname.log")
		Expect(string(tarballFileContents("uname.log"))).
			To(ContainSubstring("Linux"))

		By("collecting monit summary")
		Expect(session).To(gbytes.Say("## Monit Summary"))

		By("collecting the number of containers")
		tarballShouldContainFile("num-containers.log")
		Expect(string(tarballFileContents("num-containers.log"))).
			To(ContainSubstring("1"))

		By("collecting the number of open files")
		tarballShouldContainFile("num-open-files.log")
		Expect(tarballFileContents("num-open-files.log")).ToNot(BeEmpty())

		By("collecting the max number of open files")
		tarballShouldContainFile("file-max.log")
		Expect(tarballFileContents("file-max.log")).ToNot(BeEmpty())

		By("collecting the disk usage")
		tarballShouldContainFile("df.log")
		Expect(tarballFileContents("df.log")).To(ContainSubstring("Filesystem"))

		By("collecting the open files")
		tarballShouldContainFile("lsof.log")
		Expect(tarballFileContents("lsof.log")).To(ContainSubstring("COMMAND"))

		By("collecting the process information")
		tarballShouldContainFile("ps-info.log")
		Expect(tarballFileContents("ps-info.log")).To(ContainSubstring("PID"))

		By("collecting the process forest information")
		tarballShouldContainFile("ps-forest.log")
		Expect(tarballFileContents("ps-forest.log")).To(ContainSubstring("USER"))

		By("collecting the dmesg")
		tarballShouldContainFile("dmesg.log")
		Expect(tarballFileContents("dmesg.log")).
			To(MatchRegexp(dateRegexp))

		By("collecting the network interfaces")
		tarballShouldContainFile("ifconfig.log")
		Expect(tarballFileContents("ifconfig.log")).To(ContainSubstring("Link"))

		By("collecting the firewall configuration")
		tarballShouldContainFile("iptables-L.log")
		Expect(tarballFileContents("iptables-L.log")).To(ContainSubstring("Chain"))

		By("collecting the NAT info")
		tarballShouldContainFile("iptables-tnat.log")
		Expect(tarballFileContents("iptables-tnat.log")).To(ContainSubstring("Chain"))

		By("collecting the mount table")
		Expect(session).To(gbytes.Say("## Mount Table"))

		By("collecting Garden depot contents")
		tarballShouldContainFile("depot-contents.log")
		Expect(tarballFileContents("depot-contents.log")).To(ContainSubstring("depot"))

		By("collecting XFS fragmentation info")
		Expect(session).To(gbytes.Say("## XFS Fragmentation"))

		By("collecting XFS info")
		Expect(session).To(gbytes.Say("## XFS Info"))

		By("collecting Slabinfo")
		tarballShouldContainFile("slabinfo.log")
		Expect(tarballFileContents("slabinfo.log")).To(ContainSubstring("active_objs"))

		By("collecting Meminfo")
		tarballShouldContainFile("meminfo.log")
		Expect(tarballFileContents("meminfo.log")).To(ContainSubstring("MemTotal"))

		By("collecting iostat")
		tarballShouldContainFile("iostat.log")
		Expect(tarballFileContents("iostat.log")).To(ContainSubstring("Linux"))

		By("collecting vm statistics")
		tarballShouldContainFile("vmstat-s.log")
		Expect(tarballFileContents("vmstat-s.log")).To(ContainSubstring("memory"))

		By("collecting disk statistics")
		tarballShouldContainFile("vmstat-d.log")
		Expect(tarballFileContents("vmstat-d.log")).To(ContainSubstring("disk"))

		By("collecting active and inactive memory statistics")
		tarballShouldContainFile("vmstat-a.log")
		Expect(tarballFileContents("vmstat-a.log")).To(ContainSubstring("memory"))

		By("collecting mass process data")
		currentPid := os.Getpid()
		tarballShouldContainFile(filepath.Join("process-data", strconv.Itoa(currentPid), "fd"))
		tarballShouldContainFile(filepath.Join("process-data", strconv.Itoa(currentPid), "ns"))
		tarballShouldContainFile(filepath.Join("process-data", strconv.Itoa(currentPid), "cgroup"))
		tarballShouldContainFile(filepath.Join("process-data", strconv.Itoa(currentPid), "stack"))
		tarballShouldContainFile(filepath.Join("process-data", strconv.Itoa(currentPid), "status"))

		By("collecting the kernel logs")
		tarballShouldContainFile("kernel-logs/kern.log")
		tarballShouldContainFile("kernel-logs/kern.log.1")
		tarballShouldContainFile("kernel-logs/kern.log.2.gz")

		By("collecting monit log")
		tarballShouldContainFile("monit.log")
		Expect(tarballFileContents("monit.log")).To(ContainSubstring("monit"))

		By("collecting the syslogs")
		tarballShouldContainFile("syslogs/syslog")
		tarballShouldContainFile("syslogs/syslog.1")
		tarballShouldContainFile("syslogs/syslog.2.gz")

		By("collecting all garden config")
		tarballShouldContainFile("config/config.ini")
		Expect(string(tarballFileContents("config/config.ini"))).To(Equal("hi"))
		tarballShouldContainFile("config/grootfs_config.yml")
		tarballShouldContainFile("config/privileged_grootfs_config.yml")
		tarballShouldContainFile("config/bpm.yml")
		tarballShouldContainFile("config/containerd.toml")

		By("collecting the garden logs")
		tarballShouldContainFile("garden/garden.log")
		tarballShouldContainFile("garden/garden.log.1")
		tarballShouldContainFile("garden/garden.log.2.gz")
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
