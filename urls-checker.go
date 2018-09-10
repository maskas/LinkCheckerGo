package main

import (
	"net/http"
	"crypto/tls"
	"io/ioutil"
	"fmt"
	"strings"
)

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.116 Safari/537.36"

func checkUrl(url string, source string, resultChan chan Result) {
	go func() {
		req, err := createRequest(url)

		if err != nil {
			resultChan <- Result{
				url: url,
				source: source,
				status: -2,
				message: "Fatal error " + err.Error(),
				body: "",
			}

			return
		}

		r, err := doRequest(req)

		if err != nil {
			resultChan <- Result{
				url: url,
				source: source,
				status: -1,
				message: "Fatal error " + err.Error(),
				body: "",
			}

			return
		}

		contentType := r.Header["Content-Type"][0]

		result := Result{
			url: url,
			source: source,
			status: r.StatusCode,
			message: "",
			body: "",
			contentType: contentType,
		}

		if false == isHtmlContentType(contentType) {
			resultChan <- result

			return
		}

		body, err := getRequestBody(r)

		if err != nil {
			result.status = -2
			result.message = "Fatal error " + err.Error()

			resultChan <- result

			return
		}

		result.body = fmt.Sprintf("%s", body)

		resultChan <- result
	}()
}

func createRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return req, err
	}

	req.Header.Set("User-Agent", userAgent)

	return req, nil
}

func doRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{
		Transport: &http.Transport{ //we ignore ssl errors. This tool is for testing 404, not ssl.
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	r, err := client.Do(req)

	if err != nil || r.StatusCode == 504 {
		//retrying in a simple way in case we had any network issues
		r, err = client.Do(req)
		if err != nil || r.StatusCode == 504 {
			//retry
			r, err = client.Do(req)
		}
	}

	return r, err
}

func isHtmlContentType(header string) bool {
	return strings.Contains(header, "text/html")
}

func getRequestBody(r *http.Response) ([]byte, error) {
	defer r.Body.Close()

	return ioutil.ReadAll(r.Body)
}
