package passbolt

import (
	"context"
	"fmt"
	"net/http"

	"github.com/passbolt/go-passbolt/api"
	"github.com/passbolt/go-passbolt/helper"
)

type Client struct {
	passboltClient *api.Client

	// secretCache represents a cache of NAME -> UUID mappings.
	// This is used to avoid unnecessary API calls.
	secretCache map[string]string
}

// NewClient initializes a new passbolt client and logs in.
// The client is configured to use the given URL, username and password.
func NewClient(ctx context.Context, url, username, password string) (*Client, error) {
	clnt, err := api.NewClient(&http.Client{}, "", url, username, password)
	if err != nil {
		return nil, fmt.Errorf("failed to create passbolt client: %w", err)
	}
	err = clnt.Login(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to login to passbolt: %w", err)
	}
	return &Client{
		passboltClient: clnt,
		secretCache:    map[string]string{},
	}, nil
}

// LoadCache fills the secret cache with all secret names and IDs.
// This should be called before any secrets are retrieved.
// This is necessary because the passbolt API does not allow for searching secrets by name.
// Instead, we must retrieve all secrets and their UUIDs.
// This is not ideal, but it is the only way to retrieve secrets by name.
func (c *Client) LoadCache(ctx context.Context) error {
	resources, err := c.passboltClient.GetResources(ctx, &api.GetResourcesOptions{})
	if err != nil {
		return fmt.Errorf("failed to get secrets: %w", err)
	}
	for _, sctr := range resources {
		c.secretCache[sctr.Name] = sctr.ID
	}
	return nil
}

// Close logs out of the passbolt client.
// This should be called when the client is no longer needed.
func (c *Client) Close(ctx context.Context) error {
	return c.passboltClient.Logout(ctx)
}

// GetSecret retrieves the secret value for the given secret ID.
// The secret value is returned as a string.
func (c *Client) GetSecret(ctx context.Context, name string) (string, error) {
	if _, ok := c.secretCache[name]; !ok {
		return "", fmt.Errorf("unable to find secret in cache with name %q", name)
	}
	// retrieve the secret
	_, _, _, _, pw, _, err := helper.GetResource(ctx, c.passboltClient, c.secretCache[name])
	if err != nil {
		return "", fmt.Errorf("failed to get secret with name %q: %w", name, err)
	}
	return pw, nil
}
