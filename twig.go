package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func crashAndBurn(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

type Options struct {
	Url    string
	Method string
	Unsafe bool
	Output string
}

func (opts *Options) Print() {
	fmt.Printf("Options: { \n\tURL: %s, ", opts.Url)
	fmt.Printf("\n\tMethod: %s, ", opts.Method)
	fmt.Printf("\n\tUnsafe: %t, ", opts.Unsafe)
	fmt.Printf("\n\tOutput: %s\n}\n", opts.Output)
}

func parseOutput(opts *Options) {
	if opts.Output == "" {
		if strings.Contains(opts.Url, "/") {
			opts.Output = strings.SplitN(opts.Url, "/", 1)[1]
		} else {
			opts.Output = "index.html"
		}
	}
}

func readArgs() (Options, error) {
	opts := Options{}

	flag.StringVar(&(opts.Output), "o", "", "Specify output file name.")
	flag.StringVar(&(opts.Method), "m", "GET", "Specify HTTP method.")
	flag.BoolVar(&(opts.Unsafe), "unsafe", false, "Specify security of the connection.")

	flag.Parse()

	if !flag.Parsed() {
		return opts, errors.New("twig: ERROR - Flags could not be parsed.")
	}

	opts.Url = flag.Arg(0)

	if opts.Url == "" {
		return opts, errors.New("twig: ERROR - No URL supplied.")
	}

	parseOutput(&opts)

	return opts, nil
}

func ensureProtocol(opts *Options) {
	if !strings.Contains(opts.Url, "http") {
		if opts.Unsafe {
			opts.Url = "http://" + opts.Url
		} else {
			opts.Url = "https://" + opts.Url
		}
	}
}

func getOptions() Options {
	opts, err := readArgs()
	crashAndBurn(err)
	ensureProtocol(&opts)

	return opts
}

func writeOutputToFile(content *string, filename string) {
	fp, err := os.Create(filename)
	crashAndBurn(err)

	l, err := fp.WriteString(*content)
	crashAndBurn(err)
	fmt.Println("SUCCESS - ", l, "bytes written to file:", filename)
	err = fp.Close()
	crashAndBurn(err)
}

func handleGET(resp *http.Response, opts *Options) {
	if resp.StatusCode != 200 {
		crashAndBurn(errors.New("twig: unsuccessful request - received response with status: " + resp.Status))
	}
	body, err := io.ReadAll(resp.Body)
	crashAndBurn(err)
	content := string(body[:])

	if opts.Output != "" {
		writeOutputToFile(&content, opts.Output)
	} else {
		fmt.Println(content)
	}
}

func dispatchAndExecute(opts *Options) {
	var resp *http.Response
	var err error

	switch strings.ToUpper(opts.Method) {
	case "GET":
		resp, err = http.Get(opts.Url)
		crashAndBurn(err)
		handleGET(resp, opts)
	default:
		crashAndBurn(errors.New("twig: ERROR - Unrecognized method: " + opts.Method))
	}
}

func main() {
	opts := getOptions()
	dispatchAndExecute(&opts)
}
