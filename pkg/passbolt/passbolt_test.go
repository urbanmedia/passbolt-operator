package passbolt

import (
	"context"
	"testing"

	passboltv1alpha1 "github.com/urbanmedia/passbolt-operator/api/v1alpha1"
)

const (
	// passboltURL is the URL of the passbolt server.
	passboltURL = "http://localhost:8088"
	// passboltUsername is the username of the passbolt user.
	passboltUsername = `-----BEGIN PGP PRIVATE KEY BLOCK-----

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
=Z2i3
-----END PGP PRIVATE KEY BLOCK-----`
	// passboltPassword is the password of the passbolt user.
	passboltPassword = "TestTest123!"
)

func TestNewClient(t *testing.T) {
	type args struct {
		ctx      context.Context
		url      string
		username string
		password string
	}
	tests := []struct {
		name          string
		args          args
		wantNilClient bool
		wantErr       bool
	}{
		{
			name: "happy path",
			args: args{
				ctx:      context.Background(),
				url:      passboltURL,
				username: passboltUsername,
				password: passboltPassword,
			},
			wantNilClient: false,
			wantErr:       false,
		},
		{
			name: "openpgp error",
			args: args{
				ctx: context.Background(),
				url: "",
				username: `-----BEGIN PGP PRIVATE KEY BLOCK-----
				-----END PGP PRIVATE KEY BLOCK-----`,
				password: "a",
			},
			// we expect an error so both wantNilClient and wantErr are true
			wantNilClient: true,
			wantErr:       true,
		},
		{
			name: "login error",
			args: args{
				ctx:      context.Background(),
				url:      "http://passbolt.local",
				username: passboltUsername,
				password: passboltPassword,
			},
			// we expect an error so both wantNilClient and wantErr are true
			wantNilClient: true,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.ctx, tt.args.url, tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil {
				// we need to close the client to avoid a leak
				defer got.Close(tt.args.ctx)
			}

			if (got == nil) != tt.wantNilClient {
				t.Errorf("NewClient() = %v, want nil client", got)
			}
		})
	}
}

func TestClient_LoadCache(t *testing.T) {
	type fields struct {
		client *Client
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "happy path",
			fields: fields{
				client: func() *Client {
					clnt, err := NewClient(
						context.Background(),
						passboltURL,
						passboltUsername,
						passboltPassword)
					if err != nil {
						t.Errorf("failed to create passbolt client: %v", err)
						return nil
					}
					return clnt
				}(),
			},
			args: args{
				ctx: context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields.client
			defer tt.fields.client.Close(context.Background())

			if err := c.LoadCache(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Client.LoadCache() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Close(t *testing.T) {
	type fields struct {
		client *Client
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "happy path",
			fields: fields{
				client: func() *Client {
					clnt, err := NewClient(
						context.Background(),
						passboltURL,
						passboltUsername,
						passboltPassword)
					if err != nil {
						t.Errorf("failed to create passbolt client: %v", err)
						return nil
					}
					return clnt
				}(),
			},
			args: args{
				ctx: context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields.client

			if err := c.Close(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Client.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_GetSecret(t *testing.T) {
	type fields struct {
		client *Client
	}
	type args struct {
		ctx   context.Context
		name  string
		field passboltv1alpha1.FieldName
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "happy path username",
			fields: fields{
				client: func() *Client {
					clnt, err := NewClient(
						context.Background(),
						passboltURL,
						passboltUsername,
						passboltPassword)
					if err != nil {
						t.Errorf("failed to create passbolt client: %v", err)
						return nil
					}
					return clnt
				}(),
			},
			args: args{
				ctx:   context.Background(),
				name:  "APP_EXAMPLE",
				field: passboltv1alpha1.FieldNameUsername,
			},
			want:    "admin",
			wantErr: false,
		},
		{
			name: "happy path password",
			fields: fields{
				client: func() *Client {
					clnt, err := NewClient(
						context.Background(),
						passboltURL,
						passboltUsername,
						passboltPassword)
					if err != nil {
						t.Errorf("failed to create passbolt client: %v", err)
						return nil
					}
					return clnt
				}(),
			},
			args: args{
				ctx:   context.Background(),
				name:  "APP_EXAMPLE",
				field: passboltv1alpha1.FieldNamePassword,
			},
			want:    "admin",
			wantErr: false,
		},
		{
			name: "happy path uri",
			fields: fields{
				client: func() *Client {
					clnt, err := NewClient(
						context.Background(),
						passboltURL,
						passboltUsername,
						passboltPassword)
					if err != nil {
						t.Errorf("failed to create passbolt client: %v", err)
						return nil
					}
					return clnt
				}(),
			},
			args: args{
				ctx:   context.Background(),
				name:  "APP_EXAMPLE",
				field: passboltv1alpha1.FieldNameUri,
			},
			want:    "https://app.example.com",
			wantErr: false,
		},
		{
			name: "key not in cache",
			fields: fields{
				client: func() *Client {
					clnt, err := NewClient(
						context.Background(),
						passboltURL,
						passboltUsername,
						passboltPassword)
					if err != nil {
						t.Errorf("failed to create passbolt client: %v", err)
						return nil
					}
					return clnt
				}(),
			},
			args: args{
				ctx:  context.Background(),
				name: "example",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields.client

			// load cache
			if err := c.LoadCache(tt.args.ctx); err != nil {
				t.Errorf("failed to load cache: %v", err)
				return
			}

			got, err := c.GetSecret(tt.args.ctx, tt.args.name, tt.args.field)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetSecret() error = %v, wantErr %v\nCache data:\n%+v", err, tt.wantErr, tt.fields.client.secretCache)
				return
			}
			if got != tt.want {
				t.Errorf("Client.GetSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}
