package dispatch

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	//"codeberg.org/HackTheOxidation/twig/options"
	"options"
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

func handleGET(client *http.Client, opts *options.Options) {
	resp, err := client.Get(opts.Url)
	crashAndBurn(err)

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

func handlePOST(opts *options.Options, client *http.Client) {
	if opts.Content == "" {
		crashAndBurn(errors.New("twig: ERROR - Content cannot be empty when using method: POST"))
	}

	payload, err := json.Marshal(opts.Content)
	crashAndBurn(err)

	resp, err := client.Post(opts.Url, "application/json", bytes.NewBuffer(payload))
	crashAndBurn(err)

	fmt.Println("twig: Received response with status: ", resp.Status)
}

func createConfiguredClient() *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &http.Client{Transport: transport}
}

func DispatchAndExecute(opts *options.Options) {
	client := createConfiguredClient()

	switch strings.ToUpper(opts.Method) {
	case "GET":
		handleGET(client, opts)
	case "POST":
		handlePOST(opts, client)
	default:
		crashAndBurn(errors.New("twig: ERROR - Unrecognized method: " + opts.Method))
	}
}
