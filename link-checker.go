package main
import (
	"fmt"
	"net/http"
	"log"
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
}

var (
	knownUrls = make(map[string]bool)
)

/*
Check an array of urls for 404 errors
first channel returns results
next channel returns true when all urls are checked.
*/
func checkUrls(urls []string) (<-chan Result, <-chan bool) {
	resultChan := make(chan Result)
	finishChan := make(chan bool)

	go func() {
		internalResultChan := make(chan Result)
		internalFinishChan := make(chan bool)

		go func() { //pass trhough all results
			result := <-internalResultChan
			resultChan <- result
			}()
 		
 		for _,url := range urls { //initialize checking of all URLs
			checkUrl(url, internalResultChan, internalFinishChan)
  		}
  		for i := 0; i < len(urls); i++ { //wait till all urls all hecked
  			<-internalFinishChan
  		}
  		finishChan <- true
	}()


  	return resultChan, finishChan
}

func checkUrl(url string, resultChan chan Result, finishChan chan bool) {
	go func() {
		r, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
			fmt.Println("Error")
		} else {
			defer r.Body.Close()

			body, err := ioutil.ReadAll(r.Body)

			if err != nil {
				fmt.Println(err)
			} else {
				stringBody := fmt.Sprintf("%s", body)
				utf8.RuneCountInString(stringBody)
				newUrls := findUrls(stringBody)
				registerNewUrls(newUrls)
				// for _,newUrl := range newUrls {
				// 	fmt.Println(newUrl)
				// }
  			}
			fmt.Println(url + " Succeeded")
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



func main() {
	urls := []string{"https://www.tgstatic.com/lt", "https://www.tgstatic.com/en"}
	resultChan, finishChan := checkUrls(urls)
	go func() {
		result := <-resultChan
		fmt.Println(result)
	}()

	registerNewUrls([]string {"aaa", "bbb"})
	<-finishChan
	fmt.Println(knownUrls)
}
