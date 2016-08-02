package main

import (
	"net/http"
	"net/url"
	"fmt"
	"io/ioutil"
)

func main() {
	resp, err := http.PostForm("http://localhost:8080/jobreceiver",
		url.Values{"job": {"{123}"}})

	if err != nil {
		fmt.Printf("error:", err)
	}else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("res:", string(body))
	}
	defer resp.Body.Close()
}
