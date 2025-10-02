package internal

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRewindScanner_Rewind(t *testing.T) {
	reader := strings.NewReader(`안녕하세요
코딩 냄비입니다.`)
	scanner := NewRewindScanner(reader)

	scanner.Scan()
	assert.Equal(t, "안녕하세요", scanner.Text())

	err := scanner.Rewind()
	assert.NoError(t, err)
	scanner.Scan()
	assert.Equal(t, "안녕하세요", scanner.Text())
}
