package swiss

import (
	"github.com/robfig/cron"
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
		log.Info("scheduler: Adding job[%s] to cron",job.Name)

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