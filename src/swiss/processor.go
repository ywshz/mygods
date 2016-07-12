package swiss

import "time"

type Processor struct {
	Job ReadyToRunJob
}

func (p *Processor) Run(){
	log.Info("run the job by worker: ", p.Job)
	//get job script
	var script string
	//设置超时提醒
	slowTimer := time.AfterFunc(2*time.Hour, func() {
		log.Warnf("proc: Script '%s' slow, execution exceeding %v", script, 2*time.Hour)
	})
	//运行job

	//结束超时计时器
	slowTimer.Stop()
}
