// Package download helps with download files and keeping track of where you have downloaded them.
// It's designed to be embedded in another struct. Composition ahoy!
package download

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/juju/errors"
)

// Resource is an embeddable struct that keep tracks of where a resource is being downloaded.
type Resource struct {
	URL  string
	Name string

	where []string
}

// Opts is a struct of options to be passed to the Download function
type Opts struct {
	// Client is the http client used to fetch the resource
	Client *http.Client
	// Cache is a flag that tells the Download func not to download twice the resource in the same location
	Cache bool
	// Sha256Sum is the sum that's used to check that the download was correct
	Sha256Sum string
	// Handler is the function that saves or extracts the resource downloaded
	Handler func(body io.Reader, name, location string) error
}

// Download will retrieve the resource at .URL and save it on disk. Its behaviour
// can be modified with some options
func (r *Resource) Download(location string, opts *Opts) error {
	if opts == nil {
		opts = &Opts{}
	}

	// Check if already downloaded
	if opts.Cache && in(r.where, location) {
		return nil
	}

	client := &http.Client{Timeout: 1 * time.Second}
	// Check the http Client
	if opts.Client != nil {
		client = opts.Client
	}

	// Request the file
	resp, err := client.Get(r.URL)
	if err != nil {
		return errors.Annotatef(err, "Get %s", r.URL)
	}
	defer resp.Body.Close()

	var body io.Reader
	body = resp.Body

	// Checksum
	if opts.Sha256Sum != "" {
		var buf bytes.Buffer
		body = io.TeeReader(body, &buf)
		hasher := sha256.New()

		_, err := io.Copy(hasher, body)
		if err != nil {
			return errors.Annotate(err, "Calculate Sha256Sum")
		}

		if opts.Sha256Sum != hex.EncodeToString(hasher.Sum(nil)) {
			return errors.New("Sha256sum check failed")
		}

		body = &buf
	}

	// Execute function
	if opts.Handler != nil {
		err = opts.Handler(body, r.Name, location)
	} else {
		err = defaultHandler(body, r.Name, location)
	}

	if err != nil {
		return errors.Annotate(err, "Handling")
	}

	// Save the location
	r.where = append(r.where, location)

	return nil
}

func defaultHandler(body io.Reader, name, location string) error {
	err := os.MkdirAll(location, 0755)
	if err != nil {
		return errors.Annotatef(err, "Create %s", location)
	}

	file, err := os.OpenFile(filepath.Join(location, name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return errors.Annotatef(err, "Create %s", filepath.Join(location, name))
	}

	_, err = io.Copy(file, body)
	if err != nil {
		return errors.Annotatef(err, "Write %s", filepath.Join(location, name))
	}
	return nil
}

// Where returns a list of all the places where the resource was downloaded to
func (r *Resource) Where() []string {
	return r.where
}

func in(slice []string, el string) bool {
	for i := range slice {
		if slice[i] == el {
			return true
		}
	}
	return false
}
