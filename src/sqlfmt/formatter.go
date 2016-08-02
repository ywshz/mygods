package main

import (
	"io/ioutil"
	"fmt"
)

func main() {
	dat, _ := ioutil.ReadFile("/Users/yws/Documents/idea_workspace/mygods/src/sqlfmt/demo.sql")
	sql := string(dat)

	for _, value := range sql {
		if value != ' ' {
			fmt.Printf("%c", value)
		}
	}
	//sqlarr := strings.Split(string(dat), " ")
	//
	//for _, value := range sqlarr {
	//	fmt.Println(value)
	//}
}


