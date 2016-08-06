package jse

import (
	"github.com/robertkrimen/otto"
	"fmt"
	"time"
	"encoding/json"
	"net/http"
	"net/url"
	"io/ioutil"
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
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	fmt.Printf("Run script:\n %s \nwith:\n %v \n\n", script, params)
	start := time.Now()

	vm := otto.New()

	//增加post方法
	vm.Set("$post", func(call otto.FunctionCall) otto.Value {
		result, _ := vm.ToValue(post(call.Argument(0).String(), call.Argument(1).String()))
		return result
	})

	vm.Set("$get", func(call otto.FunctionCall) otto.Value {
		result, _ := vm.ToValue(get(call.Argument(0).String(), call.Argument(1).String()))
		return result
	})

	vm.Run(`
		var $getForJSON = function(url,params){
			return JSON.parse($get(url,params));
		}
		var $postForJSON = function(url,params){
			return JSON.parse($post(url,params));
		}
	`)

	//vm.Set("getForArray", func(call otto.FunctionCall) otto.Value {
	//	result, _ := vm.ToValue(getForArray(call.Argument(0).String(), call.Argument(1).String()))
	//	return result
	//})
	//
	//vm.Set("postForMap", func(call otto.FunctionCall) otto.Value {
	//	result, _ := vm.ToValue(postForMap(call.Argument(0).String(), call.Argument(1).String()))
	//	return result
	//})
	//
	//vm.Set("postForArray", func(call otto.FunctionCall) otto.Value {
	//	result, _ := vm.ToValue(postForArray(call.Argument(0).String(), call.Argument(1).String()))
	//	return result
	//})

	for key, val := range params {
		vm.Set(key, val)
		fmt.Printf("\nSet key:%v value:%v, type:%T \n", key, val, val)
	}

	v, err := vm.Run(script)

	if err != nil {
		errStr := fmt.Sprintf("Run JS Exception : %s", err)
		fmt.Println("Run JS Exception : ", err)
		return errStr,err
	}

	fmt.Printf("\nCost %v seconds.\n", (time.Now().Sub(start)))

	rs, err := v.Export()
	//rs, err := vm.Get("$result")
	if err != nil {
		return "", err
	}
	jsonString, _ := json.Marshal(rs)
	fmt.Println(string(jsonString))
	return string(jsonString), nil
}

func post(postUrl, params string) string {
	resp, err := http.PostForm(postUrl, url.Values{"params": {params}})

	if err != nil {
		fmt.Printf("error:", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return (string)(body)
}

func get(postUrl, params string) string  {
	resp, err := http.Get(postUrl)
	if err != nil {
		fmt.Printf("error:%v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body)
}

//func getForArray(postUrl, params string) []interface{}  {
//	resp, err := http.Get(postUrl)
//
//	if err != nil {
//		fmt.Printf("error:%v", err)
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	js,_ := simplejson.NewJson(body)
//	return js.MustArray()
//}
//
//func postForArray(postUrl, params string) []interface{}  {
//	resp, err := http.PostForm(postUrl, url.Values{"params": {params}})
//
//	if err != nil {
//		fmt.Printf("error:%v", err)
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	js,_ := simplejson.NewJson(body)
//	return js.MustArray()
//}
//
//func getForMap(postUrl, params string) map[string]interface{} {
//	resp, err := http.Get(postUrl)
//
//	if err != nil {
//		fmt.Printf("error:%v", err)
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	js,_ := simplejson.NewJson(body)
//	return js.MustMap()
//}
//
//func postForMap(postUrl, params string) map[string]interface{} {
//	resp, err := http.PostForm(postUrl, url.Values{"params": {params}})
//
//	if err != nil {
//		fmt.Printf("error:%v", err)
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	js,_ := simplejson.NewJson(body)
//	return js.MustMap()
//}