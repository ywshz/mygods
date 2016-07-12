package swiss

import (
	"net/http"
	"github.com/samuel/go-zookeeper/zk"
	"time"
	"encoding/json"
)

type Worker struct {
	name          string
	zk            *zk.Conn
	stopChn       chan struct{}
	maxProcessors chan struct{}
}

type WorkerConnectionInfo struct {
	Ip   string        `json:ip`
	Port string        `json:port`
}

func NewWorker() *Worker {
	conn, _, err := zk.Connect([]string{"127.0.0.1"}, 10 * time.Second)
	if err != nil {
		panic(err)
		log.Error("Zk连接创建失败...", err)
	}
	log.Info("Zk连接创建成功...")

	//最多10个任务同时运行
	return &Worker{
		zk:conn,
		stopChn: make(chan struct{}),
		maxProcessors:make(chan struct{}, 10),
	}
}

func (w *Worker) Start() {
	w.name = "w1"
	w.startHttpServer()
	w.registeToZk()

	<-w.stopChn
}

func (w *Worker) startHttpServer() {
	log.Info("Set up http server for receive job request on port:", "8080")
	http.HandleFunc("/jobreceiver", func(res http.ResponseWriter, req *http.Request) {
		job := ReadyToRunJob{}
		json.Unmarshal([]byte(req.PostFormValue("job")), &job)
		w.runJob(job)
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
	workConnInfo, _ := json.Marshal(WorkerConnectionInfo{
		Ip: "127.0.0.1",
		Port : "8080",
	})
	log.Info("####", string(workConnInfo))
	log.Info(w.zk.Create("/swiss/workers/" + w.name, workConnInfo, zk.FlagEphemeral, WorldACLPermAll))
}

func (w *Worker) runJob(job ReadyToRunJob) {
	w.maxProcessors <- struct{}{}
	defer func() {
		<-w.maxProcessors
	}()

	p := Processor{
		Job : job,
	}

	p.Run()
}

func (w *WorkerConnectionInfo) ToUrl() string {
	return "http://" + w.Ip + ":" + w.Port
}