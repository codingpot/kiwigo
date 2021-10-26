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

func TestAddWordFail(t *testing.T) {
	kb := NewBuilder("./ModelGenerator", 1, KIWI_BUILD_INTEGRATE_ALLOMORPH)
	add := kb.AddWord("아버지가", "SKO", 0)
	assert.Equal(t, -1, add)
	assert.Equal(t, 0, kb.Close())
}

func TestAddWord(t *testing.T) {
	kb := NewBuilder("./ModelGenerator", 1, KIWI_BUILD_INTEGRATE_ALLOMORPH)
	add := kb.AddWord("아버지가", "NNG", 0)
	kiwi := kb.Build()
	res, _ := kiwi.Analyze("아버지가 방에 들어가신다", 1, KIWI_MATCH_ALL)

	expected := []TokenResult{
		{
			Tokens: []TokenInfo{
				{
					Position: 0,
					Tag:      "NNG",
					Form:     "아버지가",
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
			Score: -36.959194,
		},
	}

	assert.Equal(t, 0, add)
	assert.Equal(t, expected, res)
	assert.Equal(t, 0, kiwi.Close())
	assert.Equal(t, 0, kb.Close())
}

func TestLoadDict(t *testing.T) {
	kb := NewBuilder("./ModelGenerator", 1, KIWI_BUILD_INTEGRATE_ALLOMORPH)
	add := kb.LoadDict("example/user_dict.txt")
	err := KiwiError()
	kiwi := kb.Build()
	res := kiwi.Analyze("사용자사전 뭄바이오로사 케토톱", 1, KIWI_MATCH_ALL)

	expected := []TokenResult{
		{
			Tokens: []TokenInfo{
				{
					Position: 0,
					Tag:      "NNG",
					Form:     "아버지가",
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
			Score: -36.959194,
		},
	}

	assert.Equal(t, 0, add)
	assert.Equal(t, 0, err)
	assert.Equal(t, expected, res)
	assert.Equal(t, 0, kiwi.Close())
	assert.Equal(t, 0, kb.Close())
}

func TestAddWordFail(t *testing.T) {
	kb := NewBuilder("./ModelGenerator", 1, KIWI_BUILD_INTEGRATE_ALLOMORPH)
	add := kb.AddWord("아버지가", "SKO", 0)
	assert.Equal(t, -1, add)
	assert.Equal(t, 0, kb.Close())
}

func TestAddWord(t *testing.T) {
	kb := NewBuilder("./ModelGenerator", 1, KIWI_BUILD_INTEGRATE_ALLOMORPH)
	add := kb.AddWord("아버지가", "SKO", 0)
	kiwi := kb.Build()
	res := kiwi.Analyze("아버지가 방에 들어가신다", 1, KIWI_MATCH_ALL)

	expected := []TokenResult{
		{
			Tokens: []TokenInfo{
				{
					Position: 0,
					Tag:      "NNG",
					Form:     "아버지가",
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
			Score: -36.959194,
		},
	}

	assert.Equal(t, 0, add)
	assert.Equal(t, expected, res)
	assert.Equal(t, 0, kiwi.Close())
	assert.Equal(t, 0, kb.Close())
}
