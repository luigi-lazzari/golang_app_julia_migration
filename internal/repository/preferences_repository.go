package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

// UserPreferencesRepository handles Cosmos DB operations for user preferences
type UserPreferencesRepository struct {
	client    *azcosmos.Client
	database  string
	container string
}

// NewUserPreferencesRepository creates a new UserPreferencesRepository
func NewUserPreferencesRepository(client *azcosmos.Client, database string) *UserPreferencesRepository {
	return &UserPreferencesRepository{
		client:    client,
		database:  database,
		container: "user_preferences",
	}
}

// GetPreferences retrieves user preferences from Cosmos DB
func (r *UserPreferencesRepository) GetPreferences(ctx context.Context, userID string) ([]byte, error) {
	containerClient, err := r.client.NewContainer(r.database, r.container)
	if err != nil {
		return nil, fmt.Errorf("failed to get container client: %w", err)
	}

	pk := azcosmos.NewPartitionKeyString(userID)
	resp, err := containerClient.ReadItem(ctx, pk, userID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to read preferences: %w", err)
	}

	return resp.Value, nil
}

// CreatePreferences creates new user preferences in Cosmos DB
func (r *UserPreferencesRepository) CreatePreferences(ctx context.Context, userID string, preferences interface{}) error {
	containerClient, err := r.client.NewContainer(r.database, r.container)
	if err != nil {
		return fmt.Errorf("failed to get container client: %w", err)
	}

	marshalledItem, err := json.Marshal(preferences)
	if err != nil {
		return fmt.Errorf("failed to marshal preferences: %w", err)
	}

	pk := azcosmos.NewPartitionKeyString(userID)
	_, err = containerClient.CreateItem(ctx, pk, marshalledItem, nil)
	if err != nil {
		return fmt.Errorf("failed to create preferences: %w", err)
	}

	return nil
}

// UpdatePreferences updates existing user preferences in Cosmos DB
func (r *UserPreferencesRepository) UpdatePreferences(ctx context.Context, userID string, preferences interface{}) error {
	containerClient, err := r.client.NewContainer(r.database, r.container)
	if err != nil {
		return fmt.Errorf("failed to get container client: %w", err)
	}

	marshalledItem, err := json.Marshal(preferences)
	if err != nil {
		return fmt.Errorf("failed to marshal preferences: %w", err)
	}

	pk := azcosmos.NewPartitionKeyString(userID)
	_, err = containerClient.ReplaceItem(ctx, pk, userID, marshalledItem, nil)
	if err != nil {
		return fmt.Errorf("failed to replace preferences: %w", err)
	}

	return nil
}

// DeletePreferences deletes user preferences from Cosmos DB
func (r *UserPreferencesRepository) DeletePreferences(ctx context.Context, userID string) error {
	containerClient, err := r.client.NewContainer(r.database, r.container)
	if err != nil {
		return fmt.Errorf("failed to get container client: %w", err)
	}

	pk := azcosmos.NewPartitionKeyString(userID)
	_, err = containerClient.DeleteItem(ctx, pk, userID, nil)
	if err != nil {
		return fmt.Errorf("failed to delete preferences: %w", err)
	}

	return nil
}
