// Package kiwi is a Go binding for Kiwi (https://github.com/bab2min/Kiwi) project.
package kiwi

/*
#cgo LDFLAGS: -l kiwi
#include <stdlib.h>
#include <string.h>
#include <stdint.h> // for uintptr_t

#include <kiwi/capi.h>

extern int KiwiReaderBridge(int lineNumber, char *buffer, void *userData);
*/
import "C"

import (
	"io"
	"runtime/cgo"
	"unsafe"

	"github.com/codingpot/kiwigo/internal"
)

// BuildOption is a bitwise OR of the KiwiBuildOption values.
type BuildOption int

const (
	KIWI_BUILD_LOAD_DEFAULT_DICT   BuildOption = C.KIWI_BUILD_LOAD_DEFAULT_DICT
	KIWI_BUILD_INTEGRATE_ALLOMORPH BuildOption = C.KIWI_BUILD_INTEGRATE_ALLOMORPH
	KIWI_BUILD_DEFAULT             BuildOption = C.KIWI_BUILD_DEFAULT
)

// AnalyzeOption is a bitwise OR of the KiwiAnalyzeOption values.
type AnalyzeOption int

const (
	KIWI_MATCH_URL                  AnalyzeOption = C.KIWI_MATCH_URL
	KIWI_MATCH_EMAIL                AnalyzeOption = C.KIWI_MATCH_EMAIL
	KIWI_MATCH_HASHTAG              AnalyzeOption = C.KIWI_MATCH_HASHTAG
	KIWI_MATCH_MENTION              AnalyzeOption = C.KIWI_MATCH_MENTION
	KIWI_MATCH_ALL                  AnalyzeOption = C.KIWI_MATCH_ALL
	KIWI_MATCH_NORMALIZE_CODA       AnalyzeOption = C.KIWI_MATCH_NORMALIZE_CODA
	KIWI_MATCH_ALL_WITH_NORMALIZING AnalyzeOption = C.KIWI_MATCH_ALL_WITH_NORMALIZING
)

// KiwiVersion returns the version of the kiwi library.
func KiwiVersion() string {
	return C.GoString(C.kiwi_version())
}

// KiwiError returns the Error messages.
func KiwiError() string {
	return C.GoString(C.kiwi_error())
}

// KiwiClearError clear error.
func KiwiClearError() {
	C.kiwi_clear_error()
}

// Kiwi is a wrapper for the kiwi C library.
type Kiwi struct {
	handler C.kiwi_h
}

// New returns a new Kiwi instance.
// Don't forget to call Close after this.
func New(modelPath string, numThread int, options BuildOption) *Kiwi {
	return &Kiwi{
		handler: C.kiwi_init(C.CString(modelPath), C.int(numThread), C.int(options)),
	}
}

// TokenInfo returns the token info for the given token(Str).
type TokenInfo struct {
	// Position is the index of this token appears in the original text.
	Position int

	// Tag represents a type of this token (e.g. VV, NNG, ...).
	Tag POSType

	// Form is the actual string of this token.
	Form string
}

// TokenResult is a result for Analyze.
type TokenResult struct {
	Tokens []TokenInfo
	Score  float32
}

// Analyze returns the result of the analysis.
func (k *Kiwi) Analyze(text string, topN int, options AnalyzeOption) ([]TokenResult, error) {
	kiwiResH := C.kiwi_analyze(k.handler, C.CString(text), C.int(topN), C.int(options))

	defer C.kiwi_res_close(kiwiResH)

	resSize := int(C.kiwi_res_size(kiwiResH))
	res := make([]TokenResult, resSize)

	for i := 0; i < resSize; i++ {
		tokens := make([]TokenInfo, int(C.kiwi_res_word_num(kiwiResH, C.int(i))))

		for j := 0; j < len(tokens); j++ {
			pos, err := ParsePOSType(C.GoString(C.kiwi_res_tag(kiwiResH, C.int(i), C.int(j))))
			if err != nil {
				return nil, err
			}
			tokens[j] = TokenInfo{
				Form:     C.GoString(C.kiwi_res_form(kiwiResH, C.int(i), C.int(j))),
				Tag:      pos,
				Position: int(C.kiwi_res_position(kiwiResH, C.int(i), C.int(j))),
			}
		}

		res[i] = TokenResult{
			Tokens: tokens,
			Score:  float32(C.kiwi_res_prob(kiwiResH, C.int(i))),
		}
	}

	return res, nil
}

