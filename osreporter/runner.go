package osreporter

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	osReportDirPattern = "os-report-%s-%s"
)

type RegisteredPlugin struct {
	streamPlugin StreamPlugin
	name         string
	filename     string
}

type Runner struct {
	hostname      string
	timestamp     time.Time
	baseReportDir string
	out           io.Writer
	TarballPath   string
	ReportPath    string
	plugins       []RegisteredPlugin
}

func New(baseReportDir, hostname string, now time.Time, out io.Writer) Runner {
	r := Runner{
		baseReportDir: baseReportDir,
		hostname:      hostname,
		timestamp:     now,
		out:           out,
	}

	r.SetPaths()
	return r
}

func (r *Runner) SetPaths() {
	timestamp := r.timestamp.Format("2006-01-02-15-04-05")
	reportDir := fmt.Sprintf(osReportDirPattern, r.hostname, timestamp)
	r.ReportPath = filepath.Join(r.baseReportDir, reportDir)
	r.TarballPath = r.ReportPath + ".tar.gz"
}

func (r Runner) Run() error {
	if err := r.sanityCheck(); err != nil {
		return err
	}

	fmt.Fprintln(r.out, "<Useful information below, please copy-paste from here>")

	for _, plugin := range r.plugins {
		fmt.Fprintln(r.out, "## "+plugin.name)
		if plugin.streamPlugin != nil {
			out, err := plugin.streamPlugin()
			if err != nil {
				fmt.Fprintln(r.out, "Failure:", err.Error())
				continue
			}
			outPath := filepath.Join(r.ReportPath, plugin.filename)
			err = ioutil.WriteFile(outPath, out, 0644)
			if err != nil {
				fmt.Fprintln(r.out, "Failed to write file:", err.Error())
			}
		}
	}

	if err := r.createTarball(); err != nil {
		return err
	}

	return nil
}

func (r Runner) sanityCheck() error {
	if currentUID := os.Getuid(); currentUID != 0 {
		fmt.Fprintf(os.Stderr, "Keep Calm and Re-run as Root!")
		return fmt.Errorf("must be run as root")
	}

	if err := os.MkdirAll(r.ReportPath, 0755); err != nil {
		fmt.Fprintln(os.Stderr, "Cannot create report directory - exiting", err.Error())
		return err
	}
	return nil
}

func (r Runner) createTarball() error {
	cmd := exec.Command("tar", "cf", r.TarballPath, "-C", r.baseReportDir, filepath.Base(r.ReportPath))
	return cmd.Run()
}
