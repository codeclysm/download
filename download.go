// Download this thing here
// Download this thing here, here's how to process it
// Download this thing here, only if it matches this sha1sum
// Download this thing here, but only if you haven't do it already
// Have I already downloaded this thing here?
// Where did I download this thing?

package download

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/juju/errors"
)

type Resource struct {
	URL  string
	Name string

	where []string
}

type Opts struct {
	Client    *http.Client
	Cache     bool
	Sha256Sum string
}

func (r *Resource) Download(location string, opts *Opts) error {
	if opts == nil {
		opts = &Opts{}
	}

	// Check if already downloaded
	if opts.Cache && in(r.where, location) {
		return nil
	}

	// Check the http Client
	if opts.Client == nil {
		opts.Client = &http.Client{}
	}

	// Request the file
	resp, err := opts.Client.Get(r.URL)
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
		_, err := io.Copy(hasher, &buf)
		if err != nil {
			return errors.Annotate(err, "Calculate Sha256Sum")
		}
		if opts.Sha256Sum != hex.EncodeToString(hasher.Sum(nil)) {
			return errors.New("Sha256sum check failed")
		}
	}

	err = os.MkdirAll(location, 0755)
	if err != nil {
		return errors.Annotatef(err, "Create %s", location)
	}

	file, err := os.OpenFile(filepath.Join(location, r.Name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return errors.Annotatef(err, "Create %s", filepath.Join(location, r.Name))
	}

	_, err = io.Copy(file, body)
	if err != nil {
		return errors.Annotatef(err, "Write %s", filepath.Join(location, r.Name))
	}

	r.where = append(r.where, location)

	return nil
}

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
