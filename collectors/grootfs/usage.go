package grootfs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type privilegeType string

const (
	Privileged   privilegeType = "privileged"
	Unprivileged privilegeType = "unprivileged"
)

type UsageCollector struct {
	storeDir      string
	depotDir      string
	privilegeType privilegeType
	runner        CommandRunner
}

//go:generate counterfeiter . CommandRunner

type CommandRunner interface {
	Run(context.Context, string, ...string) ([]byte, error)
}

const dirName = "grootfs"

func NewUsageCollector(storeDir, depotDir string, privilegeType privilegeType, cmdRunner CommandRunner) UsageCollector {
	return UsageCollector{
		storeDir:      storeDir,
		depotDir:      depotDir,
		privilegeType: privilegeType,
		runner:        cmdRunner,
	}
}

func (c UsageCollector) Run(ctx context.Context, reportDir string, stdout io.Writer) error {
	grootfsDir := filepath.Join(reportDir, dirName)
	err := os.MkdirAll(grootfsDir, 0755)
	if err != nil {
		return err
	}

	outputPath := filepath.Join(grootfsDir, string(c.privilegeType)+"-usage.txt")
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	total, err := c.totalVolumesSize()
	if err != nil {
		return err
	}

	size, err := c.usedVolumesSize()
	if err != nil {
		return err
	}

	fmt.Fprintf(outputFile, "%-24s %12d bytes\n", "volumes-total:", total)
	fmt.Fprintf(outputFile, "%-24s %12d bytes\n", "volumes-used:", size)
	fmt.Fprintf(outputFile, "%-24s %12d bytes\n", "volumes-cleanable:", total-size)

	imagesSize, err := c.imagesSize()
	if err != nil {
		return err
	}

	fmt.Fprintf(outputFile, "%-24s %12d bytes\n", "images-exclusive:", imagesSize)

	return nil
}

func (c UsageCollector) imagesSize() (int64, error) {
	containerIds, err := c.getContainerIds()
	if err != nil {
		return 0, err
	}

	var size int64

	for _, id := range containerIds {
		containerImageSize, err := c.getContainerImageSize(id)
		if err != nil {
			return 0, err
		}
		size += containerImageSize
	}

	return size, nil
}

func (c UsageCollector) getContainerIds() ([]string, error) {
	ids := []string{}

	entries, err := ioutil.ReadDir(c.depotDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		ids = append(ids, entry.Name())
	}

	return ids, nil
}

type stats struct {
	DiskUsage struct {
		ExclusiveBytesUsed int64 `json:"exclusive_bytes_used"`
	} `json:"disk_usage"`
}

func (c UsageCollector) getContainerImageSize(id string) (int64, error) {
	if !c.imageExists(id) {
		return 0, nil
	}
	grootConfig := c.getGrootConfigPath()
	output, err := c.runner.Run(context.Background(), "/var/vcap/packages/grootfs/bin/grootfs", "--config", grootConfig, "stats", id)
	if err != nil {
		return 0, err
	}

	var usage stats
	err = json.Unmarshal(output, &usage)
	if err != nil {
		return 0, err
	}

	return usage.DiskUsage.ExclusiveBytesUsed, nil
}

func (c UsageCollector) imageExists(id string) bool {
	imageDir := filepath.Join(c.storeDir, string(c.privilegeType), "images", id)
	stat, err := os.Stat(imageDir)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

func (c UsageCollector) getGrootConfigPath() string {
	if c.privilegeType == Privileged {
		return "/var/vcap/jobs/garden/config/privileged_grootfs_config.yml"
	}
	return "/var/vcap/jobs/garden/config/grootfs_config.yml"
}

type metaData struct {
	Size int64 `json:"Size"`
}

func (c UsageCollector) volumeSize(id string) (int64, error) {
	metaFile := filepath.Join(c.storeDir, string(c.privilegeType), "meta", "volume-"+id)
	contents, err := ioutil.ReadFile(metaFile)
	if err != nil {
		return 0, err
	}
	var sizedata metaData
	err = json.Unmarshal(contents, &sizedata)
	if err != nil {
		return 0, err
	}
	return sizedata.Size, nil
}

func (c UsageCollector) getUsedVolumes() ([]string, error) {
	linkDir := filepath.Join(c.storeDir, string(c.privilegeType), "l")
	usedVolumes := map[string]bool{}
	ids := []string{}

	err := filepath.Walk(linkDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink == 0 {
			return nil
		}
		target, err := os.Readlink(path)
		if err != nil {
			return err
		}
		volumeId := filepath.Base(target)
		if !usedVolumes[volumeId] {
			usedVolumes[volumeId] = true
			ids = append(ids, volumeId)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (c UsageCollector) usedVolumesSize() (int64, error) {
	var usedSize int64

	volumes, err := c.getUsedVolumes()
	if err != nil {
		return 0, err
	}

	for _, volumeId := range volumes {
		size, err := c.volumeSize(volumeId)
		if err != nil {
			return 0, err
		}
		usedSize += size
	}

	return usedSize, nil
}

func (c UsageCollector) totalVolumesSize() (int64, error) {
	output, err := c.runner.Run(context.Background(), "du", "-b", "-s", filepath.Join(c.storeDir, string(c.privilegeType), "volumes"))
	if err != nil {
		return 0, err
	}
	parts := strings.Split(string(output), "\t")
	if len(parts) < 2 {
		return 0, fmt.Errorf("unexpected `du` output %q", output)
	}
	size, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, err
	}
	return size, nil
}
