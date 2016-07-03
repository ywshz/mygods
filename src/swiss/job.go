package swiss

import (
	"time"
	"fmt"
	"encoding/json"
	"strconv"
)

type Job struct {
	Id           int `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	//Group        string `json:"group"`
	ScheduleType int `json:"scheduleType"`
	Cron         string `json:"cron"`
	Dependency   []int `json:"dependency"`
	ScriptType   int `json:"scriptType"`
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
	Tags  []string `json:tags`
	Ip    []string `json:ip`
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
			Tags : job.WorkerTag,
			Ip: job.WorkerIp,
		})

		path, error := job.Server.Store.Create("/swiss/readytorun/" + strconv.Itoa(executionId), exeJobInfo)
		if error != nil {
			log.Error(error)
		}else {
			log.Info("Job save to ", path)
		}
	}
}

// Friendly format a job
func (j *Job) String() string {
	return fmt.Sprintf("\"Job: %s, scheduled at: %s, tags:%v\"", j.Name, j.Cron, j.WorkerTag)
}
