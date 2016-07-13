package swiss

import (
	"time"
	"fmt"
	"encoding/json"
	"net/http"
	"net/url"
	"math/rand"
)

const (
	hive = "hive"
	python = "python"
	shell = "shell"
)
type Job struct {
	Id           int `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	//Group        string `json:"group"`
	ScheduleType int `json:"scheduleType"`
	Cron         string `json:"cron"`
	Dependency   []int `json:"dependency"`
	ScriptType   string `json:"scriptType"`
	WorkerTag    []string `json:"workerTag"`
	WorkerIp     []string `json:"workerIp"`
	CreateTime   time.Time `json:"createTime"`
	ModifyTime   time.Time `json:"modifyTime"`
	Status       int `json:"status"`
	Retries      uint `json:"retries"`
	SuccessCount int `json:"success_count"`
	ErrorCount   int `json:"error_count"`
	LastSuccess  time.Time `json:"last_success"`
	LastError    time.Time `json:"last_error"`
	Server       *Server `json:"-"`
}

type ReadyToRunJob struct {
	JobId int `json:jobId`
	ExeId int `json:ExeId`
	ScriptType string `json:scriptType`
	Tags  []string `json:tags`
	Ip    []string `json:ip`
	CreateTime time.Time `json:createTime`
}

func (job *Job) Run() {
	//如果是Leader, 往readytorun里塞记录
	if job.Server.candidate.IsLeader() {
		//取executionId
		executionId, _ := job.Server.Store.NextExecutionId()
		log.Info("Get job exec id -> ", executionId)
		exeJobInfo, _ := json.Marshal(ReadyToRunJob{
			JobId: job.Id,
			ExeId: executionId,
			ScriptType: job.ScriptType,
			Tags : job.WorkerTag,
			Ip: job.WorkerIp,
			CreateTime: time.Now(),
		})

		//get worker
		workers, err := job.Server.Store.GetWorkers(job.WorkerIp)
		var successSendToWorker bool = false

		if err != nil {
			successSendToWorker = false
			panic(err)
		}

		for i := len(workers); i > 0; i-- {
			worker := workers[rand.Intn(len(workers))]

			resp, err := http.PostForm(worker.ToUrl() + "/jobreceiver",
				url.Values{"job": {string(exeJobInfo)}})

			if err != nil {
				fmt.Printf("error:", err)
				continue
			}
			resp.Body.Close()

			successSendToWorker = true
			break;
		}

		if successSendToWorker {
			job.Server.Store.Create("/swiss/runningjobs/", exeJobInfo)
		} else {
			//TODO: notify send failed error
			log.Info("Job can not send to worker")
		}

	}
}

// Friendly format a job
func (j *Job) String() string {
	return fmt.Sprintf("\"Job: %s, scheduled at: %s, tags:%v\"", j.Name, j.Cron, j.WorkerTag)
}
