package file_test

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"

	"code.cloudfoundry.org/dontpanic/collectors/file"
)

var _ = Describe("file.Collector", func() {
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

		fileToCollect, err := ioutil.TempFile("", "")
		Expect(err).NotTo(HaveOccurred())
		err = ioutil.WriteFile(fileToCollect.Name(), []byte("file-to-collect"), 0755)
		Expect(err).NotTo(HaveOccurred())
		sourcePath = fileToCollect.Name()
	})

	AfterEach(func() {
		Expect(os.RemoveAll(sourcePath)).To(Succeed())
		Expect(os.RemoveAll(destinationPath)).To(Succeed())
	})

	JustBeforeEach(func() {
		err = file.NewCollector(sourcePath, "destination_file").Run(ctx, destinationPath, stdout)
	})

	It("copies the specified file to the destination directory", func() {
		destinationFilePath := filepath.Join(destinationPath, "destination_file")

		Expect(err).NotTo(HaveOccurred())
		Expect(destinationFilePath).To(BeAnExistingFile())

		destinationFileContents, err := ioutil.ReadFile(destinationFilePath)
		Expect(err).NotTo(HaveOccurred())
		Expect(destinationFileContents).To(Equal([]byte("file-to-collect")))
	})

	When("the copy fails", func() {
		BeforeEach(func() {
			sourcePath = "/i/do/not/exist"
		})

		It("returns an error", func() {
			Expect(err).To(MatchError(ContainSubstring("No such file or directory")))
		})
	})

	When("the destination path contains folders that don't exist", func() {
		BeforeEach(func() {
			destinationPath = filepath.Join(destinationPath, "foo", "bar")
		})

		It("creates the missing directories", func() {
			Expect(err).NotTo(HaveOccurred())

			Expect(destinationPath).To(BeADirectory())
			Expect(filepath.Join(destinationPath, "destination_file")).To(BeAnExistingFile())
		})
	})
})
