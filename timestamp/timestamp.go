package timestamp

import (
	"fmt"

	"io"

	"github.com/sgravrock/flickr-to-go-go/clock"
	"github.com/sgravrock/flickr-to-go-go/storage"
)

func Read(fs storage.Storage, stderr io.Writer) uint32 {
	result, err := read2(fs, stderr)
	if err != nil {
		fmt.Fprintf(stderr, "Error reading timestamp: %s\n", err.Error())
		return 0
	}

	return result
}

func read2(fs storage.Storage, stderr io.Writer) (uint32, error) {
	f, err := fs.Open("timestamp")
	if err != nil {
		return 0, err
	}
	defer f.Close()
	var result uint32
	_, err = fmt.Fscan(f, &result)
	return result, err
}

func Write(clock clock.Clock, fs storage.Storage) error {
	f, err := fs.Create("timestamp")
	if err != nil {
		return err
	}
	defer f.Close()
	s := fmt.Sprint(clock.Now().Unix()) + "\n"
	f.Write([]byte(s))
	return nil
}
