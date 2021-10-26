package kiwi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePOSType(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    POSType
		wantErr bool
	}{
		{
			name:    "NNG is a POSType",
			arg:     "NNG",
			want:    POS_NNG,
			wantErr: false,
		},
		{
			name:    "WHATEVER is not a valid POSType",
			arg:     "WHATEVER",
			want:    POS_UNKNOWN,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePOSType(tt.arg)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
