package jobs

import (
	"log"
	"time"
)

type Orchestrator interface {
	OrchestrateNotificationPreferencesUpdate() error
}

type NotificationJob struct {
	orchestrator Orchestrator
	maxRetries   int
}

func NewNotificationJob(orchestrator Orchestrator, maxRetries int) *NotificationJob {
	return &NotificationJob{
		orchestrator: orchestrator,
		maxRetries:   maxRetries,
	}
}

func (j *NotificationJob) Run() {
	log.Println("*************** Starting NotificationJob execution *****************")
	maxRetries := j.maxRetries
	retryCount := 0

	for {
		log.Printf("Max Retries Property: %d", maxRetries)
		log.Printf("Retry Count Property: %d", retryCount)
		log.Printf("*************** Actual Retry Count Property: %d ********************", retryCount)

		err := j.orchestrator.OrchestrateNotificationPreferencesUpdate()
		if err == nil {
			log.Println("NotificationJob completed successfully")
			return
		}

		if retryCount < maxRetries {
			retryCount++
			log.Printf("Error occurred: %v. Trying to refire job execution", err)
			// Small delay before refire (Go doesn't have native "refire" as Quartz does, so we loop)
			time.Sleep(1 * time.Second)
			continue
		} else {
			log.Printf("Max retries reached for NotificationJob. No more attempts will be made. Error: %v", err)
			return
		}
	}
}
