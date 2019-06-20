package archive_test

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"

	"code.cloudfoundry.org/dontpanic/collectors/archive"
)

var _ = Describe("archive.Collector", func() {
	var (
		ctx             context.Context
		destinationPath string
		stdout          io.Writer
		sourcePath      string
		err             error
	)

	BeforeEach(func() {
		ctx = context.TODO()

		var err error
		destinationPath, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		stdout = gbytes.NewBuffer()

		sourcePath, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())
		_, err = os.Create(filepath.Join(sourcePath, "file_to_check"))
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(os.RemoveAll(sourcePath)).To(Succeed())
		Expect(os.RemoveAll(destinationPath)).To(Succeed())
	})

	JustBeforeEach(func() {
		err = archive.NewCollector(sourcePath, "destination_tar.tgz").Run(ctx, destinationPath, stdout)
	})

	It("tars the specified directory to the destination directory", func() {
		Expect(err).NotTo(HaveOccurred())

		destinationTar := filepath.Join(destinationPath, "destination_tar.tgz")
		Expect(destinationTar).To(BeAnExistingFile())

		cmd := exec.Command("tar", "xf", destinationTar, "-C", destinationPath)
		Expect(cmd.Run()).To(Succeed())

		Expect(filepath.Join(destinationPath, "file_to_check")).To(BeAnExistingFile())
	})

	When("tar fails", func() {
		BeforeEach(func() {
			sourcePath = "/i/do/not/exist"
		})

		It("returns an error", func() {
			Expect(err).To(HaveOccurred())
		})
	})

	When("the destination path contains folders that don't exist", func() {
		BeforeEach(func() {
			destinationPath = filepath.Join(destinationPath, "foo", "bar")
		})

		It("creates the missing directories", func() {
			Expect(destinationPath).To(BeADirectory())
			Expect(filepath.Join(destinationPath, "destination_tar.tgz")).To(BeAnExistingFile())
		})
	})
})
