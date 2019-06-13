package bpmchecker

import (
	"bytes"
	"io/ioutil"
)

type Checker struct {
	File string
}

func New(file ...string) Checker {
	path := "/proc/1/cmdline"
	if len(file) > 0 {
		path = file[0]
	}
	return Checker{File: path}
}

func (c Checker) HasGardenPid1() (bool, error) {
	contents, err := ioutil.ReadFile(c.File)
	if err != nil {
		return false, err
	}

	if bytes.Contains(contents, []byte("garden_start")) {
		return true, nil
	}

	return false, nil
}
