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
	"strconv"
)

//routine

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
// func checkUrls(urls []string, resultChan chan Result, urlDiscoveryChan chan string, finishChan chan bool) {

// 	go func() {
// 		singleUrlFinishChan := make(chan bool)
 		
//  		for _,url := range urls { //initialize checking of all URLs
// 			go checkUrl(url, resultChan, urlDiscoveryChan, singleUrlFinishChan)
//   		}
//   		for i := 0; i < len(urls); i++ { //wait till all urls all hecked
//   			<-singleUrlFinishChan
//   		}
//   		finishChan <- true
// 	}()
// }

func checkUrl(url string, resultChan chan Result, urlDiscoveryChan chan string, singleUrlFinishChan chan bool) {
	go func() {
		//fmt.Println("Checking url " + url)
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
				for _,newUrl := range newUrls {
					urlDiscoveryChan <- newUrl
				}
				//checkUrls(newUrls, resultChan, finishChan)
				resultChan <- Result{url: url, status: r.StatusCode, message: ""}
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

func registerNewUrls(newUrls []string) {
	for _,newUrl := range newUrls {
		knownUrls[newUrl] = true
	}
}
func find404Errors(url string, limit int) bool {
	resultChan := make(chan Result)
	urlDiscoveryChan := make(chan string)
	finishChan := make(chan bool)
	pendingChecks := 1
	//checkUrls(urls, resultChan, urlDiscoveryChan, finishChan)
	count := 0

	finishOrLimitChan := make(chan bool)


	go func(urlDiscoveryChan chan string) {
		for {
			newUrl := <-urlDiscoveryChan
			//fmt.Println("New URL detected " + newUrl)
			pendingChecks++
			go checkUrl(newUrl, resultChan, urlDiscoveryChan, finishChan)
		}
	}(urlDiscoveryChan)

	go func(finishOrLimitChan chan bool) {
		for {
			if count == limit {
				fmt.Println("Limit reached")
				finishOrLimitChan <- true
			}
			result := <-resultChan
			fmt.Println(strconv.Itoa(result.status) + " " + result.url)
			count++
		}
	}(finishOrLimitChan)

	go func(finishChan <-chan bool, finishOrLimitChan chan bool) {
		<-finishChan
		pendingChecks--
		if pendingChecks == 0 {
			finishOrLimitChan <- true	
		}
	}(finishChan, finishOrLimitChan)

	//init the first check
	go checkUrl(url, resultChan, urlDiscoveryChan, finishChan)


	<-finishOrLimitChan
	return true
}


func main() {
	results := find404Errors("https://www.tgstatic.com/en", 100)
	fmt.Println(results)
}
