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
	passboltv1 "github.com/urbanmedia/passbolt-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

func TestMain(m *testing.M) {
	GetSecretID = func(name string) (string, error) {
		return "example-id", nil
	}
	GetSecretName = func(id string) (string, error) {
		return "example-name", nil
	}
	m.Run()
}

func TestPassboltSecret_ConvertTo(t *testing.T) {
	type fields struct {
		TypeMeta   metav1.TypeMeta
		ObjectMeta metav1.ObjectMeta
		Spec       PassboltSecretSpec
		Status     PassboltSecretStatus
	}
	type args struct {
		dstRaw conversion.Hub
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    conversion.Hub
		wantErr bool
	}{
		{
			name: "convert to v1alpha3 opaque",
			fields: fields{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example-passboltsecret",
					Namespace: "default",
				},
				Spec: PassboltSecretSpec{
					LeaveOnDelete: false,
					SecretType:    corev1.SecretTypeOpaque,
					Secrets: []SecretSpec{
						{
							PassboltSecret: PassboltSpec{
								Name:  "APP_EXAMPLE",
								Field: FieldNameUsername,
							},
							KubernetesSecretKey: "amqp_dsn",
						},
					},
				},
				Status: PassboltSecretStatus{
					SyncErrors: []SyncError{},
				},
			},
			args: args{
				dstRaw: &passboltv1.PassboltSecret{},
			},
			want: &passboltv1.PassboltSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example-passboltsecret",
					Namespace: "default",
				},
				Spec: passboltv1.PassboltSecretSpec{
					LeaveOnDelete: false,
					SecretType:    corev1.SecretTypeOpaque,
					PassboltSecrets: map[string]passboltv1.PassboltSecretRef{
						"amqp_dsn": {
							ID:    "example-id",
							Field: passboltv1.FieldNameUsername,
						},
					},
				},
				Status: passboltv1.PassboltSecretStatus{
					SyncErrors: []passboltv1.SyncError{},
				},
			},
			wantErr: false,
		},
		{
			name: "convert to v1alpha3 dockerconfigjson",
			fields: fields{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example-passboltsecret",
					Namespace: "default",
				},
				Spec: PassboltSecretSpec{
					LeaveOnDelete:      false,
					SecretType:         corev1.SecretTypeDockerConfigJson,
					PassboltSecretName: func() *string { s := "APP_EXAMPLE"; return &s }(),
				},
				Status: PassboltSecretStatus{
					SyncErrors: []SyncError{},
				},
			},
			args: args{
				dstRaw: &passboltv1.PassboltSecret{},
			},
			want: &passboltv1.PassboltSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example-passboltsecret",
					Namespace: "default",
				},
				Spec: passboltv1.PassboltSecretSpec{
					LeaveOnDelete:    false,
					SecretType:       corev1.SecretTypeDockerConfigJson,
					PassboltSecretID: func() *string { s := "example-id"; return &s }(),
				},
				Status: passboltv1.PassboltSecretStatus{
					SyncErrors: []passboltv1.SyncError{},
				},
			},
			wantErr: false,
		},
		{
			name: "convert to v1alpha3 with value set",
			fields: fields{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example-passboltsecret",
					Namespace: "default",
				},
				Spec: PassboltSecretSpec{
					LeaveOnDelete: false,
					SecretType:    corev1.SecretTypeOpaque,
					Secrets: []SecretSpec{
						{
							PassboltSecret: PassboltSpec{
								Name:  "APP_EXAMPLE",
								Value: func() *string { s := "{{.Username}}"; return &s }(),
							},
							KubernetesSecretKey: "amqp_dsn",
						},
					},
				},
				Status: PassboltSecretStatus{
					SyncErrors: []SyncError{},
				},
			},
			args: args{
				dstRaw: &passboltv1.PassboltSecret{},
			},
			want: &passboltv1.PassboltSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example-passboltsecret",
					Namespace: "default",
				},
				Spec: passboltv1.PassboltSecretSpec{
					LeaveOnDelete: false,
					SecretType:    corev1.SecretTypeOpaque,
					PassboltSecrets: map[string]passboltv1.PassboltSecretRef{
						"amqp_dsn": {
							ID:    "example-id",
							Value: func() *string { s := "{{.Username}}"; return &s }(),
						},
					},
				},
				Status: passboltv1.PassboltSecretStatus{
					SyncErrors: []passboltv1.SyncError{},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := &PassboltSecret{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			got := tt.args.dstRaw
			if err := src.ConvertTo(got); (err != nil) != tt.wantErr {
				t.Errorf("PassboltSecret.ConvertTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			diff := cmp.Diff(tt.want, got)
			if diff != "" {
				t.Errorf("PassboltSecret.ConvertTo() (-want, +got) = %v", diff)
				return
			}
		})
	}
}

