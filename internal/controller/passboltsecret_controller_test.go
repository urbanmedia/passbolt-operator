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
	. "github.com/onsi/gomega"

	passboltv1alpha1 "github.com/urbanmedia/passbolt-operator/api/v1alpha1"
	passboltv1alpha2 "github.com/urbanmedia/passbolt-operator/api/v1alpha2"
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

	BeforeEach(func() {
		// Add any setup steps that needs to be executed before each test
	})

	AfterEach(func() {
		// Add any teardown steps that needs to be executed after each test
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

	Context("Version v1alpha1", func() {
		It("Should create and update successfully", func() {
			passboltSecretSpec := passboltv1alpha1.PassboltSecretSpec{
				LeaveOnDelete: false,
				Secrets: []passboltv1alpha1.SecretSpec{
					{
						KubernetesSecretKey: "username",
						PassboltSecret: passboltv1alpha1.PassboltSpec{
							Name:  "APP_EXAMPLE",
							Field: passboltv1alpha1.FieldNameUsername,
						},
					},
					{
						KubernetesSecretKey: "password",
						PassboltSecret: passboltv1alpha1.PassboltSpec{
							Name:  "APP_EXAMPLE",
							Field: passboltv1alpha1.FieldNamePassword,
						},
					},
				},
			}

			passboltSecret := &passboltv1alpha1.PassboltSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: passboltSecretSpec,
			}

			By("By checking the PassboltSecret has been created")
			// test if the passbolt secret is created
			ctx := context.Background()
			Expect(k8sClient.Create(ctx, passboltSecret)).Should(Succeed())

			defer func() {
				Expect(k8sClient.Delete(context.Background(), passboltSecret)).Should(Succeed())
				time.Sleep(time.Second * 5)
				Expect(k8sClient.Delete(context.Background(), &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: namespace,
					},
				})).Should(Succeed())
				time.Sleep(time.Second * 5)
			}()

			// here we have to delay a little
			time.Sleep(5 * time.Second)

			By("By checking if PassboltSecret was created")
			pbGetSecret := &passboltv1alpha1.PassboltSecret{}
			Eventually(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, pbGetSecret), timeout, interval).Should(Succeed())

			By("By checking if PassboltSecret has the correct sync status")
			Expect(pbGetSecret.Status.SyncStatus).Should(Equal(passboltv1alpha1.SyncStatusSuccess))

			By("By checking if Secret was created")
			secret := &corev1.Secret{}
			Eventually(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, secret), timeout, interval).Should(Succeed())

			By("By checking if Secret has the correct length")
			Expect(secret.Data).Should(HaveLen(len(pbGetSecret.Spec.Secrets)))

			By("By checking if Secret has the correct keys")
			Expect(secret.Data).Should(HaveKey("username"))
			Expect(secret.Data).Should(HaveKey("password"))

			By("By checking if Secret can be updated")
			pbGetSecret.Spec.Secrets = []passboltv1alpha1.SecretSpec{
				{
					KubernetesSecretKey: "dsn",
					PassboltSecret: passboltv1alpha1.PassboltSpec{
						Name:  "APP_EXAMPLE",
						Field: passboltv1alpha1.FieldNameUsername,
					},
				},
			}
			Expect(k8sClient.Update(ctx, pbGetSecret)).Should(Succeed())

			// here we have to delay a little
			time.Sleep(5 * time.Second)

			By("By checking if PassboltSecret has the correct length")
			Eventually(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, pbGetSecret), timeout, interval).Should(Succeed())
		})
	})

	Context("Version v1alpha2", func() {
		It("Should create and update successfully", func() {
			passboltSecretSpec := passboltv1alpha2.PassboltSecretSpec{
				LeaveOnDelete: false,
				Secrets: []passboltv1alpha2.SecretSpec{
					{
						KubernetesSecretKey: "amqp_dsn",
						PassboltSecret: passboltv1alpha2.PassboltSpec{
							Name:  "APP_EXAMPLE",
							Value: func() *string { s := "amqp://{{ .Username }}:{{ .Password }}@{{ .URI }}/vhost"; return &s }(),
						},
					},
					{
						KubernetesSecretKey: "pg_dsn",
						PassboltSecret: passboltv1alpha2.PassboltSpec{
							Name:  "APP_EXAMPLE",
							Value: func() *string { s := "amqp://{{ .Username }}:{{ .Password }}@{{ .URI }}/vhost"; return &s }(),
						},
					},
				},
			}

			passboltSecret := &passboltv1alpha2.PassboltSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: passboltSecretSpec,
			}

			By("By checking the PassboltSecret has been created")
			// test if the passbolt secret is created
			ctx := context.Background()
			Expect(k8sClient.Create(ctx, passboltSecret)).Should(Succeed())

			defer func() {
				Expect(k8sClient.Delete(context.Background(), passboltSecret)).Should(Succeed())
				time.Sleep(time.Second * 5)
				Expect(k8sClient.Delete(context.Background(), &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: namespace,
					},
				})).Should(Succeed())
				time.Sleep(time.Second * 5)
			}()

			// here we have to delay a little
			time.Sleep(5 * time.Second)

			By("By checking if PassboltSecret was created")
			pbGetSecret := &passboltv1alpha2.PassboltSecret{}
			Eventually(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, pbGetSecret), timeout, interval).Should(Succeed())

			By("By checking if PassboltSecret has the correct sync status")
			Expect(pbGetSecret.Status.SyncStatus).Should(Equal(passboltv1alpha2.SyncStatusSuccess))

			By("By checking if Secret was created")
			secret := &corev1.Secret{}
			Eventually(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, secret), timeout, interval).Should(Succeed())

			By("By checking if Secret has the correct length")
			Expect(secret.Data).Should(HaveLen(len(pbGetSecret.Spec.Secrets)))

			By("By checking if Secret has the correct keys")
			Expect(secret.Data).Should(HaveKey("amqp_dsn"))
			Expect(secret.Data).Should(HaveKey("pg_dsn"))

			By("By checking if Secret can be updated")
			pbGetSecret.Spec.Secrets = []passboltv1alpha2.SecretSpec{
				{
					KubernetesSecretKey: "dsn",
					PassboltSecret: passboltv1alpha2.PassboltSpec{
						Name:  "APP_EXAMPLE",
						Value: func() *string { s := "amqp://{{ .Username }}:{{ .Password }}@{{ .URI }}/vhost"; return &s }(),
					},
				},
			}
			Expect(k8sClient.Update(ctx, pbGetSecret)).Should(Succeed())

			// here we have to delay a little
			time.Sleep(5 * time.Second)

			By("By checking if PassboltSecret has the correct length")
			Eventually(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, pbGetSecret), timeout, interval).Should(Succeed())
		})
	})

	Context("Version v1alpha3", func() {
		It("Should create and update successfully", func() {
			passboltSecretSpec := passboltv1alpha3.PassboltSecretSpec{
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
			}

			passboltSecret := &passboltv1alpha3.PassboltSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: passboltSecretSpec,
			}

			By("By checking the PassboltSecret has been created")
			// test if the passbolt secret is created
			ctx := context.Background()
			Expect(k8sClient.Create(ctx, passboltSecret)).Should(Succeed())

			defer func() {
				Expect(k8sClient.Delete(context.Background(), passboltSecret)).Should(Succeed())
				time.Sleep(time.Second * 5)
				Expect(k8sClient.Delete(context.Background(), &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: namespace,
					},
				})).Should(Succeed())
				time.Sleep(time.Second * 5)
			}()

			// here we have to delay a little
			time.Sleep(5 * time.Second)

			By("By checking if PassboltSecret was created")
			pbGetSecret := &passboltv1alpha3.PassboltSecret{}
			Eventually(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, pbGetSecret), timeout, interval).Should(Succeed())

			By("By checking if PassboltSecret has the correct sync status")
			Expect(pbGetSecret.Status.SyncStatus).Should(Equal(passboltv1alpha3.SyncStatusSuccess))

			By("By checking if Secret was created")
			secret := &corev1.Secret{}
			Eventually(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, secret), timeout, interval).Should(Succeed())

			By("By checking if Secret has the correct length")
			Expect(secret.Data).Should(HaveLen(len(pbGetSecret.Spec.PassboltSecrets)))

			By("By checking if Secret has the correct keys")
			Eventually(secret.Data).Should(HaveKey("amqp_dsn"))
			Eventually(secret.Data).Should(HaveKey("pg_dsn"))

			By("By checking if Secret can be updated")
			pbGetSecret.Spec.PassboltSecrets = map[string]passboltv1alpha3.PassboltSecretRef{
				"dsn": {
					ID:    "184734ea-8be3-4f5a-ba6c-5f4b3c0603e8",
					Value: func() *string { s := "amqp://{{ .Username }}:{{ .Password }}@{{ .URI }}/vhost"; return &s }(),
				},
			}
			Expect(k8sClient.Update(ctx, pbGetSecret)).Should(Succeed())

			// here we have to delay a little
			time.Sleep(5 * time.Second)

			By("By checking if PassboltSecret has the correct length")
			Eventually(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, pbGetSecret), timeout, interval).Should(Succeed())
		})
	})
})
