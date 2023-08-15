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
	"testing"

	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	//ctrl "sigs.k8s.io/controller-runtime"
)

/*
func TestPassboltSecret_SetupWebhookWithManager(t *testing.T) {
	type fields struct {
		TypeMeta   metav1.TypeMeta
		ObjectMeta metav1.ObjectMeta
		Spec       PassboltSecretSpec
		Status     PassboltSecretStatus
	}
	type args struct {
		mgr ctrl.Manager
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &PassboltSecret{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			if err := r.SetupWebhookWithManager(tt.args.mgr); (err != nil) != tt.wantErr {
				t.Errorf("PassboltSecret.SetupWebhookWithManager() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
*/

func TestPassboltSecret_Default(t *testing.T) {
	type fields struct {
		TypeMeta   metav1.TypeMeta
		ObjectMeta metav1.ObjectMeta
		Spec       PassboltSecretSpec
		Status     PassboltSecretStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   PassboltSecret
	}{
		{
			name: "secret type not set",
			fields: fields{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
				},
				Spec:   PassboltSecretSpec{},
				Status: PassboltSecretStatus{},
			},
			want: PassboltSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
				},
				Spec: PassboltSecretSpec{
					SecretType: corev1.SecretTypeOpaque,
				},
			},
		},
		{
			name: "secret type is opaque",
			fields: fields{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
				},
				Spec: PassboltSecretSpec{
					SecretType: corev1.SecretTypeOpaque,
				},
				Status: PassboltSecretStatus{},
			},
			want: PassboltSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
				},
				Spec: PassboltSecretSpec{
					SecretType: corev1.SecretTypeOpaque,
				},
			},
		},
		{
			name: "secret type is SecretTypeDockerConfigJson",
			fields: fields{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
				},
				Spec: PassboltSecretSpec{
					SecretType: corev1.SecretTypeDockerConfigJson,
				},
				Status: PassboltSecretStatus{},
			},
			want: PassboltSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
				},
				Spec: PassboltSecretSpec{
					SecretType: corev1.SecretTypeDockerConfigJson,
				},
			},
		},
		{
			name: "secret type is SecretTypeTLS",
			fields: fields{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
				},
				Spec: PassboltSecretSpec{
					SecretType: corev1.SecretTypeTLS,
				},
				Status: PassboltSecretStatus{},
			},
			want: PassboltSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
				},
				Spec: PassboltSecretSpec{
					SecretType: corev1.SecretTypeOpaque,
				},
			},
		},
		{
			name: "secret type is SecretTypeBasicAuth",
			fields: fields{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
				},
				Spec: PassboltSecretSpec{
					SecretType: corev1.SecretTypeBasicAuth,
				},
				Status: PassboltSecretStatus{},
			},
			want: PassboltSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
				},
				Spec: PassboltSecretSpec{
					SecretType: corev1.SecretTypeOpaque,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &PassboltSecret{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			r.Default()
			if diff := cmp.Diff(*r, tt.want); diff != "" {
				t.Errorf("PassboltSecret.Default() diff = %s", diff)
			}
		})
	}
}

