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
)

type Result struct {
    url string
    status int
}


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
  			fmt.Println("Single Url finished")
  		}
  		finishChan <- true
	}()


  	return resultChan, finishChan
}

func checkUrl(url string, resultChan chan Result, finishChan chan bool) {


	go func() {
		fmt.Println("Start")
		//	r, err := http.Get("https://www.transfergo.com/lt")
		r, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
			fmt.Println("Error")
		} else {
			fmt.Println("Success")
			
			defer r.Body.Close()

			body, err := ioutil.ReadAll(r.Body)

			if err != nil {
				fmt.Println(err)
			} else {
				stringBody := fmt.Sprintf("%s", body)
				fmt.Println(utf8.RuneCountInString(stringBody))
				//fmt.Println(stringBody)		
			}
			//fmt.Println(reflect.TypeOf(r.Body))
			//fmt.Printf("%#v", r.Body)
			//io.Copy(os.Stdout, r.Body)
	//		fmt.Println(r.Body)
	//	    fmt.Println("hello world")	
			fmt.Println("Success")
			fmt.Println(url)
		}

		//res := new(Result)
		finishChan <- true
	}()
}

func main() {

	urls := []string{"https://www.tgstatic.com/lt", "https://www.tgstatic.com/en"}
	resultChan, finishChan := checkUrls(urls)
	go func() {
		result := <-resultChan
		fmt.Println(result)
	}()
	<-finishChan
}