// Close frees the resource allocated for Kiwi and returns the exit status.
// This must be called after New.
// Returns 0 if successful.
// Safe to call multiple times.
func (k *Kiwi) Close() int {
	if k.handler != nil {
		out := int(C.kiwi_close(k.handler))
		k.handler = nil
		return out
	}
	return 0
}

// KiwiBuilder is a wrapper for the kiwi C library.
type KiwiBuilder struct {
	handler C.kiwi_builder_h
}

// NewBuilder returns a new KiwiBuilder instance.
// Don't forget to call Close after this.
func NewBuilder(modelPath string, numThread int, options BuildOption) *KiwiBuilder {
	return &KiwiBuilder{
		handler: C.kiwi_builder_init(C.CString(modelPath), C.int(numThread), C.int(options)),
	}
}

// AddWord set custom word with word, pos, score.
func (kb *KiwiBuilder) AddWord(word string, pos POSType, score float32) int {
	return int(C.kiwi_builder_add_word(kb.handler, C.CString(word), C.CString(string(pos)), C.float(score)))
}

// LoadDict loads user dict with dict file path.
func (kb *KiwiBuilder) LoadDict(dictPath string) int {
	return int(C.kiwi_builder_load_dict(kb.handler, C.CString(dictPath)))
}

// Build creates kiwi instance with user word etc.
func (kb *KiwiBuilder) Build() *Kiwi {
	h := C.kiwi_builder_build(kb.handler)
	defer kb.Close()
	return &Kiwi{
		handler: h,
	}
}

// Close frees the resource allocated for KiwiBuilder and returns the exit status.
// This must be called after New but not need to called after Build.
// Returns 0 if successful.
// Safe to call multiple times.
func (kb *KiwiBuilder) Close() int {
	if kb.handler != nil {
		out := int(C.kiwi_builder_close(kb.handler))
		kb.handler = nil
		return out
	}
	return 0
}

// WordInfo returns the token info for the given token(Str).
type WordInfo struct {
	Form     string
	Freq     int
	POSScore float32
	Score    float32
}

//export KiwiReaderImpl
func KiwiReaderImpl(lineNumber C.int, buffer *C.char, userData unsafe.Pointer) C.int {
	scanner := cgo.Handle(userData).Value().(*internal.RewindScanner)

	if buffer == nil {
		if lineNumber == 0 {
			scanner.Rewind()
		}

		if !scanner.Scan() {
			return C.int(0)
		}

		text := scanner.Text()
		return C.int(len([]byte(text)) + 1)
	}

	textCString := C.CString(scanner.Text())
	defer C.free(unsafe.Pointer(textCString))

	C.strcpy(buffer, textCString)
	return C.int(0)
}

// ExtractWords returns the result of extract word.
func (kb *KiwiBuilder) ExtractWords(readSeeker io.ReadSeeker, minCnt int, maxWordLen int, minScore float32, posThreshold float32) ([]WordInfo, error) {
	scanner := internal.NewRewindScanner(readSeeker)
	h := cgo.NewHandle(scanner)
	defer h.Delete()

	kiwiWsH := C.kiwi_builder_extract_words(
		kb.handler,
		C.kiwi_reader_t(C.KiwiReaderBridge),
		unsafe.Pointer(h),
		C.int(minCnt), C.int(maxWordLen), C.float(minScore), C.float(posThreshold))
	defer C.kiwi_ws_close(kiwiWsH)

	resSize := int(C.kiwi_ws_size(kiwiWsH))

	if resSize < 0 {
		resSize = 0
	}

	res := make([]WordInfo, resSize)

	for i := 0; i < resSize; i++ {
		res[i] = WordInfo{
			Form:     C.GoString(C.kiwi_ws_form(kiwiWsH, C.int(i))),
			Freq:     int(C.kiwi_ws_freq(kiwiWsH, C.int(i))),
			POSScore: float32(C.kiwi_ws_pos_score(kiwiWsH, C.int(i))),
			Score:    float32(C.kiwi_ws_score(kiwiWsH, C.int(i))),
		}
	}

	return res, nil
}
