// Package kiwi is a Go binding for Kiwi (https://github.com/bab2min/Kiwi) project.
package kiwi

/*
#cgo LDFLAGS: -l kiwi
#include <kiwi/capi.h>
*/
import "C"

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
	// TODO: convert string to enum
	Tag string

	// Form is the actual string of this token.
	Form string
}

// TokenResult is a result for Analyze.
type TokenResult struct {
	Tokens  []TokenInfo
	WordNum int
	Score   float32
}

// Analyze returns the result of the analysis.
func (k *Kiwi) Analyze(text string, topN int, options AnalyzeOption) []TokenResult {
	kiwiResH := C.kiwi_analyze(k.handler, C.CString(text), C.int(topN), C.int(options))

	defer C.kiwi_res_close(kiwiResH)

	resSize := int(C.kiwi_res_size(kiwiResH))
	res := make([]TokenResult, resSize)

	for i := 0; i < resSize; i++ {
		tokens := make([]TokenInfo, int(C.kiwi_res_word_num(kiwiResH, C.int(i))))

		for j := 0; j < len(tokens); j++ {
			tokens[j] = TokenInfo{
				Form:     C.GoString(C.kiwi_res_form(kiwiResH, C.int(i), C.int(j))),
				Tag:      C.GoString(C.kiwi_res_tag(kiwiResH, C.int(i), C.int(j))),
				Position: int(C.kiwi_res_position(kiwiResH, C.int(i), C.int(j))),
			}
		}

		res[i] = TokenResult{
			Tokens:  tokens,
			WordNum: int(C.kiwi_res_word_num(kiwiResH, C.int(i))),
			Score:   float32(C.kiwi_res_prob(kiwiResH, C.int(i))),
		}
	}

	return res
}

// Close frees the resource allocated for Kiwi and returns the exit status.
// This must be called after New.
//
// Returns 0 if successful.
func (k *Kiwi) Close() int {
	return int(C.kiwi_close(k.handler))
}

