package azure

import (
	"crypto/tls"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/comune-roma/bff-julia-profile-api/internal/config"
)

// NewCosmosClient creates a new Cosmos DB client
func NewCosmosClient(cfg *config.Config) (*azcosmos.Client, error) {
	cred, err := azcosmos.NewKeyCredential(cfg.CosmosDB.Key)
	if err != nil {
		return nil, err
	}

	clientOptions := &azcosmos.ClientOptions{}

	// Disable SSL verification for emulator
	if cfg.CosmosDB.Emulator {
		clientOptions.ClientOptions.Transport = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	}

	client, err := azcosmos.NewClientWithKey(cfg.CosmosDB.Endpoint, cred, clientOptions)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// NewAppConfigClient creates a new Azure App Configuration client
func NewAppConfigClient(cfg *config.Config) (interface{}, error) {
	// TODO: Implement Azure App Configuration client
	// For now, return nil as placeholder
	return nil, nil
}
