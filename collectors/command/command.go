package command

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

type Collector struct {
	cmd      string
	filename string
}

func New(cmd, filename string) Collector {
	return Collector{cmd: cmd, filename: filename}
}

func (c Collector) Run(ctx context.Context, destPath string, stdout io.Writer) error {
	cmd := exec.CommandContext(ctx, "sh", "-c", c.cmd)
	out, err := cmd.Output()

	if ctx.Err() == context.DeadlineExceeded {
		return ctx.Err()
	}

	if exitErr, ok := err.(*exec.ExitError); ok {
		return errors.New(string(exitErr.Stderr))
	}

	if err != nil {
		return err
	}

	outPath := filepath.Join(destPath, c.filename)

	err = os.MkdirAll(filepath.Dir(outPath), 0755)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(outPath, out, 0644)
	if err != nil {
		return err
	}

	_, err = stdout.Write(out)
	if err != nil {
		return err
	}

	return nil
}
