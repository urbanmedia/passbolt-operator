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

	"github.com/passbolt/go-passbolt/api"
	"github.com/passbolt/go-passbolt/helper"
	"github.com/prometheus/client_golang/prometheus"
	passboltv1alpha2 "github.com/urbanmedia/passbolt-operator/api/v1alpha2"
	"github.com/urbanmedia/passbolt-operator/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	passboltSecretGetAttemptsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "passbolt_secret_get_attempts_total",
			Help: "Number of attempts to get a secret from passbolt.",
		},
	)
	passboltSecretGetFailureAttemptsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "passbolt_secret_get_failure_attempts_total",
			Help: "Number of failure attempts to get a secret from passbolt.",
		},
	)
	passboltReLogins = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "passbolt_relogins_total",
			Help: "Number of re-logins attempts to passbolt.",
		},
	)
	passboltReLoginFailures = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "passbolt_relogin_errors_total",
			Help: "Number of re-login error attempts to passbolt.",
		},
	)
	passboltCacheSync = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "passbolt_cache_sync_total",
			Help: "Number of cache syncs with passbolt.",
		},
	)
	passboltCacheFailures = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "passbolt_cache_sync_errors_total",
			Help: "Number of cache sync errors.",
		},
	)
)

func init() {
	// Register custom metrics with the global prometheus registry
	metrics.Registry.MustRegister(
		passboltSecretGetAttemptsTotal,
		passboltSecretGetFailureAttemptsTotal,
		passboltReLogins,
		passboltReLoginFailures,
		passboltCacheSync,
		passboltCacheFailures,
	)
}

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

	// cache is the secret cache.
	cache cache.Cacher
}

// NewClient initializes a new passbolt client and logs in.
// The client is configured to use the given URL, username and password.
func NewClient(ctx context.Context, cache cache.Cacher, url, username, password string) (*Client, error) {
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
		cache:          cache,
	}, nil
}

// LoadCache fills the secret cache with all secret names and IDs.
// This should be called before any secrets are retrieved.
// This is necessary because the passbolt API does not allow for searching secrets by name.
// Instead, we must retrieve all secrets and their UUIDs.
// This is not ideal, but it is the only way to retrieve secrets by name.
func (c *Client) LoadCache(ctx context.Context) error {
	// increase cache sync metric
	passboltCacheSync.Inc()

	// retrieve all secrets
	resources, err := c.passboltClient.GetResources(ctx, &api.GetResourcesOptions{})
	if err != nil {
		passboltCacheFailures.Inc()
		return fmt.Errorf("failed to get secrets: %w", err)
	}
	// fill the cache
	for _, sctr := range resources {
		c.cache.Set(ctx, sctr.Name, sctr.ID, 0)
	}
	return nil
}

// Close logs out of the passbolt client.
// This should be called when the client is no longer needed.
func (c *Client) Close(ctx context.Context) error {
	return c.passboltClient.Logout(ctx)
}

// GetSecret retrieves the secret value for the given secret name.
// Under the hook, this function queries the internal cache for the secret ID by name.
// If the secret is not in the cache, an error is returned.
// If the secret is in the cache, the secret is retrieved from passbolt.
func (c *Client) GetSecret(ctx context.Context, name string) (*PassboltSecretDefinition, error) {
	passboltSecretGetAttemptsTotal.Inc()

	// retrieve the secret ID from the cache
	val, err := c.cache.Get(ctx, name)
	if err != nil {
		passboltSecretGetFailureAttemptsTotal.Inc()
		return nil, err
	}

	if val == nil {
		passboltSecretGetFailureAttemptsTotal.Inc()
		return nil, fmt.Errorf("secret with name %q not found in cache", name)
	}

	// retrieve the secret
	folderParentID, name, username, uri, pw, description, err := helper.GetResource(ctx, c.passboltClient, val.(string))
	if err != nil {
		passboltSecretGetFailureAttemptsTotal.Inc()
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

// ReLogin logs out of the passbolt client and logs in again.
// This is useful if the session has expired.
// This function should be called before any other function.
func (c *Client) ReLogin(ctx context.Context) error {
	passboltReLogins.Inc()
	err := c.passboltClient.Login(ctx)
	if err != nil {
		passboltReLoginFailures.Inc()
		return fmt.Errorf("failed to re-login to passbolt: %w", err)
	}
	return nil
}
