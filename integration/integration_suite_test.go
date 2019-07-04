package integration_test

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

const (
	gardenConfigDir = "/var/vcap/jobs/garden/config"
	gardenLogDir    = "/var/vcap/sys/log/garden"
	varLogDir       = "/var/log"
	gardenDepotDir  = "/var/vcap/data/garden/depot"
	monitDir        = "/var/vcap/monit"
)

var (
	dontPanicBin string
)

var _ = SynchronizedBeforeSuite(func() []byte {
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
	return []byte{}
}, func(data []byte) {
	bin, err := gexec.Build("code.cloudfoundry.org/dontpanic")
	Expect(err).NotTo(HaveOccurred())
	dontPanicBin = makeBinaryAccessibleToEveryone(bin)
})

var _ = SynchronizedAfterSuite(func() {
	gexec.CleanupBuildArtifacts()
	Expect(os.RemoveAll(dontPanicBin)).To(Succeed())
}, func() {
	Expect(os.RemoveAll(gardenConfigDir)).To(Succeed())
	Expect(os.RemoveAll(gardenLogDir)).To(Succeed())
})

func makeBinaryAccessibleToEveryone(binaryPath string) string {
	binaryName := path.Base(binaryPath)

	tempDir, err := ioutil.TempDir("", binaryName)
	Expect(err).NotTo(HaveOccurred())
	Expect(os.Chmod(tempDir, 0755)).To(Succeed())

	newBinaryPath := filepath.Join(tempDir, binaryName)
	Expect(os.Rename(binaryPath, newBinaryPath)).To(Succeed())
	Expect(os.Chmod(newBinaryPath, 0755)).To(Succeed())

	return newBinaryPath
}
