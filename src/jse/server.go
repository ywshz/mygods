package jse

import (
	"net/http"
	"encoding/json"
)

type Server struct {
	jsEngine *JsEngine
}

func NewServer() *Server {
	return &Server{
		jsEngine: NewJsEngine(),
	}
}
func (j *Server) Start(port string) {
	http.HandleFunc("/runjs", func(w http.ResponseWriter, req *http.Request) {
		script := req.PostFormValue("script")
		params := req.PostFormValue("params")
		var paramsMap map[string]interface{}
		json.Unmarshal([]byte(params), &paramsMap)

		value, err := j.jsEngine.Run(script, paramsMap)

		if err != nil {
			panic(err)
		}

		w.Write([]byte(value))
	})

	http.ListenAndServe(":" + port, nil)
}


