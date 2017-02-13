# Download
Download helps with download files and keeping track of where you have downloaded them. It's designed to be embedded in another struct. Composition ahoy!

Imagine you have a list of books like this:

```go
type Book struct {
	Name string
	Author string
	URL string
}
```

You may want to download the book from the URL, but also keep track of where you've downloaded it.
Well, that's what this library is for.

Go get it:

```bash
$ go get github.com/codeclysm/download
```

And embed it in your struct

```go
type Book struct {
	download.Resource
	Author string
}

book := Book{}
book.Name = "The silly lives of gophers"
book.Author = "G. O. Fer"
book.URL = "https://download.example.com/silly_gophers.epub"

book.Download("books", nil)

log.Println(book.Where()) // Will print out ["books"]
```

Download comes with a few options you may want to use:

```go
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
```

Check out the documentation (I just love godoc): https://godoc.org/github.com/codeclysm/download
