package grootfs_test

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"code.cloudfoundry.org/dontpanic/collectors/grootfs"
	"code.cloudfoundry.org/dontpanic/collectors/grootfs/grootfsfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Grootfs", func() {

	var (
		collector     grootfs.UsageCollector
		tmpDir        string
		fakeRunner    *grootfsfakes.FakeCommandRunner
		ctx           context.Context
		stdout        io.Writer
		runError      error
		usageFilePath string
	)

	BeforeEach(func() {
		var err error
		tmpDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		depotDir := filepath.Join(tmpDir, "depot")
		Expect(os.MkdirAll(filepath.Join(depotDir, "container1"), 0755)).To(Succeed())
		Expect(os.MkdirAll(filepath.Join(depotDir, "container2"), 0755)).To(Succeed())

		fakeRunner = new(grootfsfakes.FakeCommandRunner)
		collector = grootfs.NewUsageCollector(tmpDir, depotDir, grootfs.Unprivileged, fakeRunner)
		ctx = context.TODO()
		stdout = gbytes.NewBuffer()
		usageFilePath = filepath.Join(tmpDir, "grootfs", "unprivileged-usage.txt")

		fakeRunner.RunStub = func(ctx context.Context, cmd string, args ...string) ([]byte, error) {
			switch cmd {
			case "du":
				return []byte("123456\tfoo/bar/sha\n"), nil
			case "/var/vcap/packages/grootfs/bin/grootfs":
				return []byte(`{"disk_usage": {"exclusive_bytes_used": 3040}}`), nil
			default:
				Fail("unexpected " + args[0])
			}
			return nil, nil
		}
	})

	AfterEach(func() {
		os.RemoveAll(tmpDir)
	})

	JustBeforeEach(func() {
		runError = collector.Run(ctx, tmpDir, stdout)
		Expect(runError).NotTo(HaveOccurred())
	})

	Context("volumes", func() {
		BeforeEach(func() {

			linkDir := filepath.Join(tmpDir, "unprivileged", "l")
			metaDir := filepath.Join(tmpDir, "unprivileged", "meta")
			imageDir1 := filepath.Join(tmpDir, "unprivileged", "images", "container1")
			imageDir2 := filepath.Join(tmpDir, "unprivileged", "images", "container2")
			volumeDir1 := filepath.Join(tmpDir, "unprivileged", "volumes", "abc123")
			volumeDir2 := filepath.Join(tmpDir, "unprivileged", "volumes", "abc456")

			Expect(os.MkdirAll(linkDir, 0755)).To(Succeed())
			Expect(os.MkdirAll(metaDir, 0755)).To(Succeed())
			Expect(os.MkdirAll(imageDir1, 0755)).To(Succeed())
			Expect(os.MkdirAll(imageDir2, 0755)).To(Succeed())
			Expect(os.MkdirAll(volumeDir1, 0755)).To(Succeed())
			Expect(os.MkdirAll(volumeDir2, 0755)).To(Succeed())
			Expect(os.Symlink(volumeDir1, filepath.Join(linkDir, "link1"))).To(Succeed())
			Expect(os.Symlink(volumeDir2, filepath.Join(linkDir, "link2"))).To(Succeed())
			Expect(ioutil.WriteFile(filepath.Join(metaDir, "volume-abc123"), []byte(`{"Size": 1024}`), 0644)).To(Succeed())
			Expect(ioutil.WriteFile(filepath.Join(metaDir, "volume-abc456"), []byte(`{"Size": 2048}`), 0644)).To(Succeed())
		})

		It("adds volume dir size to usage", func() {
			Expect(contents(usageFilePath)).To(ContainSubstring("volumes-total:                 123456 bytes"))
			Expect(contents(usageFilePath)).To(ContainSubstring("volumes-used:                    3072 bytes"))
			Expect(contents(usageFilePath)).To(ContainSubstring("volumes-cleanable:             120384 bytes"))
		})

		It("adds image used size to usage", func() {
			Expect(contents(usageFilePath)).To(ContainSubstring("images-exclusive:                6080 bytes"))
		})

	})
})

func contents(path string) string {
	b, err := ioutil.ReadFile(path)
	Expect(err).NotTo(HaveOccurred())
	return string(b)
}
