package kiwi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKiwiVersion(t *testing.T) {
	assert.Equal(t, KiwiVersion(), "0.10.2")
}

func TestParsePOSType(t *testing.T) {
	res, _ := ParsePOSType("NNG")
	assert.Equal(t, res, POS_NNG)

	if _, err := ParsePOSType("NNK"); err != nil {
		assert.EqualError(t, err, "POS type parse err. input type: NNK")
	}
}

func TestAnalyze(t *testing.T) {
	kiwi := New("./ModelGenerator", 1, KIWI_BUILD_DEFAULT)
	res, _ := kiwi.Analyze("아버지가 방에 들어가신다", 1, KIWI_MATCH_ALL)

	expected := []TokenResult{
		{
			Tokens: []TokenInfo{
				{
					Position: 0,
					Tag:      POS_NNG,
					Form:     "아버지",
				},
				{
					Position: 3,
					Tag:      POS_JKS,
					Form:     "가",
				},
				{
					Position: 5,
					Tag:      POS_NNG,
					Form:     "방",
				},
				{
					Position: 6,
					Tag:      POS_JKB,
					Form:     "에",
				},
				{
					Position: 8,
					Tag:      POS_VV,
					Form:     "들어가",
				},
				{
					Position: 11,
					Tag:      POS_EP,
					Form:     "시",
				},
				{
					Position: 12,
					Tag:      POS_EF,
					Form:     "ᆫ다",
				},
			},
			Score: -38.967132568359375,
		},
	}

	assert.Equal(t, expected, res)
	assert.Equal(t, 0, kiwi.Close())
}
