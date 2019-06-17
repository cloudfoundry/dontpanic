package uptime

import (
	"context"
	"os/exec"
)

func Run(ctx context.Context) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "uptime")
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return out, ctx.Err()
	}
	return out, err
}
