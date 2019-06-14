package osreporter

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	osReportDirPattern = "os-report-%s-%s"
)

type Runner struct {
	baseReportDir string
	hostname      string
	timestamp     time.Time
	ReportPath    string
	TarballPath   string
}

func New(baseReportDir, hostname string, now time.Time) Runner {
	r := Runner{
		baseReportDir: baseReportDir,
		hostname:      hostname,
		timestamp:     now,
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
	if currentUID := os.Getuid(); currentUID != 0 {
		fmt.Fprintf(os.Stderr, "Keep Calm and Re-run as Root!")
		return fmt.Errorf("must be run as root")
	}

	if err := os.MkdirAll(r.ReportPath, 0755); err != nil {
		return err
	}

	fmt.Println("<Useful information below, please copy-paste from here>")

	if err := r.writeDate(); err != nil {
		return err
	}

	if err := r.createTarball(); err != nil {
		return err
	}

	return nil
}

func (r Runner) writeDate() error {
	f := filepath.Join(r.ReportPath, "date.log")
	d := r.timestamp.Format(time.UnixDate)
	return ioutil.WriteFile(f, []byte(d), 0644)
}

func (r Runner) createTarball() error {
	cmd := exec.Command("tar", "cf", r.TarballPath, "-C", r.baseReportDir, filepath.Base(r.ReportPath))
	return cmd.Run()
}
