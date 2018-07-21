package timestamp_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"

	"path"

	"bytes"

	"time"

	"github.com/sgravrock/flickr-to-go-go/clock/clockfakes"
	"github.com/sgravrock/flickr-to-go-go/storage"
	"github.com/sgravrock/flickr-to-go-go/timestamp"
)

var _ = Describe("Timestamp", func() {
	var (
		dir    string
		stderr *bytes.Buffer
		fs     storage.Storage
	)

	BeforeEach(func() {
		d, err := ioutil.TempDir("", "storage_test")
		dir = d
		Expect(err).To(BeNil())
		stderr = new(bytes.Buffer)
		fs = storage.NewFileStorage(dir)
	})

	Describe("Read", func() {
		var result uint32

		JustBeforeEach(func() {
			result = timestamp.Read(fs, stderr)
		})

		Describe("When the timestamp file exists", func() {
			BeforeEach(func() {
				err := ioutil.WriteFile(path.Join(dir, "timestamp"),
					[]byte("12345\n"), 0600)
				Expect(err).To(BeNil())
			})

			It("returns the stored timestamp", func() {
				Expect(result).To(Equal(uint32(12345)))
				Expect(stderr.Len()).To(Equal(0))
			})
		})

		Describe("When the timestamp file does not exist", func() {
			It("returns 0", func() {
				Expect(result).To(Equal(uint32(0)))
			})

			It("logs an error", func() {
				Expect(stderr.String()).To(HavePrefix("Error reading timestamp: open"))
			})
		})
	})

	Describe("Write", func() {
		It("writes the timestamp", func() {
			clock := new(clockfakes.FakeClock)
			clock.NowReturns(time.Unix(1257894000, 0))

			err := timestamp.Write(clock, fs)
			Expect(err).To(BeNil())

			b, err := ioutil.ReadFile(path.Join(dir, "timestamp"))
			Expect(err).To(BeNil())
			Expect(string(b)).To(Equal("1257894000\n"))
		})
	})
})
