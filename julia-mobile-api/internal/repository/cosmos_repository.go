package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

// CosmosRepository handles Cosmos DB operations
type CosmosRepository struct {
	client   *azcosmos.Client
	database string
}

// NewCosmosRepository creates a new CosmosRepository
func NewCosmosRepository(client *azcosmos.Client, database string) *CosmosRepository {
	return &CosmosRepository{
		client:   client,
		database: database,
	}
}

// GetItem retrieves an item from Cosmos DB
func (r *CosmosRepository) GetItem(ctx context.Context, container, id, partitionKey string) ([]byte, error) {
	containerClient, err := r.client.NewContainer(r.database, container)
	if err != nil {
		return nil, fmt.Errorf("failed to get container client: %w", err)
	}

	pk := azcosmos.NewPartitionKeyString(partitionKey)
	resp, err := containerClient.ReadItem(ctx, pk, id, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to read item: %w", err)
	}

	return resp.Value, nil
}

// CreateItem creates a new item in Cosmos DB
func (r *CosmosRepository) CreateItem(ctx context.Context, container string, item interface{}, partitionKey string) error {
	containerClient, err := r.client.NewContainer(r.database, container)
	if err != nil {
		return fmt.Errorf("failed to get container client: %w", err)
	}

	marshalledItem, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	pk := azcosmos.NewPartitionKeyString(partitionKey)
	_, err = containerClient.CreateItem(ctx, pk, marshalledItem, nil)
	if err != nil {
		return fmt.Errorf("failed to create item: %w", err)
	}

	return nil
}

// UpdateItem updates an existing item in Cosmos DB
func (r *CosmosRepository) UpdateItem(ctx context.Context, container, id string, item interface{}, partitionKey string) error {
	containerClient, err := r.client.NewContainer(r.database, container)
	if err != nil {
		return fmt.Errorf("failed to get container client: %w", err)
	}

	marshalledItem, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	pk := azcosmos.NewPartitionKeyString(partitionKey)
	_, err = containerClient.ReplaceItem(ctx, pk, id, marshalledItem, nil)
	if err != nil {
		return fmt.Errorf("failed to replace item: %w", err)
	}

	return nil
}

// DeleteItem deletes an item from Cosmos DB
func (r *CosmosRepository) DeleteItem(ctx context.Context, container, id, partitionKey string) error {
	containerClient, err := r.client.NewContainer(r.database, container)
	if err != nil {
		return fmt.Errorf("failed to get container client: %w", err)
	}

	pk := azcosmos.NewPartitionKeyString(partitionKey)
	_, err = containerClient.DeleteItem(ctx, pk, id, nil)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	return nil
}
