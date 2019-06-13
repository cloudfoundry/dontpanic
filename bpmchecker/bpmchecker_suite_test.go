package bpmchecker_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBpmchecker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bpmchecker Suite")
}
