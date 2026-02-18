package scheduler

import (
	"log"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron *cron.Cron
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		cron: cron.New(),
	}
}

func (s *Scheduler) AddJob(spec string, cmd cron.Job) (cron.EntryID, error) {
	return s.cron.AddJob(spec, cmd)
}

func (s *Scheduler) Start() {
	log.Println("Starting scheduler...")
	s.cron.Start()
}

func (s *Scheduler) Stop() {
	log.Println("Stopping scheduler...")
	s.cron.Stop()
}
