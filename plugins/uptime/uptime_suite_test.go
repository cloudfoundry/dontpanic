package uptime_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestUptime(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Uptime Suite")
}
