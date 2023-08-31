package cache

import (
	"context"
	"reflect"
	"testing"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
)

func TestNewInMemoryCache(t *testing.T) {
	tests := []struct {
		name string
		want Cacher
	}{
		{
			name: "TestNewInMemoryCache",
			want: NewInMemoryCache(ctrl.Log),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInMemoryCache(ctrl.Log); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInMemoryCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inMemory_Set(t *testing.T) {
	type fields struct {
		cacheData map[string]any
	}
	type args struct {
		key   string
		value any
		ttl   time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Test_inMemory_Set",
			fields: fields{cacheData: make(map[string]any)},
			args: args{
				key:   "test",
				value: "test",
				ttl:   0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &inMemory{
				cacheData: tt.fields.cacheData,
			}
			if err := c.Set(context.TODO(), tt.args.key, tt.args.value, tt.args.ttl); (err != nil) != tt.wantErr {
				t.Errorf("inMemory.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_inMemory_Get(t *testing.T) {
	type fields struct {
		cacheData map[string]any
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "Test_inMemory_Get",
			fields: fields{cacheData: map[string]any{
				"test": "test",
			}},
			args: args{
				key: "test",
			},
			want:    "test",
			wantErr: false,
		},
		{
			name: "Test_inMemory_Get cache miss",
			fields: fields{cacheData: map[string]any{
				"test": "test",
			}},
			args: args{
				key: "test2",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &inMemory{
				cacheData: tt.fields.cacheData,
			}
			got, err := c.Get(context.TODO(), tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("inMemory.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("inMemory.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inMemory_Delete(t *testing.T) {
	type fields struct {
		cacheData map[string]any
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test_inMemory_Delete",
			fields: fields{cacheData: map[string]any{
				"test": "test",
			}},
			args: args{
				ctx: context.TODO(),
				key: "test",
			},
			wantErr: false,
		},
		{
			name: "Test_inMemory_Delete cache miss",
			fields: fields{cacheData: map[string]any{
				"test": "test",
			}},
			args: args{
				ctx: context.TODO(),
				key: "test2",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &inMemory{
				cacheData: tt.fields.cacheData,
			}
			if err := c.Delete(tt.args.ctx, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("inMemory.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
