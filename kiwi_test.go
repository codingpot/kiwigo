package kiwi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKiwiVersion(t *testing.T) {
	assert.Equal(t, KiwiVersion(), "0.10.2")
}

func TestAnalyze(t *testing.T) {
	kiwi := New("./ModelGenerator", 1, KIWI_BUILD_DEFAULT)
	res := kiwi.Analyze("아버지가 방에 들어가신다", 1, KIWI_MATCH_ALL)

	expected := []TokenResult{
		{
			Tokens: []TokenInfo{
				{
					Position: 0,
					Tag:      "NNG",
					Form:     "아버지",
				},
				{
					Position: 3,
					Tag:      "JKS",
					Form:     "가",
				},
				{
					Position: 5,
					Tag:      "NNG",
					Form:     "방",
				},
				{
					Position: 6,
					Tag:      "JKB",
					Form:     "에",
				},
				{
					Position: 8,
					Tag:      "VV",
					Form:     "들어가",
				},
				{
					Position: 11,
					Tag:      "EP",
					Form:     "시",
				},
				{
					Position: 12,
					Tag:      "EF",
					Form:     "ᆫ다",
				},
			},
			Score: -38.967132568359375,
		},
	}

	assert.Equal(t, expected, res)
	assert.Equal(t, 0, kiwi.Close())
}
