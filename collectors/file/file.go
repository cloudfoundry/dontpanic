package file

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"code.cloudfoundry.org/dontpanic/commandrunner"
)

type Collector struct {
	sourcePath      string
	destinationPath string
	runner          commandrunner.CommandRunner
}

func NewCollector(sourcePath, destinationPath string) Collector {
	return Collector{
		sourcePath:      sourcePath,
		destinationPath: destinationPath,
		runner:          commandrunner.CommandRunner{},
	}
}

func (c Collector) Run(ctx context.Context, reportDir string, stdout io.Writer) error {
	fullDestinationPath := filepath.Join(reportDir, c.destinationPath)
	if err := os.MkdirAll(filepath.Dir(fullDestinationPath), 0755); err != nil {
		return err
	}

	_, err := c.runner.Run(ctx, "cp", c.sourcePath, fullDestinationPath)
	return err
}
