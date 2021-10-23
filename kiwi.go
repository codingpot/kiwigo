package kiwi

/*
#cgo LDFLAGS: -L libs/kiwi/build -l kiwi
#include "libs/kiwi/include/kiwi/capi.h"
*/
import "C"

// KiwiBuildOption is a bitwise OR of the KiwiBuildOption values.
type KiwiBuildOption int

const (
	KIWI_BUILD_LOAD_DEFAULT_DICT   KiwiBuildOption = 1
	KIWI_BUILD_INTEGRATE_ALLOMORPH KiwiBuildOption = 2
	KIWI_BUILD_DEFAULT             KiwiBuildOption = 3
)

// KiwiAnalyzeOption is a bitwise OR of the KiwiAnalyzeOption values.
type KiwiAnalyzeOption int

const (
	KIWI_MATCH_URL     KiwiAnalyzeOption = 1
	KIWI_MATCH_EMAIL   KiwiAnalyzeOption = 2
	KIWI_MATCH_HASHTAG KiwiAnalyzeOption = 4
	KIWI_MATCH_MENTION KiwiAnalyzeOption = 8
	KIWI_MATCH_ALL     KiwiAnalyzeOption = 15
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
func New(modelPath string, numThread int, options KiwiBuildOption) *Kiwi {
	return &Kiwi{
		handler: C.kiwi_init(C.CString(modelPath), C.int(numThread), C.int(options)),
	}
}

// TokenInfo returns the token info for the given token(Str).
type TokenInfo struct {
	Str         string
	Position    int
	Length      int
	WorPosition int
}

// KiwiTokenResult is a result for Analyze.
type KiwiTokenResult struct {
	Tokens []TokenInfo
	Score  float64
}

// Analyze returns the result of the analysis.
func (k *Kiwi) Analyze(text string, topN int, options KiwiAnalyzeOption) KiwiTokenResult {
	kiwiResH := C.kiwi_analyze(k.handler, C.CString(text), C.int(topN), C.int(options))

	score := float64(C.kiwi_res_prob(kiwiResH, 0))

	return KiwiTokenResult{
		Score: score,
	}
}
