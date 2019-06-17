package command

import (
	"context"
	"os/exec"

	"code.cloudfoundry.org/dontpanic/osreporter"
)

func New(cmd string) osreporter.StreamPlugin {
	return func(ctx context.Context) ([]byte, error) {
		c := exec.CommandContext(ctx, "sh", "-c", cmd)
		out, err := c.CombinedOutput()
		if ctx.Err() == context.DeadlineExceeded {
			return out, ctx.Err()
		}
		return out, err
	}
}
