package main
import (
	"fmt"
	"net/http"
	"os"
	//"io"
	//"reflect"
	"io/ioutil"
	//"time"
	"regexp"
	"strconv"
	"crypto/tls"
	"log"
	"strings"
)

//routine

type Result struct {
    url string
    status int
    message string
    body string
}

type DiscoveredUrl struct {
	url string
	source string
}

func checkUrl(url string, resultChan chan Result, singleUrlFinishChan chan bool) {
	go func() {
		//fmt.Println("Checking url " + url)
		tr := &http.Transport{ //we ignore ssl errors. This tool is for testing 404, not ssl.
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}

		r, err := client.Get(url)
		if err != nil {
			resultChan <- Result{url: url, status: -1, message: "Fatal error " + err.Error(), body: ""}
		} else {
			defer r.Body.Close()
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				resultChan <- Result{url: url, status: -2, message: "Fatal error " + err.Error(), body: ""}
			} else {
				stringBody := fmt.Sprintf("%s", body)
				resultChan <- Result{url: url, status: r.StatusCode, message: "", body: stringBody}
  			}
		}
		singleUrlFinishChan <- true
	}()
}

func findUrls(html string) []string {
	re := regexp.MustCompile("<a .* href=\"(https://www.tgstatic.com/[^\"]*)")
	matches := re.FindAllStringSubmatch(html, -1)
	var urls = []string{}
	for _,match := range matches {
		urls = append(urls, match[1])
	}
	return urls
}

func findRoot(url string) string {
	return "https://www.tgstatic.com/"
}

func find404Errors(url string, limit int) bool {
	resultChan := make(chan Result)
	urlDiscoveryChan := make(chan DiscoveredUrl)
	finishChan := make(chan bool)
	pendingChecks := 1
	count := 0
	knownUrls := make(map[string]string)

	finishOrLimitChan := make(chan bool)

	urlRoot := findRoot(url)

	go func(finishOrLimitChan chan bool, urlDiscoveryChan chan DiscoveredUrl) {
		for {
			if count == limit {
				fmt.Println("Limit reached")
				finishOrLimitChan <- true
			}
			result := <-resultChan
			fmt.Println(strconv.Itoa(result.status) + " " + result.url + " (" + knownUrls[result.url] + ")" + " " + result.message)
			newUrls := findUrls(result.body)
			for _,newUrl := range newUrls {
				if _, ok := knownUrls[newUrl]; ok {
					continue
				}
				if !strings.HasPrefix(newUrl, urlRoot) {
					continue
				}
				knownUrls[newUrl] = result.url
				//fmt.Println("New URL detected " + newUrl)
				pendingChecks++
				go checkUrl(newUrl, resultChan, finishChan)
			}
			count++
		}
	}(finishOrLimitChan, urlDiscoveryChan)

	go func(finishChan <-chan bool, finishOrLimitChan chan bool) {
		for {
			<-finishChan
			pendingChecks--
			// fmt.Println("Check finished " + strconv.Itoa(pendingChecks))
			if pendingChecks == 0 {
				finishOrLimitChan <- true	
			}	
		}
	}(finishChan, finishOrLimitChan)

	//init the first check
	go checkUrl(url, resultChan, finishChan)

	<-finishOrLimitChan
	return true
}


func main() {
	if len(os.Args) != 3 {
		log.Fatal("Invalid number of arguments.\nUsage example:\n\"go run link-checker.go http://example.com 100\"")
	}
	url := os.Args[1]
 	limit, _ := strconv.Atoi(os.Args[2])
	results := find404Errors(url, limit)
	fmt.Println(results)
}
