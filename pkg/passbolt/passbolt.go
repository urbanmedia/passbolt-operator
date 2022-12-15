package passbolt

import (
	"context"
	"fmt"
	"net/http"

	"github.com/passbolt/go-passbolt/api"
)

type Client struct {
	passboltClient *api.Client
}

// NewClient initializes a new passbolt client and logs in.
// The client is configured to use the given URL, username and password.
func NewClient(ctx context.Context, url, username, password string) (*Client, error) {
	clnt, err := api.NewClient(&http.Client{}, "", "https://passbolt.example.com", "", "")
	if err != nil {
		return nil, fmt.Errorf("failed to create passbolt client: %w", err)
	}
	err = clnt.Login(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to login to passbolt: %w", err)
	}
	return &Client{
		passboltClient: clnt,
	}, nil
}

// Close logs out of the passbolt client.
// This should be called when the client is no longer needed.
func (c *Client) Close(ctx context.Context) error {
	return c.passboltClient.Logout(ctx)
}

// GetSecret retrieves the secret value for the given secret ID.
// The secret value is returned as a string.
func (c *Client) GetSecret(ctx context.Context, id string) (string, error) {
	scrt, err := c.passboltClient.GetSecret(ctx, id)
	if err != nil {
		return "", fmt.Errorf("failed to get secret: %s", err)
	}
	return scrt.Data, nil
}
