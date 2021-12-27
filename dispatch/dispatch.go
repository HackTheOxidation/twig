package dispatch

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"codeberg.org/HackTheOxidation/twig/options"
)

func crashAndBurn(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
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

func handleGET(resp *http.Response, opts *options.Options) {
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

func createConfiguredClient() *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &http.Client{Transport: transport}
}

func DispatchAndExecute(opts *options.Options) {
	var resp *http.Response
	var err error
	client := createConfiguredClient()

	switch strings.ToUpper(opts.Method) {
	case "GET":
		resp, err = client.Get(opts.Url)
		crashAndBurn(err)
		handleGET(resp, opts)
	default:
		crashAndBurn(errors.New("twig: ERROR - Unrecognized method: " + opts.Method))
	}
}