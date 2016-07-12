package main

import (
	"encoding/json"
	"fmt"
)

type WorkerConnectionInfo struct {
	Ip   string
	Port string
}

func main() {
	workConnInfo, _ := json.Marshal(WorkerConnectionInfo{
		Ip: "127.0.0.1",
		Port : "8080",
	})

	fmt.Println(string(workConnInfo))

	var workConnInfo2 WorkerConnectionInfo
	json.Unmarshal(workConnInfo,&workConnInfo2)

	fmt.Println(workConnInfo2)
}
