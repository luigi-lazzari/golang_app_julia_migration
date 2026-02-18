package service

import (
	"fmt"
	"log"
)

type ExternalGateway interface {
	GetNotificationPreferences() ([]NewsExternalItem, error)
}

type InternalGateway interface {
	UpdateNotificationNewsPreferences(news []NewsItem) error
}

type OrchestratorService struct {
	externalGateway ExternalGateway
	internalGateway InternalGateway
}

func NewOrchestratorService(ext ExternalGateway, int InternalGateway) *OrchestratorService {
	return &OrchestratorService{
		externalGateway: ext,
		internalGateway: int,
	}
}

func (s *OrchestratorService) OrchestrateNotificationPreferencesUpdate() error {
	log.Println("Starting orchestration of notification preferences update")

	externalNews, err := s.externalGateway.GetNotificationPreferences()
	if err != nil {
		return fmt.Errorf("failed to fetch external news: %w", err)
	}

	newsItems := MapExternalToInternal(externalNews)

	err = s.internalGateway.UpdateNotificationNewsPreferences(newsItems)
	if err != nil {
		return fmt.Errorf("failed to update notification news: %w", err)
	}

	log.Println("Completed orchestration of notification preferences update")
	return nil
}

func MapExternalToInternal(external []NewsExternalItem) []NewsItem {
	internal := make([]NewsItem, len(external))
	for i, e := range external {
		internal[i] = NewsItem{
			ID:          e.ID,
			Description: e.Description,
			Channel:     e.Channel,
		}
	}
	return internal
}
