package service

import (
	"log"
	"sync"
	"time"
)

const deduplicationTTL = 24 * time.Hour

type DeduplicationService struct {
	processedMessages map[string]time.Time
	mu                sync.RWMutex
}

func NewDeduplicationService() *DeduplicationService {
	s := &DeduplicationService{
		processedMessages: make(map[string]time.Time),
	}
	go s.startCleanupTask()
	return s
}

// IsDuplicate checks if a messageId has already been processed within the TTL.
func (s *DeduplicationService) IsDuplicate(messageID string) bool {
	if messageID == "" {
		return false
	}

	s.mu.RLock()
	processedAt, found := s.processedMessages[messageID]
	s.mu.RUnlock()

	if !found {
		return false
	}

	// Double check TTL
	if time.Since(processedAt) < deduplicationTTL {
		log.Printf("Duplicate message detected: MessageId=%s, originallyProcessedAt=%v", messageID, processedAt)
		return true
	}

	// Expired, should be removed by cleanup but we can ignore it here
	return false
}

// MarkAsProcessed marks a messageId as successfully processed.
func (s *DeduplicationService) MarkAsProcessed(messageID string) {
	if messageID == "" {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.processedMessages[messageID] = time.Now()
	log.Printf("Message marked as processed in deduplication cache: MessageId=%s", messageID)
}

func (s *DeduplicationService) startCleanupTask() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C
		s.cleanup()
	}
}

func (s *DeduplicationService) cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	cutoff := time.Now().Add(-deduplicationTTL)
	removed := 0

	for id, processedAt := range s.processedMessages {
		if processedAt.Before(cutoff) {
			delete(s.processedMessages, id)
			removed++
		}
	}

	if removed > 0 {
		log.Printf("Cleaned up %d expired deduplication entries (older than %v)", removed, deduplicationTTL)
	}
}

func (s *DeduplicationService) GetCacheSize() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.processedMessages)
}
