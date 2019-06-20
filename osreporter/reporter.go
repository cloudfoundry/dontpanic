package osreporter

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"time"
)

type Reporter struct {
	stdout     io.Writer
	reportPath string
	collectors []RegisteredCollector
}

//go:generate counterfeiter . Collector
type Collector interface {
	Run(context.Context, string, io.Writer) error
}

func New(reportPath string, stdout io.Writer) Reporter {
	return Reporter{
		reportPath: reportPath,
		stdout:     stdout,
	}
}

func (r *Reporter) RegisterCollector(name string, collector Collector, timeout ...time.Duration) {
	r.registerCollector(name, collector, false, timeout...)
}

func (r *Reporter) RegisterNoisyCollector(name string, collector Collector, timeout ...time.Duration) {
	r.registerCollector(name, collector, true, timeout...)
}

func (r *Reporter) registerCollector(name string, collector Collector, echoOutput bool, timeout ...time.Duration) {
	maxDuration := 10 * time.Second
	if len(timeout) > 0 {
		maxDuration = timeout[0]
	}

	registeredCollector := RegisteredCollector{
		collector:  collector,
		name:       name,
		timeout:    maxDuration,
		echoOutput: echoOutput,
	}
	r.collectors = append(r.collectors, registeredCollector)
}

func (r Reporter) Run() error {
	fmt.Fprintln(r.stdout, "<Useful information below, please copy-paste from here>")

	for _, collector := range r.collectors {
		fmt.Fprintln(r.stdout, "## "+collector.name)

		out := ioutil.Discard
		if collector.echoOutput {
			out = r.stdout
		}

		err := collector.Run(r.reportPath, out)
		if err != nil {
			fmt.Fprintln(r.stdout, "Failure:", err.Error())
		}
	}

	if err := r.createTarball(); err != nil {
		return err
	}

	return nil
}

func (r Reporter) createTarball() error {
	return exec.Command("tar", "cf", r.reportPath+".tar.gz", "-C", filepath.Dir(r.reportPath), filepath.Base(r.reportPath)).Run()
}

type RegisteredCollector struct {
	collector  Collector
	name       string
	echoOutput bool
	timeout    time.Duration
}

func (p RegisteredCollector) Run(dstPath string, out io.Writer) error {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	err := p.collector.Run(ctx, dstPath, out)
	if err == context.DeadlineExceeded {
		return fmt.Errorf("timed out after %s", p.timeout)
	}

	return err
}
