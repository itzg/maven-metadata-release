package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type Metadata struct {
	Versioning struct {
		Release string `xml:"release"`
	} `xml:"versioning"`
}

func main() {
	var in io.Reader

	if len(os.Args) > 1 {
		location, err := url.Parse(os.Args[1])
		if err == nil {
			response, err := http.Get(location.String())
			if err != nil {
				failWithError("unable to retrieve maven metadata: %s", err)
			}
			defer response.Body.Close()
			in = response.Body
		} else {
			file, err := os.Open(os.Args[1])
			if err != nil {
				failWithError("unable to open input file %s: %s", os.Args[1], err)
			}
			defer file.Close()

			in = file
		}
	} else {
		in = os.Stdin
	}

	decoder := xml.NewDecoder(in)
	var metadata Metadata
	err := decoder.Decode(&metadata)
	if err != nil {
		failWithError("unable to decode input: %s", err)
	}

	fmt.Println(metadata.Versioning.Release)
}

func failWithError(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}
