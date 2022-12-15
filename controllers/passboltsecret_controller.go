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

package controllers

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	corev1 "k8s.io/api/core/v1"

	passboltv1alpha1 "github.com/urbanmedia/passbolt-operator/api/v1alpha1"
	"github.com/urbanmedia/passbolt-operator/pkg/passbolt"
)

// PassboltSecretReconciler reconciles a PassboltSecret object
type PassboltSecretReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=passbolt.tagesspiegel.de,resources=passboltsecrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=passbolt.tagesspiegel.de,resources=passboltsecrets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=passbolt.tagesspiegel.de,resources=passboltsecrets/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;create;update;delete;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PassboltSecret object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *PassboltSecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logr := log.FromContext(ctx)

	secret := &passboltv1alpha1.PassboltSecret{}
	err := r.Client.Get(ctx, req.NamespacedName, secret)
	if err != nil {
		// If the resource no longer exists, in which case we stop processing.
		logr.Error(err, "unable to fetch PassboltSecret")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// FIXME: pass real values to the passbolt client
	// create passbolt client
	clnt, err := passbolt.NewClient(ctx, "", "", "")
	if err != nil {
		logr.Error(err, "unable to create passbolt client")
		return ctrl.Result{}, err
	}
	defer clnt.Close(ctx)

	// retrieve secrets from passbolt and store them in Kubernetes secrets
	k8sSecret := &corev1.Secret{
		ObjectMeta: ctrl.ObjectMeta{
			Name:        secret.Name,
			Namespace:   secret.Namespace,
			Labels:      secret.Labels,
			Annotations: secret.Annotations,
		},
	}
	for _, secret := range secret.Spec.Secrets {
		secretData, err := clnt.GetSecret(ctx, secret.Name)
		if err != nil {
			logr.Error(err, "unable to retrieve secret from passbolt")
			return ctrl.Result{}, err
		}
		k8sSecret.StringData[secret.KubernetesSecretKey] = secretData
	}

	// check if the secret already exists
	err = r.Client.Get(ctx, req.NamespacedName, &corev1.Secret{})
	if err != nil && errors.IsNotFound(err) {
		// fail if the error is not a not found error
		return ctrl.Result{}, fmt.Errorf("unable to retrieve secret from kubernetes: %w", err)
	}
	// check if the secret already exists
	if errors.IsNotFound(err) {
		err = r.Client.Create(ctx, k8sSecret)
		if err != nil {
			logr.Error(err, "unable to create kubernetes secret")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}
	// secret already exists, update it
	err = r.Client.Update(ctx, k8sSecret)
	if err != nil {
		logr.Error(err, "unable to update kubernetes secret")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PassboltSecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&passboltv1alpha1.PassboltSecret{}).
		Complete(r)
}
