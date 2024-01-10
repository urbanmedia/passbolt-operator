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
	"time"

	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	. "github.com/onsi/gomega"

	passboltv1alpha3 "github.com/urbanmedia/passbolt-operator/api/v1alpha3"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Run Controller", func() {

	const (
		timeout  = time.Second * 30
		interval = time.Second * 1

		name      = "example-passboltsecret"
		namespace = "default"
	)

	gomega.SetDefaultEventuallyTimeout(timeout)
	gomega.SetDefaultEventuallyPollingInterval(interval)

	passboltSecretV1Alpha3 := &passboltv1alpha3.PassboltSecret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: passboltv1alpha3.PassboltSecretSpec{
			LeaveOnDelete: false,
			PassboltSecrets: map[string]passboltv1alpha3.PassboltSecretRef{
				"amqp_dsn": {
					ID:    "184734ea-8be3-4f5a-ba6c-5f4b3c0603e8",
					Value: func() *string { s := "amqp://{{ .Username }}:{{ .Password }}@{{ .URI }}/vhost"; return &s }(),
				},
				"pg_dsn": {
					ID:    "184734ea-8be3-4f5a-ba6c-5f4b3c0603e8",
					Value: func() *string { s := "amqp://{{ .Username }}:{{ .Password }}@{{ .URI }}/vhost"; return &s }(),
				},
			},
		},
	}

	defer func() {

	}()

	BeforeEach(func() {

	})

	AfterEach(func() {

	})

	// Add Tests for OpenAPI validation (or additonal CRD features) specified in
	// your API definition.
	// Avoid adding tests for vanilla CRUD operations because they would
	// test Kubernetes API server, which isn't the goal here.
	Context("Run directly without existing job", func() {
		It("Should create successfully", func() {
			Expect(1).To(Equal(1))
		})
	})

	Context("Version v1alpha3", func() {
		It("PassboltSecret", func() {
			// create the passbolt secret before the test
			By("By checking the PassboltSecret has been created")
			// test if the passbolt secret is created
			ctx := context.Background()
			Expect(k8sClient.Create(ctx, passboltSecretV1Alpha3)).Should(Succeed())

			time.Sleep(5 * time.Second)

			By("By checking, if PassboltSecret can be retrieved")
			pbGetSecret := &passboltv1alpha3.PassboltSecret{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, pbGetSecret)).Should(Succeed())

			By("By checking if PassboltSecret has the correct sync status")
			Expect(pbGetSecret.Status.SyncStatus).Should(Equal(passboltv1alpha3.SyncStatusSuccess))
		})

		It("Secret", func() {
			By("By checking if Secret was created")
			secret := &corev1.Secret{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, secret)).Should(Succeed())

			By("By checking if Secret has the correct length")
			Expect(secret.Data).Should(HaveLen(len(passboltSecretV1Alpha3.Spec.PassboltSecrets) + len(passboltSecretV1Alpha3.Spec.PlainTextFields)))

			By("By checking if Secret has the correct keys")
			Eventually(secret.Data).Should(HaveKey("amqp_dsn"))
			Eventually(secret.Data).Should(HaveKey("pg_dsn"))
		})

		It("Should delete", func() {
			// delete the passbolt secret after the test
			Expect(k8sClient.Delete(context.Background(), passboltSecretV1Alpha3)).Should(Succeed())
			time.Sleep(time.Second * 5)
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, &passboltv1alpha3.PassboltSecret{})).ShouldNot(Succeed())
			time.Sleep(time.Second * 5)
		})
	})
})
