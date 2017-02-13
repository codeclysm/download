package download_test

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"path/filepath"
	"testing"

	download "github.com/codeclysm/downloader"
)

func Example() {
	img := struct {
		download.Resource
		Author string
	}{}

	img.Name = "Mona Lisa"
	img.Author = "Leonardo da Vinci"
	img.URL = "https://upload.wikimedia.org/wikipedia/commons/6/6a/Mona_Lisa.jpg"

	fmt.Println(img.Download("paintings", nil))
	fmt.Println(img.Where())
	// Output:
	// <nil>
	// [paintings]
}

func Example_Cache() {
	img := struct {
		download.Resource
		Author string
	}{}

	img.Name = "Mona Lisa"
	img.Author = "Leonardo da Vinci"
	img.URL = "https://upload.wikimedia.org/wikipedia/commons/6/6a/Mona_Lisa.jpg"

	fmt.Println(img.Download("paintings", nil))

	// This time it won't be downloaded again
	fmt.Println(img.Download("paintings", &download.Opts{Cache: true}))
	fmt.Println(img.Where())
	// Output:
	// <nil>
	// <nil>
	// [paintings]
}

func Example_Sha256Sum() {
	img := struct {
		download.Resource
		Author string
	}{}

	img.Name = "Mona Lisa"
	img.Author = "Leonardo da Vinci"
	img.URL = "https://upload.wikimedia.org/wikipedia/commons/6/6a/Mona_Lisa.jpg"

	// It fails because the checksum is wrong
	fmt.Println(img.Download("paintings", &download.Opts{Sha256Sum: "wrong checksum"}))
	fmt.Println(img.Where())
	// Output:
	// Sha256sum check failed
	// []
}

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

func TestSha256Sum(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join("testdata", "test1.txt"))
	}))
	defer ts.Close()

	d := download.Resource{URL: ts.URL, Name: "test1.txt"}
	temp, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	err = d.Download(temp, &download.Opts{Cache: true, Sha256Sum: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestHandler(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join("testdata", "test1.txt"))
	}))
	defer ts.Close()

	d := download.Resource{URL: ts.URL, Name: "test1.txt"}
	temp, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	err = d.Download(temp, &download.Opts{Handler: func(body io.Reader, name, location string) error {
		return errors.New("")
	}})
	if err == nil {
		t.Fatal("Err should not be nil")
	}
}

func in(slice []string, el string) bool {
	for i := range slice {
		if slice[i] == el {
			return true
		}
	}
	return false
}
