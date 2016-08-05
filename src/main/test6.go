package main

import (
	"github.com/robertkrimen/otto"
	"fmt"
	"github.com/bitly/go-simplejson"
	"net/http"
	"io/ioutil"
	"net/url"
)

func main() {

	vm := otto.New()

	s, _ := simplejson.NewJson([]byte(`
	{"name":"[{\"name\":\"全部\",\"value\":\"all\"},{\"value\":372,\"name\":\"家居个护\"},{\"value\":437,\"name\":\"美容彩妆\"},{\"value\":438,\"name\":\"母婴\"},{\"value\":440,\"name\":\"数码家电\"},{\"value\":636,\"name\":\"【多余类目\"},{\"value\":836,\"name\":\"环球美食\"},{\"value\":837,\"name\":\"营养保健\"},{\"value\":1025,\"name\":\"箱包配饰\"},{\"value\":7578,\"name\":\"运动户外\"},{\"value\":8115,\"name\":\"生鲜\"},{\"value\":9691,\"name\":\"服装鞋靴\"}]"}
	`))
	//fmt.Println(s.Get("name").MustMap())
	vm.Set("def", s.Get("name").MustString())
	vm.Set("getArray", func(call otto.FunctionCall) otto.Value {
		result, _ := vm.ToValue(getArray(call.Argument(0).String(), call.Argument(1).String()))
		return result
	})
	vm.Set("getMap", func(call otto.FunctionCall) otto.Value {
		result, _ := vm.ToValue(getMap(call.Argument(0).String(), call.Argument(1).String()))
		return result
	})
	_, e := vm.Run(`
	//var x = [];
	//x=x.concat(eval(def))
	//	console.log(x[0].name)
		//console.log(x.push(1))
		//var x = {category_level:1};
    		//var x = JSON.parse(getData('http://localhost:8080/ds/50',JSON.stringify({name:1,age:19})));
    		//console.log(x.data)
    		//console.log(x.code)
    		console.log(getArray('http://www.jeasyui.com/demo/main/treegrid_data1.json')[0].id)
	`)
	if e != nil {
		fmt.Println(e)
	}
}

func getArray(postUrl, params string) []interface{}  {
	resp, err := http.PostForm(postUrl, url.Values{"params": {params}})

	if err != nil {
		fmt.Printf("error:%v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	js,_ := simplejson.NewJson(body)
	return js.MustArray()
}

func getMap(postUrl, params string) map[string]interface{} {
	resp, err := http.PostForm(postUrl, url.Values{"params": {params}})

	if err != nil {
		fmt.Printf("error:%v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	js,_ := simplejson.NewJson(body)
	return js.MustMap()
}