package kiwi

import (
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

// floatComparer returns a cmp.Option for floating-point comparisons with tolerance.
func floatComparer() cmp.Option {
	return cmpopts.EquateApprox(0, 1e-5)
}

func TestKiwiVersion(t *testing.T) {
	assert.Equal(t, KiwiVersion(), "0.21.0")
}

func TestAnalyze(t *testing.T) {
	kiwi := New("./base", 1, KIWI_BUILD_DEFAULT)
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
					Position: 11,
					Tag:      POS_EF,
					Form:     "ᆫ다",
				},
			},
			Score: -34.55623,
		},
	}

	if diff := cmp.Diff(expected, res, floatComparer()); diff != "" {
		t.Errorf("Analyze result mismatch (-want +got):\n%s", diff)
	}
	assert.Equal(t, 0, kiwi.Close())
}

func TestSplitSentence(t *testing.T) {
	kiwi := New("./base", 1, KIWI_BUILD_DEFAULT)
	res, _ := kiwi.SplitSentence("여러 문장으로 구성된 텍스트네 이걸 분리해줘", KIWI_MATCH_ALL)

	expected := []SplitResult{
		{
			Text:  "여러 문장으로 구성된 텍스트네",
			Begin: 0,
			End:   42,
		},
		{
			Text:  "이걸 분리해줘",
			Begin: 43,
			End:   62,
		},
	}

	assert.Equal(t, expected, res)
	assert.Equal(t, 0, kiwi.Close())
}

func TestAddWordFail(t *testing.T) {
	kb := NewBuilder("./base", 1, KIWI_BUILD_INTEGRATE_ALLOMORPH)
	add := kb.AddWord("아버지가", "SKO", 0)
	assert.Equal(t, -1, add)
	assert.Equal(t, 0, kb.Close())

	KiwiClearError()
}

func TestAddWord(t *testing.T) {
	kb := NewBuilder("./base", 1, KIWI_BUILD_INTEGRATE_ALLOMORPH)
	add := kb.AddWord("아버지가", "NNG", 0)

	assert.Equal(t, 0, add)

	kiwi := kb.Build()
	res, _ := kiwi.Analyze("아버지가 방에 들어가신다", 1, KIWI_MATCH_ALL)

	// kb should have been closed.
	assert.Equal(t, 0, kb.Close())

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
					Position: 11,
					Tag:      "EF",
					Form:     "ᆫ다",
				},
			},
			Score: -32.80881,
		},
	}

	if diff := cmp.Diff(expected, res, floatComparer()); diff != "" {
		t.Errorf("AddWord result mismatch (-want +got):\n%s", diff)
	}
	assert.Equal(t, 0, kiwi.Close())
}

func TestLoadDict(t *testing.T) {
	kb := NewBuilder("./base", 1, KIWI_BUILD_INTEGRATE_ALLOMORPH)
	add := kb.LoadDict("./example/user_dict.tsv")

	assert.Equal(t, 1, add)

	err := KiwiError()

	assert.Equal(t, "", err)

	kiwi := kb.Build()

	// kb should have been closed already.
	assert.Equal(t, 0, kb.Close())

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
					Position: 11,
					Tag:      "EF",
					Form:     "ᆫ다",
				},
			},
			Score: -32.80881,
		},
	}

	if diff := cmp.Diff(expected, res, floatComparer()); diff != "" {
		t.Errorf("LoadDict result mismatch (-want +got):\n%s", diff)
	}
	assert.Equal(t, 0, kiwi.Close())
}

func TestLoadDict2(t *testing.T) {
	kb := NewBuilder("./base", 1, KIWI_BUILD_INTEGRATE_ALLOMORPH)
	add := kb.LoadDict("./example/user_dict2.tsv")

	assert.Equal(t, 3, add)

	err := KiwiError()

	assert.Equal(t, "", err)

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
					Form:     "방에",
				},
				{
					Position: 8,
					Tag:      "NNG",
					Form:     "들어가신다",
				},
			},
			Score: -12.538677,
		},
	}

	if diff := cmp.Diff(expected, res, floatComparer()); diff != "" {
		t.Errorf("LoadDict2 result mismatch (-want +got):\n%s", diff)
	}
	assert.Equal(t, 0, kiwi.Close())
}

func TestExtractWord(t *testing.T) {
	kb := NewBuilder("./base", 1, KIWI_BUILD_DEFAULT)
	rs := strings.NewReader(`2008년에는 애국가의 작곡자 안익태가 1930년대에 독일 유학 기간 중 친일 활동을 했다는 사실이 밝혀졌다. 이후 안익태가 나치 독일 하의
베를린에서 만주국 10주년 건국 기념음악회를 지휘하는 동영상까지 발굴되어 관련 학계나 사회에 큰 충격을 주었다. 안익태가 친일 행적을 한 바
있다는 빼도박도 못할 증거가 나왔으니까. 영상물의 '만주환상곡'에는 우리가 현재 알고있는 '한국환상곡'의 두 선율("무궁화 삼천리 나의 사랑아,
영광의 태극기 길이 빛나라", "화려한 강산 한반도, 나의 사랑 한반도 너희 뿐일세")에 거의 그 모습 그대로 나타난다. # 이로 인해 애국가의
본이 되는 한국환상곡은 사실 만주국 창립 기념을 위한 만주환상곡에서 따왔다는 주장이 제기되었고, 일각에서는 현재의 애국가는 가사(공식적으론 작사
미상이나 가장 유력한 윤치호 작사를 인정 못해서 그렇다는게 중론), 곡(안익태 작곡) 모두 친일파의 산물이라고 주장하며 국가 재제정 운동을
벌이기도 하였다. 만약 애국가 악곡의 일부를 만주환상곡에서 가져왔다는 주장이 사실이라면 이 문제가 언젠가 공론화가 된다면 국가 교체가 이루어질
가능성이 크다. 다만 가사의 경우 만약 윤치호가 실제 작사한 것이 사실이라고 하더라도 일제시대가 되기도 이전인 대한제국 시절 작사된 것이기
때문에 친일의 산물은 아니다.`)
	wordInfos, _ := kb.ExtractWords(rs, 3 /*=minCnt*/, 3 /*=maxWordLen*/, 0.0 /*=minScore*/, -3.0 /*=posThreshold*/)
	expected := []WordInfo{
		{
			Form:     "안익",
			Freq:     3,
			POSScore: -1.92593,
			Score:    0,
		},
		{
			Form:     "익태",
			Freq:     4,
			POSScore: -0.23702252,
			Score:    0,
		},
	}
	if diff := cmp.Diff(expected, wordInfos, floatComparer()); diff != "" {
		t.Errorf("ExtractWord result mismatch (-want +got):\n%s", diff)
	}
	assert.Equal(t, 0, kb.Close())
}

func TestExtractWordwithFile(t *testing.T) {
	kb := NewBuilder("./base", 1, KIWI_BUILD_DEFAULT) // Use single thread for deterministic results
	file, _ := os.Open("./example/test.txt")

	wordInfos, _ := kb.ExtractWords(file, 10 /*=minCnt*/, 5 /*=maxWordLen*/, 0.0 /*=minScore*/, -25.0 /*=posThreshold*/)
	expectedWordInfo := WordInfo{
		Form: "무위원", Freq: 17, POSScore: -1.7342134, Score: 0.69981515,
	}
	if diff := cmp.Diff(expectedWordInfo, wordInfos[0], floatComparer()); diff != "" {
		t.Errorf("ExtractWordwithFile result mismatch (-want +got):\n%s", diff)
	}
	assert.Equal(t, 0, kb.Close())
}
