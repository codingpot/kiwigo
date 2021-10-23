package kiwi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKiwiVersion(t *testing.T) {
	assert.Equal(t, KiwiVersion(), "0.10.1")
}

func TestAnalyze(t *testing.T) {
	kiwi := New("./libs/kiwi/ModelGenerator", 1, KIWI_BUILD_DEFAULT)
	res := kiwi.Analyze("아버지가 방에 들어가신다", 1, KIWI_MATCH_ALL)
	assert.NotEqual(t, 0, res.Score)
}
