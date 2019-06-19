package command_test

import (
	"context"
	"io"
	"io/ioutil"
	"path/filepath"
	"time"

	"code.cloudfoundry.org/dontpanic/collectors/command"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Command Runner", func() {
	var (
		cmd      string
		ctx      context.Context
		cancel   context.CancelFunc
		dstPath  string
		stdout   io.Writer
		filename string

		err error
	)

	BeforeEach(func() {
		ctx = context.TODO()
		dstPath, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())
		stdout = gbytes.NewBuffer()
		filename = "hello"
	})

	JustBeforeEach(func() {
		err = command.New(cmd, filename).Run(ctx, dstPath, stdout)
	})

	When("cmd is a simple executable", func() {
		BeforeEach(func() {
			cmd = "echo hello world"
		})

		It("returns the byte output", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(stdout).To(gbytes.Say("hello world\n"))

			fileContents, err := ioutil.ReadFile(filepath.Join(dstPath, filename))
			Expect(err).NotTo(HaveOccurred())
			Expect(fileContents).To(Equal([]byte("hello world\n")))
		})
	})

	When("cmd is a pipeline", func() {
		BeforeEach(func() {
			cmd = "seq 2 10 | wc -l"
		})

		It("executes the pipeline", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(stdout).To(gbytes.Say("9"))

			fileContents, err := ioutil.ReadFile(filepath.Join(dstPath, filename))
			Expect(err).NotTo(HaveOccurred())
			Expect(fileContents).To(Equal([]byte("9\n")))
		})
	})

	When("cmd exceeds time limit", func() {
		BeforeEach(func() {
			ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond)
			cmd = "sleep 1"
		})

		AfterEach(func() {
			cancel()
		})

		It("times out", func() {
			Expect(err).To(Equal(context.DeadlineExceeded))
		})
	})

	When("command fails and has stdout and stderr", func() {
		BeforeEach(func() {
			cmd = "echo foo; echo bar >&2; exit 1"
		})

		It("returns error containing bar", func() {
			Expect(err).To(MatchError(ContainSubstring("bar")))
			Expect(filepath.Join(dstPath, filename)).NotTo(BeAnExistingFile())
			Expect(stdout).NotTo(gbytes.Say("foo"))
		})
	})

})
