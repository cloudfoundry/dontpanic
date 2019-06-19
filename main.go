package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"code.cloudfoundry.org/dontpanic/collectors/command"
	"code.cloudfoundry.org/dontpanic/osreporter"
)

func main() {
	checkIsRoot()
	checkIsNotBpm()
	reportDir := createReportDir("/var/vcap/data/tmp")

	osReporter := osreporter.New(reportDir, os.Stdout)

	osReporter.RegisterNoisyCollector("Date", command.New("date", "date.log"))
	osReporter.RegisterNoisyCollector("Uptime", command.New("uptime", "uptime.log"))
	osReporter.RegisterNoisyCollector("Garden Version", command.New("/var/vcap/packages/guardian/bin/gdn -v", "gdn-version.log"))
	osReporter.RegisterNoisyCollector("Hostname", command.New("hostname", "hostname.log"))

	if err := osReporter.Run(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func checkIsRoot() {
	if currentUID := os.Getuid(); currentUID != 0 {
		log.Fatalf("Keep Calm and Re-run as Root!")
	}
}

func checkIsNotBpm() {
	contents, err := ioutil.ReadFile("/proc/1/cmdline")
	if err != nil {
		log.Fatalf("Cannot determine if running in bpm: cannot read cmdline")
	}

	if bytes.Contains(contents, []byte("garden_start")) {
		log.Fatalf("Keep Calm and Re-run outside the BPM container!")
	}
}

func createReportDir(baseDir string) string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("could not determine hostname")
		hostname = "UNKNOWN-HOSTNAME"
	}
	timestamp := time.Now().Format("2006-01-02-15-04-05")
	reportDir := fmt.Sprintf("os-report-%s-%s", hostname, timestamp)
	path := filepath.Join(baseDir, reportDir)
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Fatalf("cannot create report directory %q: %s", path, err.Error())
	}
	return path
}
