package osreporter_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOsreporter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Osreporter Suite")
}
