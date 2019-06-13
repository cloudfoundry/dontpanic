package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	tempDirPath        = "/var/vcap/data/tmp"
	osReportDirPattern = "os-report-%s-%s"
	osReportTarPattern = tempDirPath + "/os-report-%s-%s.tar.gz"
)

func main() {

	if currentUID := os.Getuid(); currentUID != 0 {
		fmt.Fprintf(os.Stderr, "Keep Calm and Re-run as Root!")
		os.Exit(1)
	}

	fmt.Println("<Useful information below, please copy-paste from here>")

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "UNKNOWN-HOSTNAME"
	}
	timestamp := time.Now().Format("2006-01-02-15-04-05")

	reportDir := fmt.Sprintf(osReportDirPattern, hostname, timestamp)
	reportPath := filepath.Join(tempDirPath, reportDir)

	tarballPath := fmt.Sprintf(osReportTarPattern, hostname, timestamp)

	err = os.MkdirAll(reportPath, 0755)
	if err != nil {
		log.Fatal(err)
	}

	createTarball(reportDir, tarballPath)
}

func createTarball(reportDirName, tarballPath string) error {
	tarIt := exec.Command("tar", "cf", tarballPath, "-C", tempDirPath, reportDirName)
	fmt.Printf("tarIt = %+v\n", tarIt)
	return tarIt.Run()
}
