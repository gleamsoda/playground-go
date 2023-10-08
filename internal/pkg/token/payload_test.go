package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPayload_Valid(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name    string
		payload *Payload
		wantErr error
	}{
		{
			name: "Valid payload",
			payload: &Payload{
				ExpiredAt: now.Add(time.Hour),
			},
			wantErr: nil,
		},
		{
			name: "Expired payload",
			payload: &Payload{
				ExpiredAt: now.Add(-time.Hour),
			},
			wantErr: ErrExpiredToken,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.payload.Valid()
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
