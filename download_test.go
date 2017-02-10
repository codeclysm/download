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
		err = d.Download(temp, nil)
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

func TestWhere(t *testing.T) {
	cases := []struct {
		File string
		N    int
	}{
		{"test1.txt", 0}, {"test2.txt", 2},
	}
	for _, tc := range cases {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, path.Join("testdata", tc.File))
		}))
		defer ts.Close()

		d := download.Resource{URL: ts.URL, Name: tc.File}

		if len(d.Where()) != 0 {
			t.Errorf("d.Where() should be empty: %s", d.Where())
		}

		where := []string{}

		for i := 0; i < tc.N; i++ {
			temp, err := ioutil.TempDir("", "")
			if err != nil {
				t.Fatal(err)
			}

			where = append(where, temp)
			err = d.Download(temp, nil)
			if err != nil {
				t.Fatal(err)
			}
		}

		if len(d.Where()) != tc.N {
			t.Errorf("d.Where() should have length %d: %s", tc.N, d.Where())
		}

		for _, el := range where {
			if !in(d.Where(), el) {
				t.Errorf("%s should be in %s", el, d.Where())
			}
		}
		for _, el := range d.Where() {
			if !in(where, el) {
				t.Errorf("%s should be in %s", el, where)
			}
		}
	}
}

func TestCache(t *testing.T) {
	downloaded := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if downloaded {
			t.Errorf("It shouldn't download it again")
		}
		downloaded = true
		http.ServeFile(w, r, path.Join("testdata", "test1.txt"))
	}))
	defer ts.Close()

	d := download.Resource{URL: ts.URL, Name: "test1.txt"}
	temp, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	d.Download(temp, &download.Opts{Cache: true})
	d.Download(temp, &download.Opts{Cache: true})
}

func in(slice []string, el string) bool {
	for i := range slice {
		if slice[i] == el {
			return true
		}
	}
	return false
}
