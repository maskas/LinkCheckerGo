package main
import (
	"fmt"
	"net/http"
	//"os"
	//"io"
	//"reflect"
	"io/ioutil"
	"unicode/utf8"
	//"time"
	"regexp"
)

type Result struct {
    url string
    status int
    message string
}

var (
	knownUrls = make(map[string]bool)
)

/*
Check an array of urls for 404 errors
first channel returns results
next channel returns true when all urls are checked.
*/
func checkUrls(urls []string, resultChan chan Result, urlDiscoveryChan chan string, finishChan chan bool) {

	go func() {
		internalFinishChan := make(chan bool)
 		
 		for _,url := range urls { //initialize checking of all URLs
 			fmt.Println(url)
			go checkUrl(url, resultChan, urlDiscoveryChan, internalFinishChan)
  		}
  		for i := 0; i < len(urls); i++ { //wait till all urls all hecked
  			<-internalFinishChan
  		}
  		finishChan <- true
	}()
}

func checkUrl(url string, resultChan chan Result, urlDiscoveryChan chan string, finishChan chan bool) {
	go func() {
		fmt.Println("Checking url " + url)
		r, err := http.Get(url)
		if err != nil {
			resultChan <- Result{url: url, status: 0, message: "Fatal error " + err.Error()}
		} else {
			defer r.Body.Close()
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				resultChan <- Result{url: url, status: 0, message: "Fatal error " + err.Error()}
			} else {
				stringBody := fmt.Sprintf("%s", body)
				utf8.RuneCountInString(stringBody)
				newUrls := findUrls(stringBody)
				registerNewUrls(newUrls)
				//checkUrls(newUrls, resultChan, finishChan)
				fmt.Println("Result " + url)
				resultChan <- Result{url: url, status: r.StatusCode, message: ""}
  			}
		}
		finishChan <- true
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

func registerNewUrls(newUrls []string) {
	for _,newUrl := range newUrls {
		knownUrls[newUrl] = true
	}
}
func startChecking(urls []string, limit int) bool {
	resultChan := make(chan Result)
	urlDiscoveryChan := make(chan string)
	finishChan := make(chan bool)
	checkUrls(urls, resultChan, urlDiscoveryChan, finishChan)
	count := 0

	finishOrLimitChan := make(chan bool)

	go func(finishOrLimitChan chan bool) {
		for {
			if count == limit {
				fmt.Println("Limit reached")
				finishOrLimitChan <- true
			}
			fmt.Println("Listening")
			result := <-resultChan
			fmt.Println(result)
			count++
		}
	}(finishOrLimitChan)

	// go	 func(finishChan <-chan bool, finishOrLimitChan chan bool) {
	// 	<-finishChan
	// 	finishOrLimitChan <- true
	// }(finishChan, finishOrLimitChan)


	<-finishOrLimitChan
	return true
}


func main() {
	urls := []string{"https://www.tgstatic.com/lt", "https://www.tgstatic.com/en"}
	results := startChecking(urls, 10)
	fmt.Println(results)
}
