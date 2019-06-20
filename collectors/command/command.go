package command

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"code.cloudfoundry.org/dontpanic/commandrunner"
)

type Collector struct {
	cmd      string
	filename string
	runner   commandrunner.CommandRunner
}

func NewCollector(cmd, filename string) Collector {
	return Collector{
		cmd:      cmd,
		filename: filename,
		runner:   commandrunner.CommandRunner{},
	}
}

func (c Collector) Run(ctx context.Context, destPath string, stdout io.Writer) error {
	out, err := c.runner.Run(ctx, "sh", "-c", c.cmd)
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
