package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_collector_GetSystemData(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr error
	}{
		{"simple test", "success", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New("../adapter/sms/test_sms.data")
			got := c.GetSystemData()

			assert.Equal(t, tt.want, got)
			//assert.Equal(t, tt.wantErr, err)
		})
	}
}
