package integration_test

import (
	"fmt"
	"math/rand"
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

var (
	dontPanicBin string
)

var _ = BeforeSuite(func() {
	bin, err := gexec.Build("code.cloudfoundry.org/dontpanic")
	Expect(err).NotTo(HaveOccurred())
	dontPanicBin = makeBinaryAccessibleToEveryone(bin)
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
	Expect(os.RemoveAll(dontPanicBin)).To(Succeed())
})

func makeBinaryAccessibleToEveryone(binaryPath string) string {
	binaryName := path.Base(binaryPath)
	tempDir := fmt.Sprintf("/tmp/temp-%s-%d", binaryName, rand.Int())
	Expect(os.MkdirAll(tempDir, 0755)).To(Succeed())
	Expect(os.Chmod(tempDir, 0755)).To(Succeed())

	newBinaryPath := filepath.Join(tempDir, binaryName)
	Expect(os.Rename(binaryPath, newBinaryPath)).To(Succeed())
	Expect(os.Chmod(newBinaryPath, 0755)).To(Succeed())

	return newBinaryPath
}