func TestPassboltSecret_validatePassboltSecret(t *testing.T) {
	type fields struct {
		TypeMeta   metav1.TypeMeta
		ObjectMeta metav1.ObjectMeta
		Spec       PassboltSecretSpec
		Status     PassboltSecretStatus
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// opaque secret
		{
			name: "valid Opaque secret field name is set",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeOpaque,
					Secrets: []SecretSpec{
						{
							PassboltSecret: PassboltSpec{
								Name:  "test",
								Field: FieldNamePassword,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid Opaque secret value is set",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeOpaque,
					Secrets: []SecretSpec{
						{
							PassboltSecret: PassboltSpec{
								Name:  "test",
								Value: func() *string { s := "host={{.URI}}"; return &s }(),
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid Opaque secret passboltSecretName set",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeOpaque,
					PassboltSecretName: func() *string {
						s := "test"
						return &s
					}(),
					Secrets: []SecretSpec{
						{
							PassboltSecret: PassboltSpec{
								Name:  "test",
								Field: FieldNamePassword,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid Opaque secret length of secrets is 0",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeOpaque,
					Secrets:       []SecretSpec{},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid Opaque secret field or value is not set",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeOpaque,
					Secrets: []SecretSpec{
						{
							PassboltSecret:      PassboltSpec{},
							KubernetesSecretKey: "test",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid Opaque secret field and value is set",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeOpaque,
					Secrets: []SecretSpec{
						{
							PassboltSecret: PassboltSpec{
								Field: FieldNamePassword,
								Value: func() *string { s := "host={{.URI}}"; return &s }(),
							},
							KubernetesSecretKey: "test",
						},
					},
				},
			},
			wantErr: true,
		},
		// dockerconfigjson secret
		{
			name: "valid DockerConfigJson secret",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete:      true,
					SecretType:         corev1.SecretTypeDockerConfigJson,
					PassboltSecretName: func() *string { s := "test"; return &s }(),
				},
			},
			wantErr: false,
		},
		{
			name: "invalid DockerConfigJson secret PassboltSecretName is not set",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeDockerConfigJson,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid DockerConfigJson secret PassboltSecretName is empty",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete:      true,
					SecretType:         corev1.SecretTypeDockerConfigJson,
					PassboltSecretName: func() *string { s := ""; return &s }(),
				},
			},
			wantErr: true,
		},
		{
			name: "invalid DockerConfigJson secret secrets is set",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete:      true,
					SecretType:         corev1.SecretTypeDockerConfigJson,
					PassboltSecretName: func() *string { s := "as"; return &s }(),
					Secrets: []SecretSpec{
						{
							PassboltSecret:      PassboltSpec{},
							KubernetesSecretKey: "asd",
						},
					},
				},
			},
			wantErr: true,
		},
		// unsupported secret type
		{
			name: "invalid secret type SecretTypeBasicAuth",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeBasicAuth,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid secret type SecretTypeBootstrapToken",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeBootstrapToken,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid secret type SecretTypeDockercfg",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeDockercfg,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid secret type SecretTypeSSHAuth",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeSSHAuth,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid secret type SecretTypeServiceAccountToken",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeServiceAccountToken,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid secret type SecretTypeTLS",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeTLS,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &PassboltSecret{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			if err := r.validatePassboltSecret(); (err != nil) != tt.wantErr {
				t.Errorf("PassboltSecret.validatePassboltSecret() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPassboltSecret_ValidateCreate(t *testing.T) {
	type fields struct {
		TypeMeta   metav1.TypeMeta
		ObjectMeta metav1.ObjectMeta
		Spec       PassboltSecretSpec
		Status     PassboltSecretStatus
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// opaque secret
		{
			name: "valid Opaque secret field name is set",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeOpaque,
					Secrets: []SecretSpec{
						{
							PassboltSecret: PassboltSpec{
								Name:  "test",
								Field: FieldNamePassword,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid Opaque secret value is set",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeOpaque,
					Secrets: []SecretSpec{
						{
							PassboltSecret: PassboltSpec{
								Name:  "test",
								Value: func() *string { s := "host={{.URI}}"; return &s }(),
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid Opaque secret passboltSecretName set",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeOpaque,
					PassboltSecretName: func() *string {
						s := "test"
						return &s
					}(),
					Secrets: []SecretSpec{
						{
							PassboltSecret: PassboltSpec{
								Name:  "test",
								Field: FieldNamePassword,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid Opaque secret length of secrets is 0",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeOpaque,
					Secrets:       []SecretSpec{},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid Opaque secret field or value is not set",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeOpaque,
					Secrets: []SecretSpec{
						{
							PassboltSecret:      PassboltSpec{},
							KubernetesSecretKey: "test",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid Opaque secret field and value is set",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeOpaque,
					Secrets: []SecretSpec{
						{
							PassboltSecret: PassboltSpec{
								Field: FieldNamePassword,
								Value: func() *string { s := "host={{.URI}}"; return &s }(),
							},
							KubernetesSecretKey: "test",
						},
					},
				},
			},
			wantErr: true,
		},
		// dockerconfigjson secret
		{
			name: "valid DockerConfigJson secret",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete:      true,
					SecretType:         corev1.SecretTypeDockerConfigJson,
					PassboltSecretName: func() *string { s := "test"; return &s }(),
				},
			},
			wantErr: false,
		},
		{
			name: "invalid DockerConfigJson secret PassboltSecretName is not set",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeDockerConfigJson,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid DockerConfigJson secret PassboltSecretName is empty",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete:      true,
					SecretType:         corev1.SecretTypeDockerConfigJson,
					PassboltSecretName: func() *string { s := ""; return &s }(),
				},
			},
			wantErr: true,
		},
		{
			name: "invalid DockerConfigJson secret secrets is set",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete:      true,
					SecretType:         corev1.SecretTypeDockerConfigJson,
					PassboltSecretName: func() *string { s := "as"; return &s }(),
					Secrets: []SecretSpec{
						{
							PassboltSecret:      PassboltSpec{},
							KubernetesSecretKey: "asd",
						},
					},
				},
			},
			wantErr: true,
		},
		// unsupported secret type
		{
			name: "invalid secret type SecretTypeBasicAuth",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeBasicAuth,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid secret type SecretTypeBootstrapToken",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeBootstrapToken,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid secret type SecretTypeDockercfg",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeDockercfg,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid secret type SecretTypeSSHAuth",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeSSHAuth,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid secret type SecretTypeServiceAccountToken",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeServiceAccountToken,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid secret type SecretTypeTLS",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeTLS,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &PassboltSecret{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			if _, err := r.ValidateCreate(); (err != nil) != tt.wantErr {
				t.Errorf("PassboltSecret.ValidateCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPassboltSecret_ValidateUpdate(t *testing.T) {
	type fields struct {
		TypeMeta   metav1.TypeMeta
		ObjectMeta metav1.ObjectMeta
		Spec       PassboltSecretSpec
		Status     PassboltSecretStatus
	}
	type args struct {
		old runtime.Object
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// opaque secret
		{
			name: "valid Opaque secret field name is set",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeOpaque,
					Secrets: []SecretSpec{
						{
							PassboltSecret: PassboltSpec{
								Name:  "test",
								Field: FieldNamePassword,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid Opaque secret value is set",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeOpaque,
					Secrets: []SecretSpec{
						{
							PassboltSecret: PassboltSpec{
								Name:  "test",
								Value: func() *string { s := "host={{.URI}}"; return &s }(),
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid Opaque secret passboltSecretName set",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeOpaque,
					PassboltSecretName: func() *string {
						s := "test"
						return &s
					}(),
					Secrets: []SecretSpec{
						{
							PassboltSecret: PassboltSpec{
								Name:  "test",
								Field: FieldNamePassword,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid Opaque secret length of secrets is 0",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeOpaque,
					Secrets:       []SecretSpec{},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid Opaque secret field or value is not set",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeOpaque,
					Secrets: []SecretSpec{
						{
							PassboltSecret:      PassboltSpec{},
							KubernetesSecretKey: "test",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid Opaque secret field and value is set",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeOpaque,
					Secrets: []SecretSpec{
						{
							PassboltSecret: PassboltSpec{
								Field: FieldNamePassword,
								Value: func() *string { s := "host={{.URI}}"; return &s }(),
							},
							KubernetesSecretKey: "test",
						},
					},
				},
			},
			wantErr: true,
		},
		// dockerconfigjson secret
		{
			name: "valid DockerConfigJson secret",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete:      true,
					SecretType:         corev1.SecretTypeDockerConfigJson,
					PassboltSecretName: func() *string { s := "test"; return &s }(),
				},
			},
			wantErr: false,
		},
		{
			name: "invalid DockerConfigJson secret PassboltSecretName is not set",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeDockerConfigJson,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid DockerConfigJson secret PassboltSecretName is empty",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete:      true,
					SecretType:         corev1.SecretTypeDockerConfigJson,
					PassboltSecretName: func() *string { s := ""; return &s }(),
				},
			},
			wantErr: true,
		},
		{
			name: "invalid DockerConfigJson secret secrets is set",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete:      true,
					SecretType:         corev1.SecretTypeDockerConfigJson,
					PassboltSecretName: func() *string { s := "as"; return &s }(),
					Secrets: []SecretSpec{
						{
							PassboltSecret:      PassboltSpec{},
							KubernetesSecretKey: "asd",
						},
					},
				},
			},
			wantErr: true,
		},
		// unsupported secret type
		{
			name: "invalid secret type SecretTypeBasicAuth",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeBasicAuth,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid secret type SecretTypeBootstrapToken",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeBootstrapToken,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid secret type SecretTypeDockercfg",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeDockercfg,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid secret type SecretTypeSSHAuth",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeSSHAuth,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid secret type SecretTypeServiceAccountToken",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeServiceAccountToken,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid secret type SecretTypeTLS",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeTLS,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &PassboltSecret{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			if _, err := r.ValidateUpdate(tt.args.old); (err != nil) != tt.wantErr {
				t.Errorf("PassboltSecret.ValidateUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPassboltSecret_ValidateDelete(t *testing.T) {
	type fields struct {
		TypeMeta   metav1.TypeMeta
		ObjectMeta metav1.ObjectMeta
		Spec       PassboltSecretSpec
		Status     PassboltSecretStatus
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// opaque secret
		{
			name: "valid Opaque secret field name is set",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeOpaque,
					Secrets: []SecretSpec{
						{
							PassboltSecret: PassboltSpec{
								Name:  "test",
								Field: FieldNamePassword,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid Opaque secret value is set",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeOpaque,
					Secrets: []SecretSpec{
						{
							PassboltSecret: PassboltSpec{
								Name:  "test",
								Value: func() *string { s := "host={{.URI}}"; return &s }(),
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid Opaque secret passboltSecretName set",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeOpaque,
					PassboltSecretName: func() *string {
						s := "test"
						return &s
					}(),
					Secrets: []SecretSpec{
						{
							PassboltSecret: PassboltSpec{
								Name:  "test",
								Field: FieldNamePassword,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid Opaque secret length of secrets is 0",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeOpaque,
					Secrets:       []SecretSpec{},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid Opaque secret field or value is not set",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeOpaque,
					Secrets: []SecretSpec{
						{
							PassboltSecret:      PassboltSpec{},
							KubernetesSecretKey: "test",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid Opaque secret field and value is set",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeOpaque,
					Secrets: []SecretSpec{
						{
							PassboltSecret: PassboltSpec{
								Field: FieldNamePassword,
								Value: func() *string { s := "host={{.URI}}"; return &s }(),
							},
							KubernetesSecretKey: "test",
						},
					},
				},
			},
			wantErr: true,
		},
		// dockerconfigjson secret
		{
			name: "valid DockerConfigJson secret",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete:      true,
					SecretType:         corev1.SecretTypeDockerConfigJson,
					PassboltSecretName: func() *string { s := "test"; return &s }(),
				},
			},
			wantErr: false,
		},
		{
			name: "invalid DockerConfigJson secret PassboltSecretName is not set",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeDockerConfigJson,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid DockerConfigJson secret PassboltSecretName is empty",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete:      true,
					SecretType:         corev1.SecretTypeDockerConfigJson,
					PassboltSecretName: func() *string { s := ""; return &s }(),
				},
			},
			wantErr: true,
		},
		{
			name: "invalid DockerConfigJson secret secrets is set",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete:      true,
					SecretType:         corev1.SecretTypeDockerConfigJson,
					PassboltSecretName: func() *string { s := "as"; return &s }(),
					Secrets: []SecretSpec{
						{
							PassboltSecret:      PassboltSpec{},
							KubernetesSecretKey: "asd",
						},
					},
				},
			},
			wantErr: true,
		},
		// unsupported secret type
		{
			name: "invalid secret type SecretTypeBasicAuth",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeBasicAuth,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid secret type SecretTypeBootstrapToken",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeBootstrapToken,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid secret type SecretTypeDockercfg",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeDockercfg,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid secret type SecretTypeSSHAuth",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeSSHAuth,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid secret type SecretTypeServiceAccountToken",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeServiceAccountToken,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid secret type SecretTypeTLS",
			fields: fields{
				Spec: PassboltSecretSpec{
					LeaveOnDelete: true,
					SecretType:    corev1.SecretTypeTLS,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &PassboltSecret{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			if _, err := r.ValidateDelete(); (err != nil) != tt.wantErr {
				t.Errorf("PassboltSecret.ValidateDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
