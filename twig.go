package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"net/http"
)

func crashAndBurn(err error) {
	fmt.Println(err)
	os.Exit(-1)
}

type Options struct {
	Url string
	Method string
}

func readArgs() (Options, error) {
	opts := Options {}

	flag.StringVar(&opts.Url, "U", "", "Specify URL.")
	flag.StringVar(&opts.Method, "m", "GET", "Speci")

	flag.Parse()

	if !flag.Parsed() {
		return opts, errors.New("Flags could not be parsed.")
	}

	if opts.Url == "" {
		return opts, errors.New("No URL supplied.")
	}
	
	return opts, nil
}

func main() {
	opts, err := readArgs()

	crashAndBurn(err)
	
	resp, err := http.Get(opts.Url)

	crashAndBurn(err)

	fmt.Printf("Status: %s\n", resp.Status)
	
	body, err := io.ReadAll(resp.Body)
	content := string(body[:])

	fmt.Printf("Content: %s\n", content)
}
