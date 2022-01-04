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
	Content string
}

func (opts *Options) Print() {
	fmt.Printf("Options: { \n\tURL: %s, ", opts.Url)
	fmt.Printf("\n\tMethod: %s, ", opts.Method)
	fmt.Printf("\n\tUnsafe: %t, ", opts.Unsafe)
	fmt.Printf("\n\tOutput: %s, ", opts.Output)
	fmt.Printf("\n\tContent: %s\n}\n", opts.Content)
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

func displayHelp() {
	fmt.Println("Usage: twig <options> [url]")
	fmt.Println("\noptions:")
	flag.PrintDefaults()
	os.Exit(0)
}

func readArgs() (Options, error) {
	opts := Options{}

	flag.StringVar(&(opts.Output), "o", "", "Specify output file name.")
	flag.StringVar(&(opts.Method), "m", "GET", "Specify HTTP method.")
	flag.BoolVar(&(opts.Unsafe), "u", false, "Specify security of the connection.")
	flag.StringVar(&(opts.Content), "c", "", "Specify request content.")
	help := flag.Bool("h", false, "Display help information.")

	flag.Parse()

	if !flag.Parsed() {
		return opts, errors.New("twig: ERROR - Flags could not be parsed.")
	}

	if *help {
		displayHelp()
	}

	opts.Url = flag.Arg(0)

	if opts.Url == "" {
		return opts, errors.New("twig: ERROR - No URL supplied.")
	}

	return opts, nil
}

func GetOptions() Options {
	opts, err := readArgs()
	crashAndBurn(err)
	opts.ensureProtocol()

	return opts
}
