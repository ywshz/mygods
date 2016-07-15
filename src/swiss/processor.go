package swiss

import (
	"time"
	"github.com/samuel/go-zookeeper/zk"
	"os"
	"runtime"
	"os/exec"
	"github.com/armon/circbuf"
	"strconv"
)

type Processor struct {
	Job ReadyToRunJob
	zk  *zk.Conn
}

func (p *Processor) Run() {
	log.Debugf("Run job[%s] by worker", p.Job)

	//get job script
	scriptBts, _, _ := p.zk.Get("/swiss/jobscript/" + strconv.Itoa(p.Job.JobId))
	script := string(scriptBts)
	log.Debugf("Get the job script:%s", script)

	fullFilePath := "/Users/yws/Documents/dev/works/" + strconv.Itoa(p.Job.ExeId)
	log.Debugf("Job script full path:%s", fullFilePath)

	//write to file
	fout, err := os.Create(fullFilePath)
	if err != nil {
		log.Errorf("Job script full path,%s,%s",fullFilePath,err)
		return
	}
	fout.Write(scriptBts)
	fout.Close()
	defer func() {
		os.Remove(fullFilePath)
	}()
	//end write file

	//设置超时提醒
	slowTimer := time.AfterFunc(2 * time.Hour, func() {
		log.Warningf("proc: Script '%s' slow, execution exceeding %v", script, 2 * time.Hour)
	})

	//防止job的log太多,限制为256k, 当超过256k后,最早出现的日志将会逐步被代替
	output, _ := circbuf.NewBuffer(256000)

	//封装成cmd
	if p.Job.ScriptType == shell {

	}

	var shell, flag string

	switch p.Job.ScriptType{
	case shell :
		if runtime.GOOS == "windows" {
			shell = "cmd"
			flag = ""
		} else {
			shell = "/bin/sh"
			flag = ""
		}
	case hive :
		shell = "hive"
		flag = "-f"
	case python :
		shell = "python"
		flag = "-f"
	}

	cmd := exec.Command(shell, fullFilePath)
	cmd.Stderr = output
	cmd.Stdout = output
	//运行job
	var success, done bool

	log.Debugf("Start the cmd, Shell:%s, flag:%s", shell,flag)
	cmd.Start()
	go func() {
		err = cmd.Wait()
		slowTimer.Stop()

		if err != nil {
			log.Errorf("proc: command error output: %s",err)
			success = false
		} else {
			success = true
		}

		done = true
	}()

	for ; !done; {
		log.Info(output.String())
		time.Sleep(1 * time.Second)
	}

	log.Info("Job run completed, and success status is -> ", success)
	//execution.FinishedAt = time.Now()
	//execution.Success = success
	//execution.Output = output.Bytes()
}
