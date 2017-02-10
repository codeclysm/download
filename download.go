// Download this thing here
// Download this thing here, here's how to process it
// Download this thing here, only if it matches this sha1sum
// Download this thing here, but only if you haven't do it already
// Have I already downloaded this thing here?
// Where did I download this thing?

package download

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/juju/errors"
)

type Resource struct {
	URL    string
	Name   string
	Client *http.Client
}

func (r *Resource) Download(location string) error {
	if r.Client == nil {
		r.Client = &http.Client{}
	}
	resp, err := r.Client.Get(r.URL)
	if err != nil {
		return errors.Annotatef(err, "Get %s", r.URL)
	}

	err = os.MkdirAll(location, 0755)
	if err != nil {
		return errors.Annotatef(err, "Create %s", location)
	}

	file, err := os.OpenFile(filepath.Join(location, r.Name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return errors.Annotatef(err, "Create %s", filepath.Join(location, r.Name))
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return errors.Annotatef(err, "Write %s", filepath.Join(location, r.Name))
	}

	return nil
}
