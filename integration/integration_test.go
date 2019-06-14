package integration_test

import (
	"os/exec"
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

})
