package jse

import (
	"net/http"
	"github.com/googollee/go-socket.io"
	"github.com/bitly/go-simplejson"
	"fmt"
	"encoding/json"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("swiss")

func init() {
	var format = logging.MustStringFormatter(
		`%{level:.4s} ▶ %{shortpkg}.%{shortfile}.%{longfunc} %{color:reset} %{message}`,
	)
	logging.SetFormatter(format)
	logging.SetLevel(logging.DEBUG, "jse")
}



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
		log.Info("on connection")
		so.On("runjs", func(msg string) interface{} {
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
			log.Info("on disconnect")
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Errorf("error:", err)
	})
	http.Handle("/socket.io/", server)
	log.Infof("Serving at port:%s", port)
	http.HandleFunc("/runjs", func(w http.ResponseWriter, req *http.Request) {
		script := req.PostFormValue("script")
		params := req.PostFormValue("params")

		fmt.Println(params)
		var paramsMap map[string]interface{}
		json.Unmarshal([]byte(params), &paramsMap)

		value, err := j.jsEngine.Run(script, paramsMap)

		res, _ := json.Marshal(value)

		w.Write((res))

		if err != nil {
			log.Error(err)
		}
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}