func TestPassboltSecret_ConvertFrom(t *testing.T) {
	type args struct {
		srcRaw conversion.Hub
	}
	tests := []struct {
		name    string
		args    args
		want    *PassboltSecret
		wantErr bool
	}{
		{
			name: "convert from v1alpha3 with field name",
			args: args{
				srcRaw: &passboltv1.PassboltSecret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "example-passboltsecret",
						Namespace: "default",
					},
					Spec: passboltv1.PassboltSecretSpec{
						LeaveOnDelete: false,
						SecretType:    corev1.SecretTypeOpaque,
						PassboltSecrets: map[string]passboltv1.PassboltSecretRef{
							"amqp_dsn": {
								ID:    "example-id",
								Field: passboltv1.FieldNameUsername,
							},
						},
						PlainTextFields: map[string]string{
							"pg_dsn": "example-value",
						},
					},
				},
			},
			want: &PassboltSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example-passboltsecret",
					Namespace: "default",
				},
				Spec: PassboltSecretSpec{
					LeaveOnDelete: false,
					Secrets: []SecretSpec{
						{
							KubernetesSecretKey: "amqp_dsn",
							PassboltSecret: PassboltSpec{
								Name:  "example-name",
								Field: FieldNameUsername,
							},
						},
					},
				},
				Status: PassboltSecretStatus{
					SyncErrors: []SyncError{},
				},
			},
			wantErr: false,
		},
		{
			name: "convert from v1alpha3 dockerconfigjson",
			args: args{
				srcRaw: &passboltv1.PassboltSecret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "example-passboltsecret",
						Namespace: "default",
					},
					Spec: passboltv1.PassboltSecretSpec{
						LeaveOnDelete:    false,
						SecretType:       corev1.DockerConfigJsonKey,
						PassboltSecretID: func() *string { s := "184734ea-8be3-4f5a-ba6c-5f4b3c0603e8"; return &s }(),
					},
				},
			},
			want: &PassboltSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example-passboltsecret",
					Namespace: "default",
				},
				Spec: PassboltSecretSpec{
					LeaveOnDelete:      false,
					SecretType:         corev1.SecretTypeDockerConfigJson,
					PassboltSecretName: func() *string { s := "APP_EXAMPLE"; return &s }(),
				},
				Status: PassboltSecretStatus{
					SyncErrors: []SyncError{},
				},
			},
			wantErr: false,
		},
		{
			name: "convert from v1alpha3 with value",
			args: args{
				srcRaw: &passboltv1.PassboltSecret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "example-passboltsecret",
						Namespace: "default",
					},
					Spec: passboltv1.PassboltSecretSpec{
						LeaveOnDelete: false,
						SecretType:    corev1.SecretTypeOpaque,
						PassboltSecrets: map[string]passboltv1.PassboltSecretRef{
							"amqp_dsn": {
								ID:    "example-id",
								Value: func() *string { s := "example-value"; return &s }(),
							},
						},
						PlainTextFields: map[string]string{
							"pg_dsn": "example-value",
						},
					},
				},
			},
			want: &PassboltSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example-passboltsecret",
					Namespace: "default",
				},
				Spec: PassboltSecretSpec{
					LeaveOnDelete: false,
					Secrets:       []SecretSpec{},
				},
				Status: PassboltSecretStatus{
					SyncErrors: []SyncError{},
				},
			},
			wantErr: false,
		},
		{
			name: "convert from v1alpha3 with empty field",
			args: args{
				srcRaw: &passboltv1.PassboltSecret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "example-passboltsecret",
						Namespace: "default",
					},
					Spec: passboltv1.PassboltSecretSpec{
						LeaveOnDelete: false,
						SecretType:    corev1.SecretTypeOpaque,
						PassboltSecrets: map[string]passboltv1.PassboltSecretRef{
							"amqp_dsn": {
								ID:    "example-id",
								Field: "",
							},
						},
					},
				},
			},
			want: &PassboltSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example-passboltsecret",
					Namespace: "default",
				},
				Spec: PassboltSecretSpec{
					LeaveOnDelete: false,
					Secrets:       []SecretSpec{},
				},
				Status: PassboltSecretStatus{
					SyncErrors: []SyncError{},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.want
			if err := got.ConvertFrom(tt.args.srcRaw); (err != nil) != tt.wantErr {
				t.Errorf("PassboltSecret.ConvertFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			diff := cmp.Diff(tt.want, got)
			if diff != "" {
				t.Errorf("PassboltSecret.ConvertFrom() (-want, +got) = %v", diff)
				return
			}
		})
	}
}
