package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Result struct {
    url string
    source string
    status int
    message string
    body string
    contentType string
}

type DiscoveredUrl struct {
	url string
	source string
}

func findRoot(url string) string {
	url = url + "/" //just in case url has no leading slash. More than one won't harm
	re := regexp.MustCompile("https?://([^/]*/)")
	root := re.FindString(url)
	return root
}

func checkWebsite(url string, limit int, urlsToIgnore []string, statsChan chan Result) bool {
	resultChan := make(chan Result)
	urlDiscoveryChan := make(chan DiscoveredUrl)
	pendingChecks := 1
	count := 0
	knownUrls := make(map[string]string)

	for _, knownUrl := range urlsToIgnore {
    	knownUrls[knownUrl] = knownUrl
	}

	finishOrLimitChan := make(chan bool)

	urlRoot := findRoot(url)

	go func(finishOrLimitChan chan bool, urlDiscoveryChan chan DiscoveredUrl) {
		for {
			if count == limit {
				finishOrLimitChan <- true
				break
			}

			result := <- resultChan
			pendingChecks--
			statsChan <- result

			if strings.HasPrefix(result.url, urlRoot) && isHtmlContentType(result.contentType) {
				//parse page urls only if this page is on our domain
				newUrls := findUrls(result.body)
				for _,newUrl := range newUrls {
					if strings.HasPrefix(newUrl, "//") { //protocol relative url
						newUrl = urlRoot[0:strings.Index(urlRoot, ":") + 1] + newUrl
					}
					if string([]rune(newUrl)[0]) == "/" { //make sure we have an absolute URL
						newUrl = strings.Replace(newUrl, "/", urlRoot, 1)
					}
					if string([]rune(newUrl)[0]) == "#" { //skip if this is an internal link
						continue
					}
					if strings.HasPrefix(newUrl, "mailto:") { //skip email urls
						continue
					}
					if strings.HasPrefix(newUrl, "data:") { //skip email urls
						continue
					}

					if _, ok := knownUrls[newUrl]; ok {
						continue
					}
					if !strings.HasPrefix(newUrl, urlRoot) {
						//fmt.Println(newUrl)
						//continue
					}
					knownUrls[newUrl] = result.url
					if pendingChecks > 5  {
						time.Sleep(1)
					}
					if pendingChecks > 10 {
						time.Sleep(5)
                                        }
					pendingChecks++
					go checkUrl(newUrl, result.url, resultChan)
				}
			} else {
				//fmt.Println(result.url)
			}

			count++

			if pendingChecks == 0 {
				finishOrLimitChan <- true
				break
			}
		}
	}(finishOrLimitChan, urlDiscoveryChan)

	//init the first check
	go checkUrl(url, "", resultChan)

	<-finishOrLimitChan
	return true
}

func removeLineContent() {
	fmt.Print("\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b")
}


func main() {
	config := getConfiguration()

	fmt.Println("Website to be checked " + config.Url + ".")
	fmt.Println("Max count of URLs to be checked " + strconv.Itoa(config.Limit))

	count := 0
	errorCount := 0

	statsChan := make(chan Result)

	go func(statsChan chan Result) {
		for {
			result := <- statsChan
			count++

			removeLineContent()

			if result.status != 200 {
				errorCount++
				fmt.Println("Error: HTTP status " + strconv.Itoa(result.status) + ", Url " + result.url + ", Source " + result.source + " " + result.message)	
			}

			if config.DisplayProgress {
				fmt.Print("URLs checked " + strconv.Itoa(count))
			}

		}
	}(statsChan)
	
	checkWebsite(config.Url, config.Limit, config.UrlsToIgnore, statsChan)

	removeLineContent()	
	fmt.Println()
	fmt.Println("Total URLs checked: " + strconv.Itoa(count))

	if errorCount == 0 {
		fmt.Println("No broken URLs")
		os.Exit(0)
	}

	fmt.Println("Broken URLs: " + strconv.Itoa(errorCount))
	os.Exit(1)
}
