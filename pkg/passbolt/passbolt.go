/*
Copyright 2022 @ Verlag Der Tagesspiegel GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package passbolt

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/passbolt/go-passbolt/api"
	"github.com/passbolt/go-passbolt/helper"
	passboltv1alpha2 "github.com/urbanmedia/passbolt-operator/api/v1alpha2"
)

type PassboltSecretDefinition struct {
	FolderParentID string
	Name           string
	Username       string
	URI            string
	Password       string
	Description    string
}

// FieldValue returns the value of the given field by name.
func (p PassboltSecretDefinition) FieldValue(fieldName passboltv1alpha2.FieldName) string {
	switch fieldName {
	case passboltv1alpha2.FieldNameUsername:
		return p.Username
	case passboltv1alpha2.FieldNameUri:
		return p.URI
	case passboltv1alpha2.FieldNamePassword:
		return p.Password
	default:
		return ""
	}
}

// Client is a passbolt client.
// It is used to retrieve secrets from passbolt.
// Internally, we cache the secret names and IDs to avoid unnecessary API calls.
// This is necessary because the passbolt API does not allow for searching secrets by name.
// Instead, we must retrieve all secrets and their UUIDs.
// This is not ideal, but it is the only way to retrieve secrets by name.
type Client struct {
	// passboltClient is the underlying passbolt client.
	passboltClient *api.Client
	// mu is used to prevent concurrent access to the secret cache.
	mu sync.RWMutex
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
		mu:             sync.RWMutex{},
	}, nil
}

// LoadCache fills the secret cache with all secret names and IDs.
// This should be called before any secrets are retrieved.
// This is necessary because the passbolt API does not allow for searching secrets by name.
// Instead, we must retrieve all secrets and their UUIDs.
// This is not ideal, but it is the only way to retrieve secrets by name.
func (c *Client) LoadCache(ctx context.Context) error {
	// prevent concurrent access to the cache
	c.mu.Lock()
	defer c.mu.Unlock()
	// retrieve all secrets
	resources, err := c.passboltClient.GetResources(ctx, &api.GetResourcesOptions{})
	if err != nil {
		return fmt.Errorf("failed to get secrets: %w", err)
	}
	// fill the cache
	for _, sctr := range resources {
		c.secretCache[sctr.Name] = sctr.ID
	}
	return nil
}

// Close logs out of the passbolt client.
// This should be called when the client is no longer needed.
func (c *Client) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.passboltClient.Logout(ctx)
}

// GetSecret retrieves the secret value for the given secret ID.
// The secret value is returned as a string.
func (c *Client) GetSecret(ctx context.Context, name string, fieldName passboltv1alpha2.FieldName) (*PassboltSecretDefinition, error) {
	// prevent concurrent access to the cache
	c.mu.Lock()
	defer c.mu.Unlock()
	// check if the secret is in the cache
	if _, ok := c.secretCache[name]; !ok {
		return nil, fmt.Errorf("unable to find secret in cache with name %q", name)
	}
	// retrieve the secret
	folderParentID, name, username, uri, pw, description, err := helper.GetResource(ctx, c.passboltClient, c.secretCache[name])
	if err != nil {
		return nil, fmt.Errorf("failed to get secret with name %q: %w", name, err)
	}
	secret := &PassboltSecretDefinition{
		Username:       username,
		URI:            uri,
		Password:       pw,
		Description:    description,
		FolderParentID: folderParentID,
		Name:           name,
	}
	return secret, nil
}
