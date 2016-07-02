package swiss

import (
	"github.com/robfig/cron"
	"github.com/Sirupsen/logrus"
)

type Scheduler struct {
	Cron    *cron.Cron
	Started bool
}

func NewScheduler() *Scheduler {
	return &Scheduler{Cron: cron.New(), Started: false}
}

func (s *Scheduler) Start(jobs []*Job) {
	for _, job := range jobs {
		log.WithFields(logrus.Fields{
			"job": job.Name,
		}).Debug("scheduler: Adding job to cron")

		s.Cron.AddJob(job.Cron, job)
	}
	s.Cron.Start()
	s.Started = true
}

func (s *Scheduler) Restart(jobs []*Job) {
	s.Cron.Stop()
	s.Cron = cron.New()
	s.Start(jobs)
}

func (s *Scheduler) Stop() {
	s.Cron.Stop()
	s.Started = false
}