package download_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"path/filepath"
	"testing"

	download "github.com/codeclysm/downloader"
)

func TestDownload(t *testing.T) {
	cases := []struct {
		File string
	}{
		{"test1.txt"}, {"test2.txt"},
	}
	for _, tc := range cases {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, path.Join("testdata", tc.File))
		}))
		defer ts.Close()

		d := download.Resource{URL: ts.URL, Name: tc.File}

		temp, err := ioutil.TempDir("", "")
		if err != nil {
			t.Fatal(err)
		}
		err = d.Download(temp)
		if err != nil {
			t.Fatal(err)
		}

		data1, err := ioutil.ReadFile(filepath.Join(temp, tc.File))
		if err != nil {
			t.Fatal(err)
		}

		data2, err := ioutil.ReadFile(filepath.Join("testdata", tc.File))
		if err != nil {
			t.Fatal(err)
		}

		if string(data1) != string(data2) {
			t.Errorf("The file wasn't downloaded correctly")
		}
	}
}
