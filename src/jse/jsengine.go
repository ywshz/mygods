package jse

import (
	"github.com/robertkrimen/otto"
	"fmt"
	"time"
	"encoding/json"
)

type JsEngine struct {
	//vm *otto.Otto
}

func NewJsEngine() *JsEngine {
	//vm := otto.New()

	return &JsEngine{
		//vm : vm,
	}
}

func (j *JsEngine) Run(script string, params map[string]interface{}) (string, error) {
	fmt.Printf("Run script:\n %s \nwith:\n %v \n\n", script, params)
	start := time.Now()

	vm := otto.New()

	for key, val := range params {
		vm.Set(key, val)
	}

	v, err := vm.Run(script)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("\nCost %v seconds.\n", (time.Now().Sub(start)).Seconds())

	rs, err := v.Export()
	//rs, err := vm.Get("$result")
	if err != nil {
		return "", err
	}
	jsonString, _ := json.Marshal(rs)
	fmt.Println(string(jsonString))
	return string(jsonString),nil
}



