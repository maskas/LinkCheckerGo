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
    Field1 string
    Field2 int
}

func checkUrl(url string) <-chan bool {
	c := make(chan bool)

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
		}
		c <- true
	}()
	return c
}

func main() {
	//time.Sleep(3000 * time.Millisecond)
	c := checkUrl("https://www.transfergo.com/lt")
	a := <-c
	fmt.Println(a)

	//go checkUrl("https://www.transfergo.com/lt")
	//go checkUrl("https://www.transfergo.com/lt")
	//go checkUrl("https://www.transfergo.com/lt")
	fmt.Println("End")
	//time.Sleep(3000 * time.Millisecond)
	
}
