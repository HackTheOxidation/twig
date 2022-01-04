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
	fmt.Println("twig: SUCCESS - ", l, "bytes written to file:", filename)
	err = fp.Close()
	crashAndBurn(err)
}

func handleGET(opts *options.Options, client *http.Client) {
	resp, err := client.Get(opts.Url)
	crashAndBurn(err)

	if resp.StatusCode != 200 {
		crashAndBurn(errors.New("twig: ERROR - Unsuccessful request. Received response with status: " + resp.Status))
	}
	body, err := io.ReadAll(resp.Body)
	crashAndBurn(err)
	content := string(body[:])

	if opts.Output != "" {
		writeOutputToFile(&content, opts.Output)
	} else if content != "" {
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

	fmt.Println("twig: INFO - Received response with status: ", resp.Status)
}

func handlePUT(opts *options.Options, client *http.Client) {
	if opts.Content == "" {
		crashAndBurn(errors.New("twig: ERROR - Content cannot be empty when using method: PUT"))
	}

	payload, err := json.Marshal(opts.Content)
	crashAndBurn(err)

	req, err := http.NewRequest(http.MethodPut, opts.Url, bytes.NewBuffer(payload))
	crashAndBurn(err)

	resp, err := client.Do(req)
	crashAndBurn(err)

	fmt.Println("twig: INFO - Received response with status: ", resp.Status)
}

func handleDELETE(opts *options.Options, client *http.Client) {
	payload, err := json.Marshal(opts.Content)
	crashAndBurn(err)

	req, err := http.NewRequest(http.MethodDelete, opts.Url, bytes.NewBuffer(payload))
	crashAndBurn(err)

	resp, err := client.Do(req)
	crashAndBurn(err)

	fmt.Println("twig: INFO - Received response with status: ", resp.Status)

	body, err := io.ReadAll(resp.Body)
	crashAndBurn(err)
	content := string(body[:])

	if opts.Output != "" {
		writeOutputToFile(&content, opts.Output)
	} else if content != "" {
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
	client := createConfiguredClient()

	switch strings.ToUpper(opts.Method) {
	case "GET":
		handleGET(opts, client)
	case "POST":
		handlePOST(opts, client)
	case "PUT":
		handlePUT(opts, client)
	case "DELETE":
		handleDELETE(opts, client)
	default:
		crashAndBurn(errors.New("twig: ERROR - Unrecognized method: " + opts.Method))
	}
}
