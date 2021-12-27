package options

import (
	"errors"
	"flag"
	"fmt"
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

func (opts *Options) ensureOutput() {
	if opts.Output == "" {
		if strings.Contains(opts.Url, "/") {
			opts.Output = strings.SplitN(opts.Url, "/", 1)[1]
		} else {
			opts.Output = "index.html"
		}
	}
}

func (opts *Options) ensureProtocol() {
	if !strings.Contains(opts.Url, "http") {
		if opts.Unsafe {
			opts.Url = "http://" + opts.Url
		} else {
			opts.Url = "https://" + opts.Url
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

	opts.ensureOutput()

	return opts, nil
}

func GetOptions() Options {
	opts, err := readArgs()
	crashAndBurn(err)
	opts.ensureProtocol()

	return opts
}