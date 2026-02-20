package service

import (
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

func (s *ConversationService) GetConversation(id string, page, size int, jwt string) (*models.ConversationPage, error) {
	log.Printf("Retrieving conversation for id[%s]", id)
	return s.convClient.GetConversation(id, size, (page-1)*size, jwt)
}

func (s *ConversationService) GetUserConversations(page, size int, jwt string) (*models.ConversationSummaryPage, error) {
	log.Printf("Retrieving user conversations")
	return s.convClient.ListUserConversations(size, (page-1)*size, jwt)
}

func (s *ConversationService) AssociateUserConversation(id string, jwt string) error {
	log.Printf("Associating conversation for id[%s]", id)
	return s.convClient.AssociateConversation(id, jwt)
}

func (s *ConversationService) GetSuggestions(id string, jwt string) ([]models.Suggestion, error) {
	log.Printf("Retrieving suggestions for id[%s]", id)

	// Fetch preferences to improve suggestions
	prefs, err := s.profileClient.GetUserPreferences(jwt)
	var preferences map[string][]string
	if err == nil {
		preferences = prefs.Chat
	} else {
		log.Printf("Warning: failed to fetch user preferences for suggestions: %v", err)
	}

	return s.convClient.GenerateSuggestions(id, jwt, preferences)
}

func (s *ConversationService) DeleteConversation(id string, jwt string) error {
	log.Printf("Deleting conversation for id[%s]", id)
	// Java implementation was returning empty Mono, but let's assume it should call something if available
	return nil
}

func (s *ConversationService) ConversationInteract(request models.ConversationRequest, jwt string) (*models.ConversationResponse, error) {
	log.Printf("Handling conversation interaction for id[%s]", request.ConversationID)

	// Fetch preferences to enrich the LLM context
	prefs, err := s.profileClient.GetUserPreferences(jwt)
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

	resp, err := s.saClient.ConversationInteract(request, jwt, preferences)
	if err != nil {
		return nil, err
	}

	return &models.ConversationResponse{
		ConversationID: resp.SessionID,
		Message:        resp.Message,
	}, nil
}
