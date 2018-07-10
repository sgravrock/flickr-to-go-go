# flickr-to-go-go
Same idea as flickr-to-go, but written in Go for learning and comparison.

## Setup
1. Install Ginkgo: <http://onsi.github.io/ginkgo/>.
2. Clone the repo and ensure that the `GOPATH` environment variableis set appropriately as for any other Go project.
3. Obtain a Flickr API key and set the `FLICKR_API_KEY` and `FLICKR_API_SECRET` environment variables accordingly.
4. To run the tests, run `ginkgo -r` from the repo root.
5. To build, run `go build`.

## Generating new fakes, or updating existing ones
1. Install Counterfeiter:
    $ go get github.com/maxbrunsfeld/counterfeiter
2. Run it:
    $ counterfeiter <package path> <interface name>
    e.g.
    $ counterfeiter flickrapi Client