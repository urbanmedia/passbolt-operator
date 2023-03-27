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

package v1alpha2

import (
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var passboltsecretlog = logf.Log.WithName("passboltsecret-resource")

func (r *PassboltSecret) SetupWebhookWithManager(mgr ctrl.Manager) error {
	passboltsecretlog.V(10).Info("setting up webhook", "version", "v1alpha2")
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-passbolt-tagesspiegel-de-v1alpha2-passboltsecret,mutating=true,failurePolicy=fail,sideEffects=None,groups=passbolt.tagesspiegel.de,resources=passboltsecrets,verbs=create;update,versions=v1alpha2,name=mpassboltsecret.kb.io,admissionReviewVersions=v1

// check if we have implemented the defaulter interface
var _ webhook.Defaulter = &PassboltSecret{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *PassboltSecret) Default() {
	passboltsecretlog.Info("default", "name", r.Name)
	if r.Spec.SecretType == "" {
		r.Spec.SecretType = corev1.SecretTypeOpaque
	}
}

//+kubebuilder:webhook:path=/validate-passbolt-tagesspiegel-de-v1alpha2-passboltsecret,mutating=false,failurePolicy=fail,sideEffects=None,groups=passbolt.tagesspiegel.de,resources=passboltsecrets,verbs=create;update,versions=v1alpha2,name=vpassboltsecret.kb.io,admissionReviewVersions=v1

var (
	ErrInvalidSecretType              = errors.New("invalid secret type")
	ErrPassboltSecretNameIsRequired   = errors.New("passboltSecretName is required for secret type")
	ErrSecretsAreNotAllowed           = errors.New("secrets are not allowed")
	ErrFieldAndValueAreNotAllowed     = errors.New("field and value are not allowed")
	ErrFieldOrValueIsRequired         = errors.New("field or value is required")
	ErrSecretsAreRequired             = errors.New("secrets are required")
	ErrPassboltSecretNameIsNotAllowed = errors.New("passboltSecretName is not allowed")
)

// check if we have implemented the validator interface
var _ webhook.Validator = &PassboltSecret{}

func (r *PassboltSecret) validatePassboltSecret() error {
	switch r.Spec.SecretType {
	case corev1.SecretTypeOpaque:
		if r.Spec.PassboltSecretName != nil {
			return fmt.Errorf("%w for secret type %s", ErrPassboltSecretNameIsNotAllowed, r.Spec.SecretType)
		}
		if len(r.Spec.Secrets) == 0 {
			return fmt.Errorf("%w for secret type %s", ErrSecretsAreRequired, r.Spec.SecretType)
		}
		// check if only FieldName or Value is set
		for _, secret := range r.Spec.Secrets {
			if secret.PassboltSecret.Field == "" && secret.PassboltSecret.Value == nil {
				return fmt.Errorf("%w for secret %s and field %v", ErrFieldOrValueIsRequired, r.GetName(), secret)
			}
			if secret.PassboltSecret.Field != "" && secret.PassboltSecret.Value != nil {
				return fmt.Errorf("%w for secret %s and field %v", ErrFieldAndValueAreNotAllowed, r.GetName(), secret)
			}
		}
		return nil
	case corev1.SecretTypeDockerConfigJson:
		if r.Spec.PassboltSecretName == nil {
			return fmt.Errorf("%w: %s", ErrPassboltSecretNameIsRequired, r.Spec.SecretType)
		}
		if *r.Spec.PassboltSecretName == "" {
			return fmt.Errorf("%w: %s", ErrPassboltSecretNameIsRequired, r.Spec.SecretType)
		}
		if len(r.Spec.Secrets) > 0 {
			return fmt.Errorf("%w for secret type %s", ErrSecretsAreNotAllowed, r.Spec.SecretType)
		}
		return nil
	default:
		return fmt.Errorf("%w: %s", ErrInvalidSecretType, r.Spec.SecretType)
	}
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *PassboltSecret) ValidateCreate() error {
	passboltsecretlog.Info("validate create", "name", r.Name)
	return r.validatePassboltSecret()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *PassboltSecret) ValidateUpdate(old runtime.Object) error {
	passboltsecretlog.Info("validate update", "name", r.Name)
	return r.validatePassboltSecret()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *PassboltSecret) ValidateDelete() error {
	passboltsecretlog.Info("validate delete", "name", r.Name)
	return r.validatePassboltSecret()
}
