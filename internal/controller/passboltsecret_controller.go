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

package controller

import (
	"context"
	"errors"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	passboltv1 "github.com/urbanmedia/passbolt-operator/api/v1"
	"github.com/urbanmedia/passbolt-operator/pkg/passbolt"
	"github.com/urbanmedia/passbolt-operator/pkg/util"
)

// PassboltSecretReconciler reconciles a PassboltSecret object
type PassboltSecretReconciler struct {
	client.Client
	Scheme         *runtime.Scheme
	PassboltClient *passbolt.Client
}

var (
	errResult = ctrl.Result{
		Requeue:      true,
		RequeueAfter: 30 * time.Second,
	}
)

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
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *PassboltSecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logr := log.FromContext(ctx)
	logr.Info("starting reconciliation...", "name", req.NamespacedName)
	defer logr.Info("finished reconciliation", "name", req.NamespacedName)
	defer func() {
		if r := recover(); r != nil {
			logr.Error(errors.New("recovered from panic"), "failed to complete reconciliation", "stacktrace", r)
		}
	}()

	// get passbolt secret resource from Kubernetes
	secret := &passboltv1.PassboltSecret{}
	err := r.Client.Get(ctx, req.NamespacedName, secret)
	if err != nil {
		if err = client.IgnoreNotFound(err); err != nil {
			return errResult, err
		}
		return errResult, err
	}

	if secret.Status.FailureCount >= 3 {
		// if the secret failed to sync more than 3 times, we stop trying
		logr.Info("secret failed to sync more than 3 times. stopping sync", "name", secret.GetName(), "namespace", secret.GetNamespace())
		return ctrl.Result{}, nil
	}
	errResult.RequeueAfter = time.Duration(2^secret.Status.FailureCount) * (5 * time.Second)

	// cleanup status
	secret.Status.SyncErrors = []passboltv1.SyncError{}

	if secret.Spec.PassboltSecretID == nil && secret.Spec.PassboltSecrets == nil && secret.Spec.PlainTextFields == nil {
		return errResult, fmt.Errorf("no passbolt secret id, passbolt secret references or plain text fields defined")
	}

	// make sure that the secret type is supported
	if secret.Spec.SecretType != corev1.SecretTypeOpaque && secret.Spec.SecretType != corev1.SecretTypeDockerConfigJson {
		logr.Info("unsupported secret type", "type", secret.Spec.SecretType)
		secret.Status.SyncStatus = passboltv1.SyncStatusError
		secret.Status.FailureCount++
		secret.Status.SyncErrors = append(secret.Status.SyncErrors, passboltv1.SyncError{
			Message: fmt.Sprintf("unsupported secret type %q", secret.Spec.SecretType),
			Time:    metav1.Now(),
		})
		if err := r.Client.Status().Update(ctx, secret); err != nil {
			return errResult, err
		}
		return errResult, nil
	}

	// define Kubernetes secret to be created or updated
	k8sSecret := &corev1.Secret{
		ObjectMeta: ctrl.ObjectMeta{
			Name:        secret.Name,
			Namespace:   secret.Namespace,
			Labels:      secret.Labels,
			Annotations: secret.Annotations,
		},
		Type: secret.Spec.SecretType,
		Data: map[string][]byte{},
	}

	opRslt, err := controllerutil.CreateOrUpdate(ctx, r.Client, k8sSecret, util.UpdateSecret(ctx, r.PassboltClient, r.Scheme, secret, k8sSecret))
	if err != nil {
		if snErr, ok := err.(passboltv1.SyncError); ok {
			secret.Status.SyncStatus = passboltv1.SyncStatusError
			secret.Status.FailureCount++
			secret.Status.SyncErrors = append(secret.Status.SyncErrors, snErr)
			if err := r.Client.Status().Update(ctx, secret); err != nil {
				return errResult, err
			}
			return errResult, err
		}
		return errResult, err
	}

	// if the secret was not changed and the status is already success, we can skip the update
	if opRslt == controllerutil.OperationResultNone && secret.Status.SyncStatus == passboltv1.SyncStatusSuccess {
		// secret was not changed
		logr.V(10).Info("secret was not changed! skipping... ")
		return ctrl.Result{}, nil
	}

	// update status
	secret.Status.SyncStatus = passboltv1.SyncStatusSuccess
	secret.Status.LastSync = metav1.Now()
	err = r.Client.Status().Update(ctx, secret)
	if err != nil {
		// the secret was synced successfully but the status could not be updated
		return reconcile.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PassboltSecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&passboltv1.PassboltSecret{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}
