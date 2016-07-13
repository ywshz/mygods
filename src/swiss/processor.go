package swiss

import (
	"time"
	"github.com/samuel/go-zookeeper/zk"
	"os"
	"runtime"
	"os/exec"
	"github.com/armon/circbuf"
	"strconv"
	"github.com/Sirupsen/logrus"
)

type Processor struct {
	Job ReadyToRunJob
	zk  *zk.Conn
}

func (p *Processor) Run() {
	log.WithFields(logrus.Fields{
		"Job": p.Job,
	}).Debugln("Run job by worker")

	//get job script
	scriptBts, _, _ := p.zk.Get("/swiss/jobscript/" + strconv.Itoa(p.Job.JobId))
	script := string(scriptBts)
	log.WithFields(logrus.Fields{
		"Script": script,
	}).Debugln("Get the job script")

	fullFilePath := "/Users/yws/Documents/dev/works/" + strconv.Itoa(p.Job.ExeId)
	log.WithFields(logrus.Fields{
		"Path": fullFilePath,
	}).Debugln("Job script full path")

	//write to file
	fout, err := os.Create(fullFilePath)
	if err != nil {
		log.WithFields(logrus.Fields{
			"File": fullFilePath,
			"Error": err,
		}).Error("Job script full path")
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
		log.Warnf("proc: Script '%s' slow, execution exceeding %v", script, 2 * time.Hour)
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

	log.WithFields(logrus.Fields{
		"shell":shell,
		"flag":flag,
	}).Debug("Start the cmd")
	cmd.Start()
	go func() {
		err = cmd.Wait()
		slowTimer.Stop()

		if err != nil {
			log.WithError(err).Error("proc: command error output")
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
