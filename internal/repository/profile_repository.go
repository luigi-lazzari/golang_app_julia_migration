package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

// UserProfileRepository handles Cosmos DB operations for user profiles
type UserProfileRepository struct {
	client    *azcosmos.Client
	database  string
	container string
}

// NewUserProfileRepository creates a new UserProfileRepository
func NewUserProfileRepository(client *azcosmos.Client, database string) *UserProfileRepository {
	return &UserProfileRepository{
		client:    client,
		database:  database,
		container: "user_profiles",
	}
}

// GetProfile retrieves a user profile from Cosmos DB
func (r *UserProfileRepository) GetProfile(ctx context.Context, userID string) ([]byte, error) {
	containerClient, err := r.client.NewContainer(r.database, r.container)
	if err != nil {
		return nil, fmt.Errorf("failed to get container client: %w", err)
	}

	pk := azcosmos.NewPartitionKeyString(userID)
	resp, err := containerClient.ReadItem(ctx, pk, userID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to read profile: %w", err)
	}

	return resp.Value, nil
}

// CreateProfile creates a new user profile in Cosmos DB
func (r *UserProfileRepository) CreateProfile(ctx context.Context, userID string, profile interface{}) error {
	containerClient, err := r.client.NewContainer(r.database, r.container)
	if err != nil {
		return fmt.Errorf("failed to get container client: %w", err)
	}

	marshalledItem, err := json.Marshal(profile)
	if err != nil {
		return fmt.Errorf("failed to marshal profile: %w", err)
	}

	pk := azcosmos.NewPartitionKeyString(userID)
	_, err = containerClient.CreateItem(ctx, pk, marshalledItem, nil)
	if err != nil {
		return fmt.Errorf("failed to create profile: %w", err)
	}

	return nil
}

// UpdateProfile updates an existing user profile in Cosmos DB
func (r *UserProfileRepository) UpdateProfile(ctx context.Context, userID string, profile interface{}) error {
	containerClient, err := r.client.NewContainer(r.database, r.container)
	if err != nil {
		return fmt.Errorf("failed to get container client: %w", err)
	}

	marshalledItem, err := json.Marshal(profile)
	if err != nil {
		return fmt.Errorf("failed to marshal profile: %w", err)
	}

	pk := azcosmos.NewPartitionKeyString(userID)
	_, err = containerClient.ReplaceItem(ctx, pk, userID, marshalledItem, nil)
	if err != nil {
		return fmt.Errorf("failed to replace profile: %w", err)
	}

	return nil
}

// DeleteProfile deletes a user profile from Cosmos DB
func (r *UserProfileRepository) DeleteProfile(ctx context.Context, userID string) error {
	containerClient, err := r.client.NewContainer(r.database, r.container)
	if err != nil {
		return fmt.Errorf("failed to get container client: %w", err)
	}

	pk := azcosmos.NewPartitionKeyString(userID)
	_, err = containerClient.DeleteItem(ctx, pk, userID, nil)
	if err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}

	return nil
}

// Ping checks the connection to Cosmos DB
func (r *UserProfileRepository) Ping(ctx context.Context) error {
	database, err := r.client.NewDatabase(r.database)
	if err != nil {
		return err
	}
	_, err = database.Read(ctx, nil)
	return err
}
