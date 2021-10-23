package kiwi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKiwiVersion(t *testing.T) {
	assert.Equal(t, KiwiVersion(), "0.10.1")
}

func TestAnalyze(t *testing.T) {
	kiwi := New("./libs/kiwi/ModelGenerator", 1, 0)
	res := kiwi.Analyze("Hello World", 1, 0)
	assert.NotEqual(t, res.Score, 0)
}
