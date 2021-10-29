package internal

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestRewindScanner_Rewind(t *testing.T) {
	reader := strings.NewReader(`안녕하세요
코딩 냄비입니다.`)
	scanner := NewRewindScanner(reader)

	scanner.Scan()
	assert.Equal(t, "안녕하세요", scanner.Text())

	scanner.Rewind()
	scanner.Scan()
	assert.Equal(t, "안녕하세요", scanner.Text())
}
