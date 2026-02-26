package service

import (
	"context"
	"log"

	"julia-conversation-api/internal/api/models"
	"julia-conversation-api/internal/client"
)

type ConversationService struct {
	convClient    *client.ConversationClient
	saClient      *client.SuperAgentClient
	profileClient *client.ProfileClient
}

func NewConversationService(convClient *client.ConversationClient, saClient *client.SuperAgentClient, profileClient *client.ProfileClient) *ConversationService {
	return &ConversationService{
		convClient:    convClient,
		saClient:      saClient,
		profileClient: profileClient,
	}
}

func (s *ConversationService) GetConversation(ctx context.Context, id string, page, size int) (*models.ConversationPage, error) {
	log.Printf("Retrieving conversation for id[%s]", id)
	return s.convClient.GetConversation(ctx, id, size, (page-1)*size)
}

func (s *ConversationService) GetUserConversations(ctx context.Context, page, size int) (*models.ConversationSummaryPage, error) {
	log.Printf("Retrieving user conversations")
	return s.convClient.ListUserConversations(ctx, size, (page-1)*size)
}

func (s *ConversationService) AssociateUserConversation(ctx context.Context, id string) error {
	log.Printf("Associating conversation for id[%s]", id)
	return s.convClient.AssociateConversation(ctx, id)
}

func (s *ConversationService) GetSuggestions(ctx context.Context, id string) ([]models.Suggestion, error) {
	log.Printf("Retrieving suggestions for id[%s]", id)

	// Fetch preferences to improve suggestions
	prefs, err := s.profileClient.GetUserPreferences(ctx)
	var preferences map[string][]string
	if err == nil {
		preferences = prefs.Chat
	} else {
		log.Printf("Warning: failed to fetch user preferences for suggestions: %v", err)
	}

	return s.convClient.GenerateSuggestions(ctx, id, preferences)
}

func (s *ConversationService) DeleteConversation(ctx context.Context, id string) error {
	log.Printf("Deleting conversation for id[%s]", id)
	return s.convClient.DeleteConversation(ctx, id)
}

func (s *ConversationService) ConversationInteract(ctx context.Context, request models.ConversationRequest) (*models.ConversationResponse, error) {
	log.Printf("Handling conversation interaction for id[%s]", request.ConversationID)

	// Fetch preferences to enrich the LLM context
	prefs, err := s.profileClient.GetUserPreferences(ctx)
	var preferences map[string]interface{}
	if err == nil {
		preferences = make(map[string]interface{})
		for k, v := range prefs.Chat {
			preferences[k] = v
		}
		preferences["notifications"] = prefs.Notifications
	} else {
		log.Printf("Warning: failed to fetch user preferences for interaction: %v", err)
	}

	resp, err := s.saClient.ConversationInteract(ctx, request, preferences)
	if err != nil {
		return nil, err
	}

	return &models.ConversationResponse{
		ConversationID: resp.SessionID,
		Message:        resp.Message,
	}, nil
}
