/*
Copyright 2023 Verlag der Tagesspiegel GmbH.

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

package v1alpha3

import (
	"errors"

	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	ErrInvalidSecretType              = errors.New("invalid secret type")
	ErrPassboltSecretNameIsRequired   = errors.New("passboltSecretName is required for secret type")
	ErrSecretsAreNotAllowed           = errors.New("secrets are not allowed")
	ErrFieldAndValueAreNotAllowed     = errors.New("field and value are not allowed")
	ErrFieldOrValueIsRequired         = errors.New("field or value is required")
	ErrSecretsAreRequired             = errors.New("secrets are required")
	ErrPassboltSecretNameIsNotAllowed = errors.New("passboltSecretName is not allowed")
)

// log is for logging in this package.
var passboltsecretlog = logf.Log.WithName("passboltsecret-resource")

func (r *PassboltSecret) SetupWebhookWithManager(mgr ctrl.Manager) error {
	passboltsecretlog.V(100).Info("setting up webhook")
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}
