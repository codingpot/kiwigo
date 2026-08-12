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
		// kiwi_res_tag is implemented with tagRToString, so the -R variants
		// reach Go for regular conjugations. See #39.
		{
			name:    "VV-R is a POSType",
			arg:     "VV-R",
			want:    POS_VV_R,
			wantErr: false,
		},
		{
			name:    "XSA-R is a POSType",
			arg:     "XSA-R",
			want:    POS_XSA_R,
			wantErr: false,
		},
		{
			name:    "VV-I is a POSType",
			arg:     "VV-I",
			want:    POS_VV_I,
			wantErr: false,
		},
		{
			name:    "@ is a POSType",
			arg:     "@",
			want:    POS_PA,
			wantErr: false,
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

// TestAnalyzeRegularConjugation guards against #39. Kiwi tags a regular verb
// whose stem ends in ㄷ/ㅂ/ㅅ as VV-R rather than VV, so dropping the -R
// constants makes Analyze fail on ordinary sentences.
func TestAnalyzeRegularConjugation(t *testing.T) {
	kiwi, err := New("./base", WithNumThread(1))
	assert.NoError(t, err)

	for _, sentence := range []string{"편지를 받았다", "나는 공을 잡았다", "그는 크게 웃었다"} {
		t.Run(sentence, func(t *testing.T) {
			res, err := kiwi.Analyze(sentence)
			assert.NoError(t, err)

			var tags []POSType
			for _, result := range res {
				for _, token := range result.Tokens {
					tags = append(tags, token.Tag)
				}
			}
			assert.Contains(t, tags, POS_VV_R)
		})
	}
}
