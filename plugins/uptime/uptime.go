package uptime

import "os/exec"

func Run() ([]byte, error) {
	cmd := exec.Command("uptime")
	return cmd.CombinedOutput()
}
