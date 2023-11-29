package nginx

import (
	"context"
	"testing"
)

func TestReloadNginx(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "test1",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Reload(context.TODO()); (err != nil) != tt.wantErr {
				t.Errorf("Reload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
