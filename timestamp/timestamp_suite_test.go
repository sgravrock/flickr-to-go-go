package timestamp_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTimestamp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Timestamp Suite")
}
