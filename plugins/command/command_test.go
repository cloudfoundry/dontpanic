package command_test

import (
	"context"
	"time"

	"code.cloudfoundry.org/dontpanic/osreporter"
	"code.cloudfoundry.org/dontpanic/plugins/command"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Command Runner", func() {

	var (
		runner osreporter.StreamPlugin
		cmd    string
		ctx    context.Context
		cancel context.CancelFunc

		out []byte
		err error
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	})

	JustBeforeEach(func() {
		runner = command.New(cmd)
		out, err = runner(ctx)
	})

	AfterEach(func() {
		cancel()
	})

	When("cmd is a simple executable", func() {
		BeforeEach(func() {
			cmd = "echo hello world"
		})

		It("returns the byte output", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(string(out)).To(Equal("hello world\n"))
		})
	})

	When("cmd is a pipeline", func() {
		BeforeEach(func() {
			cmd = "seq 2 10 | wc -l"
		})

		It("executes the pipeline", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(string(out)).To(ContainSubstring("9"))
		})
	})

	When("cmd exceeds time limit", func() {
		BeforeEach(func() {
			ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond)
			cmd = "sleep 1"
		})

		It("times out", func() {
			Expect(err).To(Equal(context.DeadlineExceeded))
		})
	})

})
