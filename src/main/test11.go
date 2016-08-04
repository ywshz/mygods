package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"net/url"
	"time"
)

func main() {
	start := time.Now()
	defer func(s time.Time) {
		fmt.Printf("\nCost %v seconds.\n", time.Since(s))
	}(start)

	resp, err := http.PostForm("http://localhost:8888/runjs",
		url.Values{"script": {"newTime=day+' 12:12:12'"}, "params":{"{\"day\":\"2016-07-07\"}"}})

	//resp, err := http.Post("http://localhost:8888/runjs",
	//	"application/x-www-form-urlencoded",
	//	strings.NewReader("script=cjb"))
	if err != nil {
		fmt.Println(err)
	}
	//
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}
	//
	fmt.Println(string(body))

}
