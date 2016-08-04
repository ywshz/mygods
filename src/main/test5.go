package main

import (
	"fmt"
	"github.com/bitly/go-simplejson"
)

func main(){
	msg:="{\"script\":\"console.log(name)\",\"params\":{\"name\":\"[123,456]\"}}";

	js, err := simplejson.NewJson([]byte(msg))
	if err!=nil {

		fmt.Println(err)
	}
	fmt.Println(js.Get("script").MustString())
	fmt.Println(js.Get("params").MustMap()["name"])
}