package swiss

import (
	"net/http"
	"fmt"
	"html"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

type Worker struct {
	zk      *zk.Conn
	stopChn chan struct{}
}

func NewWorker() *Worker {
	conn, _, err := zk.Connect([]string{"127.0.0.1"}, 10*time.Second)
	if err != nil {
		panic(err)
		log.Error("Zk连接创建失败...", err)
	}

	log.Info("Zk连接创建成功...")

	return &Worker{zk:conn, stopChn: make(chan struct{})}
}

func (w *Worker) Start() {
	w.startHttpServer()
	w.registeToZk()

	<-w.stopChn
}

func (w *Worker) startHttpServer() {
	log.Info("Set up http server for receive job request on port:", "8080")
	http.HandleFunc("/jobreceiver", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	go http.ListenAndServe(":8080", nil)
}

func (w *Worker) registeToZk() {
	log.Info("准备注册到zk")
	if exist, _, _ := w.zk.Exists("/swiss/workers"); !exist {
		log.Info("/swiss/workers不存在, 创建中")
		log.Info(w.zk.Create("/swiss/workers", []byte(""), 0, WorldACLPermAll))
	}else {
		log.Info("/swiss/workers存在, 略过")
	}
	log.Info("注册Worker信息")
	log.Info(w.zk.Create("/swiss/workers/w1", []byte("127.0.0.1:8080"), zk.FlagEphemeral, WorldACLPermAll))
}
