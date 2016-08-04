package jse

import (
	"net/http"
	//"encoding/json"
	"github.com/googollee/go-socket.io"
	"github.com/bitly/go-simplejson"
	"log"
	"fmt"
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
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")
		so.On("runjs", func(msg string) string {
			fmt.Println(msg)
			js, _ := simplejson.NewJson([]byte(msg))
			value, err := j.jsEngine.Run(js.Get("script").MustString(), js.Get("params").MustMap())
			if err != nil {
				panic(err)
			}
			fmt.Println("-->",value)
			return value
		})
		so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})
	http.Handle("/socket.io/", server)
	log.Println("Serving at localhost:%s", port)
	http.HandleFunc("/runjs", func(w http.ResponseWriter, req *http.Request) {
		script := req.PostFormValue("script")
		params := req.PostFormValue("params")

		fmt.Println(params)
		var paramsMap map[string]interface{}
		json.Unmarshal([]byte(params), &paramsMap)

		value, err := j.jsEngine.Run(script, paramsMap)

		if err != nil {
			panic(err)
		}

		w.Write([]byte(value))
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}


