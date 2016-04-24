package flickrapi_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestFlickrapi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Flickrapi Suite")
}
