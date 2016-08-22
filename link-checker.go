package main
import (
	"fmt"
	"net/http"
	"os"
	"io/ioutil"
	"regexp"
	"strconv"
	"crypto/tls"
	"log"
	"strings"
)

//routine

type Result struct {
    url string
    source string
    status int
    message string
    body string
}

type DiscoveredUrl struct {
	url string
	source string
}

func checkUrl(url string, source string, resultChan chan Result, singleUrlFinishChan chan bool) {
	go func() {
		tr := &http.Transport{ //we ignore ssl errors. This tool is for testing 404, not ssl.
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}

		r, err := client.Get(url)
		if err != nil {
			resultChan <- Result{url: url, source: source, status: -1, message: "Fatal error " + err.Error(), body: ""}
		} else {
			defer r.Body.Close()
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				resultChan <- Result{url: url, source: source, status: -2, message: "Fatal error " + err.Error(), body: ""}
			} else {
				stringBody := fmt.Sprintf("%s", body)
				resultChan <- Result{url: url, source: source, status: r.StatusCode, message: "", body: stringBody}
  			}
		}
		singleUrlFinishChan <- true
	}()
}

func findUrls(html string) []string {
	re := regexp.MustCompile("<a .* href=\"([^\"]*)")
	matches := re.FindAllStringSubmatch(html, -1)
	var urls = []string{}
	for _,match := range matches {
		urls = append(urls, match[1])
	}
	return urls
}

func findRoot(url string) string {
	url = url + "/" //just in case url has no leading slash. More than one won't harm
	re := regexp.MustCompile("https?://([^/]*/)")
	root := re.FindString(url)
	return root
}

func checkWebsite(url string, limit int, statsChan chan Result) bool {
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
				fmt.Println("The imit of " + strconv.Itoa(limit) + " urls has been reached.")
				finishOrLimitChan <- true
				break
			}
			result := <- resultChan
			statsChan <- result
			newUrls := findUrls(result.body)
			for _,newUrl := range newUrls {
				if string([]rune(newUrl)[0]) == "/" { //make sure we have an absolute URL
					newUrl = strings.Replace(newUrl, "/", urlRoot, 1)
				}
				if _, ok := knownUrls[newUrl]; ok {
					continue
				}
				if !strings.HasPrefix(newUrl, urlRoot) {
					continue
				}
				knownUrls[newUrl] = result.url
				pendingChecks++
				go checkUrl(newUrl, result.url, resultChan, finishChan)
			}
			count++
			fmt.Print("\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b")
			fmt.Print("Checked URLs: " + strconv.Itoa(count))
		}
	}(finishOrLimitChan, urlDiscoveryChan)

	go func(finishChan <-chan bool, finishOrLimitChan chan bool) {
		for {
			<-finishChan
			pendingChecks--
			if pendingChecks == 0 {
				finishOrLimitChan <- true	
			}	
		}
	}(finishChan, finishOrLimitChan)

	//init the first check
	go checkUrl(url, "", resultChan, finishChan)

	<-finishOrLimitChan
	return true
}

func removeLineContent() {
	fmt.Print("\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b")
}


func main() {
	if len(os.Args) != 3 {
		log.Fatal("Invalid number of arguments.\nUsage example:\n\"go run link-checker.go http://example.com 100\"")
	}
	url := os.Args[1]
	statsChan := make(chan Result)
 	limit, _ := strconv.Atoi(os.Args[2])

	count := 0

	fmt.Println("Website to be checked " + url + ".")
	fmt.Println("Max count of URLs to be checked " + strconv.Itoa(limit))


	go func(statsChan chan Result) {
		for {
			result := <- statsChan
			count++
			removeLineContent()
			if result.status != 200 {
				fmt.Println("Error: HTTP status " + strconv.Itoa(result.status) + ", Url " + result.url + ", Source " + result.source + " " + result.message)
			}
			fmt.Print("URLs checked " + strconv.Itoa(count))
		}
	}(statsChan)

	results := checkWebsite(url, limit, statsChan)
	fmt.Println(results)
}
