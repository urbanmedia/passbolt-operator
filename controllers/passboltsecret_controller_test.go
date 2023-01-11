package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	passboltv1alpha1 "github.com/urbanmedia/passbolt-operator/api/v1alpha1"
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

	Context("Run existing job", func() {
		It("Should create successfully", func() {
			By("Create job for run")

			passboltSecretSpec := passboltv1alpha1.PassboltSecretSpec{
				LeaveOnDelete: false,
				Secrets: []passboltv1alpha1.SecretSpec{
					{
						KubernetesSecretKey: "password",
						PassboltSecret: passboltv1alpha1.PassboltSpec{
							Name:  "APP_EXAMPLE",
							Field: passboltv1alpha1.FieldNamePassword,
						},
					},
					{
						KubernetesSecretKey: "url",
						PassboltSecret: passboltv1alpha1.PassboltSpec{
							Name:  "APP_EXAMPLE",
							Field: passboltv1alpha1.FieldNameUri,
						},
					},
					{
						KubernetesSecretKey: "username",
						PassboltSecret: passboltv1alpha1.PassboltSpec{
							Name:  "APP_EXAMPLE",
							Field: passboltv1alpha1.FieldNameUsername,
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
			time.Sleep(time.Second * 60)
			defer func() {
				Expect(k8sClient.Delete(ctx, passboltSecret)).Should(Succeed())
				time.Sleep(time.Second * 10)
			}()

			By("By checking the PassboltSecret can be retrieved")
			passboltSecretKey := types.NamespacedName{Name: name, Namespace: namespace}
			passboltSecretObj := &passboltv1alpha1.PassboltSecret{}
			Eventually(func() error {
				return k8sClient.Get(ctx, passboltSecretKey, passboltSecretObj)
			}, timeout, interval).Should(Succeed())

			// By("By checking the PassboltSecret has been synced")
			// Expect(passboltSecretObj.Status.SyncStatus).Should(Equal(passboltv1alpha1.SyncStatusSuccess))

			By("By checking the PassboltSecret has been synced to Kubernetes Secret")
			kubernetesSecretKey := passboltSecretKey
			kubernetesSecretObj := &corev1.Secret{}
			Eventually(func() error {
				return k8sClient.Get(ctx, kubernetesSecretKey, kubernetesSecretObj)
			}, timeout, interval).Should(Succeed())

			By("By checking the Kubernetes Secret has three keys")
			Expect(kubernetesSecretObj.Data).Should(HaveLen(3))

			By("By checking the Kubernetes Secret has the username key")
			Expect(kubernetesSecretObj.Data).Should(HaveKey("username"))

			By("By checking the Kubernetes Secret has the url key")
			Expect(kubernetesSecretObj.Data).Should(HaveKey("url"))

			By("By checking the Kubernetes Secret has the password key")
			Expect(kubernetesSecretObj.Data).Should(HaveKey("password"))
		})
	})

	Context("Update existing secret", func() {
		It("Should update successfully", func() {
			By("Create job for run")
			passboltSecretSpec := passboltv1alpha1.PassboltSecretSpec{
				LeaveOnDelete: false,
				Secrets: []passboltv1alpha1.SecretSpec{
					{
						KubernetesSecretKey: "password",
						PassboltSecret: passboltv1alpha1.PassboltSpec{
							Name:  "APP_EXAMPLE",
							Field: passboltv1alpha1.FieldNamePassword,
						},
					},
					{
						KubernetesSecretKey: "url",
						PassboltSecret: passboltv1alpha1.PassboltSpec{
							Name:  "APP_EXAMPLE",
							Field: passboltv1alpha1.FieldNameUri,
						},
					},
					{
						KubernetesSecretKey: "username",
						PassboltSecret: passboltv1alpha1.PassboltSpec{
							Name:  "APP_EXAMPLE",
							Field: passboltv1alpha1.FieldNameUsername,
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
			time.Sleep(time.Second * 60)
			defer func() {
				Expect(k8sClient.Delete(ctx, passboltSecret)).Should(Succeed())
				time.Sleep(time.Second * 10)
			}()

			By("By checking the PassboltSecret can be retrieved")
			passboltSecretKey := types.NamespacedName{Name: name, Namespace: namespace}
			passboltSecretObj := &passboltv1alpha1.PassboltSecret{}
			Eventually(func() error {
				return k8sClient.Get(ctx, passboltSecretKey, passboltSecretObj)
			}, timeout, interval).Should(Succeed())

			// By("By checking the PassboltSecret has been synced")
			// Expect(passboltSecretObj.Status.SyncStatus).Should(Equal(passboltv1alpha1.SyncStatusSuccess))

			By("By checking the PassboltSecret has been synced to Kubernetes Secret")
			kubernetesSecretKey := passboltSecretKey
			kubernetesSecretObj := &corev1.Secret{}
			Eventually(func() error {
				return k8sClient.Get(ctx, kubernetesSecretKey, kubernetesSecretObj)
			}, timeout, interval).Should(Succeed())

			By("By checking the Kubernetes Secret has three keys")
			Expect(kubernetesSecretObj.Data).Should(HaveLen(3))

			By("By checking the Kubernetes Secret has the username key")
			Expect(kubernetesSecretObj.Data).Should(HaveKey("username"))

			By("By checking the Kubernetes Secret has the url key")
			Expect(kubernetesSecretObj.Data).Should(HaveKey("url"))

			By("By checking the Kubernetes Secret has the password key")
			Expect(kubernetesSecretObj.Data).Should(HaveKey("password"))

			//
			// do the same again but this time update the secret
			//

			By("By checking the PassboltSecret has been updated")
			// test if the passbolt secret could be updated
			// we expect that the secret can not be updated because the secret does not contain the necessary metadata
			// to update the secret
			// (Operation cannot be fulfilled on passboltsecrets.passbolt.tagesspiegel.de \"example-passboltsecret\": the object has been modified; please apply your changes to the latest version and try again")
			Expect(k8sClient.Update(ctx, passboltSecret)).ShouldNot(Succeed())
			time.Sleep(time.Second * 60)

			// this time we expect that the secret can be updated because the secret contains the necessary metadata
			By("By checking the PassboltSecret can be successfully updated")
			Eventually(func() error {
				return k8sClient.Get(ctx, passboltSecretKey, passboltSecretObj)
			}, timeout, interval).Should(Succeed())

			// update the actual value of the secret
			passboltSecret.Spec.Secrets = append(passboltSecret.Spec.Secrets, passboltv1alpha1.SecretSpec{
				PassboltSecret: passboltv1alpha1.PassboltSpec{
					Name:  "APP2_EXAMPLE",
					Field: passboltv1alpha1.FieldNamePassword,
				},
				KubernetesSecretKey: "app2_password",
			})
			passboltSecret.Spec.Secrets = append(passboltSecret.Spec.Secrets, passboltv1alpha1.SecretSpec{
				PassboltSecret: passboltv1alpha1.PassboltSpec{
					Name:  "APP2_EXAMPLE",
					Field: passboltv1alpha1.FieldNameUri,
				},
				KubernetesSecretKey: "app2_url",
			})
			passboltSecret.Spec.Secrets = append(passboltSecret.Spec.Secrets, passboltv1alpha1.SecretSpec{
				PassboltSecret: passboltv1alpha1.PassboltSpec{
					Name:  "APP2_EXAMPLE",
					Field: passboltv1alpha1.FieldNameUsername,
				},
				KubernetesSecretKey: "app2_username",
			})

			By("Expect local PassboltSecret to have six secrets")
			Expect(len(passboltSecret.Spec.Secrets)).Should(Equal(6))

			By("By checking the PassboltSecret has been updated")
			Expect(k8sClient.Update(ctx, passboltSecretObj)).Should(Succeed())
			time.Sleep(time.Second * 60)

			By("By checking the PassboltSecret has been synced to Kubernetes Secret again")
			kubernetesSecretObj = &corev1.Secret{}
			Eventually(func() error {
				return k8sClient.Get(ctx, kubernetesSecretKey, kubernetesSecretObj)
			}, timeout, interval).Should(Succeed())

			//By("By checking the Kubernetes Secret has six keys")
			//Expect(kubernetesSecretObj.Data).Should(HaveLen(6))

			/*
				By("By checking the Kubernetes Secret has six keys")
				Expect(kubernetesSecretObj.Data).Should(HaveLen(6))

				By("By checking the Kubernetes Secret has the username key again")
				Expect(kubernetesSecretObj.Data).Should(HaveKey("username"))

				By("By checking the Kubernetes Secret has the url key again")
				Expect(kubernetesSecretObj.Data).Should(HaveKey("url"))

				By("By checking the Kubernetes Secret has the password key again")
				Expect(kubernetesSecretObj.Data).Should(HaveKey("password"))

				By("By checking the Kubernetes Secret has the app2_username key again")
				Expect(kubernetesSecretObj.Data).Should(HaveKey("app2_username"))

				By("By checking the Kubernetes Secret has the app2_url key again")
				Expect(kubernetesSecretObj.Data).Should(HaveKey("app2_url"))

				By("By checking the Kubernetes Secret has the app2_password key again")
				Expect(kubernetesSecretObj.Data).Should(HaveKey("app2_password"))
			*/
		})
	})
})
