package util

import (
	"context"
	"log"
	"testing"

	"github.com/google/go-cmp/cmp"
	passboltv1 "github.com/urbanmedia/passbolt-operator/api/v1"
	"github.com/urbanmedia/passbolt-operator/pkg/passbolt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	// passboltURL is the URL of the passbolt server.
	passboltURL = "http://localhost:8088"
	// passboltUsername is the username of the passbolt user.
	passboltUsername = `-----BEGIN PGP PRIVATE KEY BLOCK-----

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
=jwOw
-----END PGP PRIVATE KEY BLOCK-----`
	// passboltPassword is the password of the passbolt user.
	passboltPassword = "TestTest123!"
)

var (
	client *passbolt.Client
	scheme *runtime.Scheme = runtime.NewScheme()
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	clnt, err := passbolt.NewClient(ctx, passboltURL, passboltUsername, passboltPassword)
	if err != nil {
		log.Fatal(err)
	}
	client = clnt
	defer client.Close(ctx)

	err = clnt.LoadCache(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = passboltv1.AddToScheme(scheme)
	if err != nil {
		log.Fatal(err)
	}

	m.Run()
}

func TestUpdateSecret(t *testing.T) {
	type args struct {
		ctx    context.Context
		clnt   *passbolt.Client
		scheme *runtime.Scheme
		pbscrt *passboltv1.PassboltSecret
		secret *corev1.Secret
	}
	tests := []struct {
		name    string
		args    args
		want    *corev1.Secret
		wantErr bool
	}{
		// docker config json
		{
			name: "success simple",
			args: args{
				ctx:    context.Background(),
				clnt:   client,
				scheme: scheme,
				pbscrt: &passboltv1.PassboltSecret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Spec: passboltv1.PassboltSecretSpec{
						SecretType: corev1.SecretTypeDockerConfigJson,
						PassboltSecretID: func() *string {
							s := "184734ea-8be3-4f5a-ba6c-5f4b3c0603e8"
							return &s
						}(),
					},
				},
				secret: &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Type: corev1.SecretTypeDockerConfigJson,
				},
			},
			want: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "passbolt.tagesspiegel.de/v1",
							Kind:               "PassboltSecret",
							Name:               "test",
							Controller:         func() *bool { b := true; return &b }(),
							BlockOwnerDeletion: func() *bool { b := true; return &b }(),
						},
					},
				},
				Type: corev1.SecretTypeDockerConfigJson,
				Data: map[string][]byte{
					corev1.DockerConfigJsonKey: []byte(`{"auths":{"https://app.example.com":{"auth":"YWRtaW46YWRtaW4="}}}`),
				},
			},
			wantErr: false,
		},
		{
			name: "secret not found",
			args: args{
				ctx:    context.Background(),
				clnt:   client,
				scheme: scheme,
				pbscrt: &passboltv1.PassboltSecret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Spec: passboltv1.PassboltSecretSpec{
						SecretType: corev1.SecretTypeDockerConfigJson,
						PassboltSecretID: func() *string {
							s := "APP_EXAMPLE_4"
							return &s
						}(),
					},
				},
				secret: &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Type: corev1.SecretTypeDockerConfigJson,
				},
			},
			want:    nil,
			wantErr: true,
		},
		// opaque
		{
			name: "secret type opaque simple",
			args: args{
				ctx:    context.Background(),
				clnt:   client,
				scheme: scheme,
				pbscrt: &passboltv1.PassboltSecret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Spec: passboltv1.PassboltSecretSpec{
						SecretType: corev1.SecretTypeOpaque,
						PassboltSecrets: map[string]passboltv1.PassboltSecretRef{
							"test": {
								ID:    "184734ea-8be3-4f5a-ba6c-5f4b3c0603e8",
								Field: passboltv1.FieldNameUsername,
							},
						},
					},
				},
				secret: &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Type: corev1.SecretTypeOpaque,
				},
			},
			want: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "passbolt.tagesspiegel.de/v1",
							Kind:               "PassboltSecret",
							Name:               "test",
							Controller:         func() *bool { b := true; return &b }(),
							BlockOwnerDeletion: func() *bool { b := true; return &b }(),
						},
					},
				},
				Type: corev1.SecretTypeOpaque,
				Data: map[string][]byte{
					"test": []byte(`admin`),
				},
			},
			wantErr: false,
		},
		{
			name: "secret type opaque templating",
			args: args{
				ctx:    context.Background(),
				clnt:   client,
				scheme: scheme,
				pbscrt: &passboltv1.PassboltSecret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Spec: passboltv1.PassboltSecretSpec{
						SecretType: corev1.SecretTypeOpaque,
						PassboltSecrets: map[string]passboltv1.PassboltSecretRef{
							"test": {
								ID:    "184734ea-8be3-4f5a-ba6c-5f4b3c0603e8",
								Value: func() *string { s := "amqp://{{ .Username }}:{{ .Password }}@{{ .URI }}/sample"; return &s }(),
							},
						},
					},
				},
				secret: &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Type: corev1.SecretTypeOpaque,
				},
			},
			want: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "passbolt.tagesspiegel.de/v1",
							Kind:               "PassboltSecret",
							Name:               "test",
							Controller:         func() *bool { b := true; return &b }(),
							BlockOwnerDeletion: func() *bool { b := true; return &b }(),
						},
					},
				},
				Type: corev1.SecretTypeOpaque,
				Data: map[string][]byte{
					"test": []byte(`amqp://admin:admin@https://app.example.com/sample`),
				},
			},
			wantErr: false,
		},
		{
			name: "secret type opaque with plain text value",
			args: args{
				ctx:    context.Background(),
				clnt:   client,
				scheme: scheme,
				pbscrt: &passboltv1.PassboltSecret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Spec: passboltv1.PassboltSecretSpec{
						SecretType: corev1.SecretTypeOpaque,
						PassboltSecrets: map[string]passboltv1.PassboltSecretRef{
							"test": {
								ID:    "184734ea-8be3-4f5a-ba6c-5f4b3c0603e8",
								Field: passboltv1.FieldNameUsername,
							},
						},
						PlainTextFields: map[string]string{
							"foo": "bar",
						},
					},
				},
				secret: &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Type: corev1.SecretTypeOpaque,
				},
			},
			want: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "passbolt.tagesspiegel.de/v1",
							Kind:               "PassboltSecret",
							Name:               "test",
							Controller:         func() *bool { b := true; return &b }(),
							BlockOwnerDeletion: func() *bool { b := true; return &b }(),
						},
					},
				},
				Type: corev1.SecretTypeOpaque,
				Data: map[string][]byte{
					"test": []byte(`admin`),
					"foo":  []byte(`bar`),
				},
			},
			wantErr: false,
		},
		{
			name: "secret field and templating not set",
			args: args{
				ctx:    context.Background(),
				clnt:   client,
				scheme: scheme,
				pbscrt: &passboltv1.PassboltSecret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Spec: passboltv1.PassboltSecretSpec{
						SecretType: corev1.SecretTypeOpaque,
						PassboltSecrets: map[string]passboltv1.PassboltSecretRef{
							"test": {
								ID: "184734ea-8be3-4f5a-ba6c-5f4b3c0603e8",
							},
						},
					},
				},
				secret: &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Type: corev1.SecretTypeOpaque,
				},
			},
			want:    nil,
			wantErr: true,
		},
		// unsupported secret type
		{
			name: "unsupported secret type SecretTypeBasicAuth",
			args: args{
				ctx:    context.Background(),
				clnt:   client,
				scheme: scheme,
				pbscrt: &passboltv1.PassboltSecret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Spec: passboltv1.PassboltSecretSpec{
						SecretType: corev1.SecretTypeBasicAuth,
						PassboltSecrets: map[string]passboltv1.PassboltSecretRef{
							"test": {
								ID:    "184734ea-8be3-4f5a-ba6c-5f4b3c0603e8",
								Field: passboltv1.FieldNameUsername,
							},
						},
					},
				},
				secret: &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Type: corev1.SecretTypeBasicAuth,
				},
			},
			want:    nil,
			wantErr: true,
		},

		{
			name: "with nil data map in secret",
			args: args{
				ctx:    context.Background(),
				clnt:   client,
				scheme: scheme,
				pbscrt: &passboltv1.PassboltSecret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Spec: passboltv1.PassboltSecretSpec{
						SecretType: corev1.SecretTypeOpaque,
						PassboltSecrets: map[string]passboltv1.PassboltSecretRef{
							"test": {
								ID:    "184734ea-8be3-4f5a-ba6c-5f4b3c0603e8",
								Field: passboltv1.FieldNameUsername,
							},
						},
					},
				},
				secret: &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Type: corev1.SecretTypeOpaque,
					Data: nil,
				},
			},
			want: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "passbolt.tagesspiegel.de/v1",
							Kind:               "PassboltSecret",
							Name:               "test",
							Controller:         func() *bool { b := true; return &b }(),
							BlockOwnerDeletion: func() *bool { b := true; return &b }(),
						},
					},
				},
				Type: corev1.SecretTypeOpaque,
				Data: map[string][]byte{
					"test": []byte(`admin`),
				},
			},
			wantErr: false,
		},
		{
			name: "with ptr to empty secret",
			args: args{
				ctx:    context.Background(),
				clnt:   client,
				scheme: scheme,
				pbscrt: &passboltv1.PassboltSecret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Spec: passboltv1.PassboltSecretSpec{
						SecretType: corev1.SecretTypeOpaque,
						PassboltSecrets: map[string]passboltv1.PassboltSecretRef{
							"test": {
								ID:    "184734ea-8be3-4f5a-ba6c-5f4b3c0603e8",
								Field: passboltv1.FieldNameUsername,
							},
						},
					},
				},
				secret: &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
				},
			},
			want: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "passbolt.tagesspiegel.de/v1",
							Kind:               "PassboltSecret",
							Name:               "test",
							Controller:         func() *bool { b := true; return &b }(),
							BlockOwnerDeletion: func() *bool { b := true; return &b }(),
						},
					},
				},
				Data: map[string][]byte{
					"test": []byte(`admin`),
				},
			},
			wantErr: false,
		},
		{
			name: "with nil secret",
			args: args{
				ctx:    context.Background(),
				clnt:   client,
				scheme: scheme,
				pbscrt: &passboltv1.PassboltSecret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Spec: passboltv1.PassboltSecretSpec{
						SecretType: corev1.SecretTypeOpaque,
						PassboltSecrets: map[string]passboltv1.PassboltSecretRef{
							"test": {
								ID:    "184734ea-8be3-4f5a-ba6c-5f4b3c0603e8",
								Field: passboltv1.FieldNameUsername,
							},
						},
					},
				},
				secret: nil,
			},
			want: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "passbolt.tagesspiegel.de/v1",
							Kind:               "PassboltSecret",
							Name:               "test",
							Controller:         func() *bool { b := true; return &b }(),
							BlockOwnerDeletion: func() *bool { b := true; return &b }(),
						},
					},
				},
				Data: map[string][]byte{
					"test": []byte(`admin`),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UpdateSecret(tt.args.ctx, tt.args.clnt, tt.args.scheme, tt.args.pbscrt, tt.args.secret)()
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateSecret() error = %v != %v", err, tt.wantErr)
				return
			}

			diff := cmp.Diff(tt.args.secret, tt.want)
			if (diff != "") != tt.wantErr {
				t.Errorf("UpdateSecret() mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}

func Test_getSecretDockerConfigJson(t *testing.T) {
	type args struct {
		secret *passbolt.PassboltSecretDefinition
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][]byte
		wantErr bool
	}{
		{
			name: "success simple",
			args: args{
				secret: &passbolt.PassboltSecretDefinition{
					Username: "test",
					Password: "test",
					URI:      "http://registry.localhost:5000",
				},
			},
			want: map[string][]byte{
				corev1.DockerConfigJsonKey: []byte(`{"auths":{"http://registry.localhost:5000":{"auth":"dGVzdDp0ZXN0"}}}`),
			},
			wantErr: false,
		},
		{
			name: "success complex",
			args: args{
				secret: &passbolt.PassboltSecretDefinition{
					Username: "984fc5ff-45f6-4663-9ca4-069293a10f06",
					Password: "d9c59199-22c7-4c8e-a064-2a777aa5c497",
					URI:      "http://registry.localhost:5000",
				},
			},
			want: map[string][]byte{
				corev1.DockerConfigJsonKey: []byte(`{"auths":{"http://registry.localhost:5000":{"auth":"OTg0ZmM1ZmYtNDVmNi00NjYzLTljYTQtMDY5MjkzYTEwZjA2OmQ5YzU5MTk5LTIyYzctNGM4ZS1hMDY0LTJhNzc3YWE1YzQ5Nw=="}}}`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getSecretDockerConfigJson(tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("getSecretDockerConfigJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Errorf("getSecretDockerConfigJson() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_getSecretTemplateValueData(t *testing.T) {
	type args struct {
		templateStr string
		secret      *passbolt.PassboltSecretDefinition
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "success simple",
			args: args{
				templateStr: "plain text",
				secret: &passbolt.PassboltSecretDefinition{
					URI:      "localhost:5672",
					Username: "guest",
					Password: "guest",
				},
			},
			want:    []byte(`plain text`),
			wantErr: false,
		},
		{
			name: "success complext",
			args: args{
				templateStr: "amqp://{{ .Username }}:{{ .Password }}@{{ .URI }}",
				secret: &passbolt.PassboltSecretDefinition{
					URI:      "localhost:5672",
					Username: "guest",
					Password: "guest",
				},
			},
			want:    []byte(`amqp://guest:guest@localhost:5672`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getSecretTemplateValueData(tt.args.templateStr, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("getSecretTemplateValueData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Errorf("getSecretTemplateValueData() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
