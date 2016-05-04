package dl_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDl(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dl Suite")
}
