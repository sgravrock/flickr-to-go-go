# flickr-to-go-go
Same idea as flickr-to-go, but written in Go for learning and comparison.

## Setup
1. Install [Ginkgo](http://onsi.github.io/ginkgo/).
2. Clone the repo and ensure that the `GOPATH` environment variableis set appropriately as for any other Go project.
  * If you've set up your Go environment "by the book", this is as simple as running `go get github.com/sgravrock/flickr-to-go-go`.
3. Obtain a Flickr API key and set the `FLICKR_API_KEY` and `FLICKR_API_SECRET` environment variables accordingly.
4. To run the tests, run `ginkgo -r` from the repo root.
5. To build, run `go build`.
