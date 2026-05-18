// Package kiwi is a Go binding for Kiwi (https://github.com/bab2min/Kiwi) project.
package kiwi

/*
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -Wl,-rpath,/usr/local/lib

#include <stdlib.h>
#include <string.h>
#include <stdint.h> // for uintptr_t

#include <kiwi/capi.h>

extern int KiwiReaderBridge(int lineNumber, char *buffer, void *userData);
*/
import "C"

import (
	"fmt"
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

// MatchOption is a bitwise OR of the KiwiMatchOption values.
type MatchOption int

const (
	KIWI_MATCH_URL                  MatchOption = C.KIWI_MATCH_URL
	KIWI_MATCH_EMAIL                MatchOption = C.KIWI_MATCH_EMAIL
	KIWI_MATCH_HASHTAG              MatchOption = C.KIWI_MATCH_HASHTAG
	KIWI_MATCH_MENTION              MatchOption = C.KIWI_MATCH_MENTION
	KIWI_MATCH_ALL                  MatchOption = C.KIWI_MATCH_ALL
	KIWI_MATCH_NORMALIZE_CODA       MatchOption = C.KIWI_MATCH_NORMALIZE_CODA
	KIWI_MATCH_ALL_WITH_NORMALIZING MatchOption = C.KIWI_MATCH_ALL_WITH_NORMALIZING
)

// Dialect represents a dialect in the Kiwi API.
type Dialect int

const (
	DialectStandard    Dialect = 0 // KIWI_DIALECT_STANDARD
	DialectGyeonggi    Dialect = 1 << 0
	DialectChungcheong Dialect = 1 << 1
	DialectGangwon     Dialect = 1 << 2
	DialectGyeongsang  Dialect = 1 << 3
	DialectJeolla      Dialect = 1 << 4
	DialectJeju        Dialect = 1 << 5
	DialectHwanghae    Dialect = 1 << 6
	DialectHamgyeong   Dialect = 1 << 7
	DialectPyeongan    Dialect = 1 << 8
	DialectArchaic     Dialect = 1 << 9
	DialectAll         Dialect = (1<<9)*2 - 1
)

const (
	// Default values derived from Kiwi C-API defaults (include/kiwi/capi.h).
	// For detailed information on these parameters, refer to:
	// https://github.com/bab2min/Kiwi/blob/main/include/kiwi/capi.h
	DefaultDialectCost   float32 = 3.0 // Default penalty for dialect words (dialect_cost)
	DefaultTypoThreshold float32 = 2.5 // Default cost threshold for typo correction (typo_threshold)
	DefaultNumThread     int     = 0   // Default number of threads (0 means auto-detect based on CPU cores)
	DefaultTopN          int     = 1   // Default number of results to return from Analyze
)

// Option represents a configuration function for Kiwi initialization.
type Option func(*kiwiOptions)

type kiwiOptions struct {
	buildOptions BuildOption
	dialects     Dialect
	numThread    int
}

// WithBuildOption sets the BuildOption for initialization.
func WithBuildOption(options BuildOption) Option {
	return func(opts *kiwiOptions) {
		opts.buildOptions = options
	}
}

// WithDialect sets the allowed dialects for initialization.
func WithDialect(dialects Dialect) Option {
	return func(opts *kiwiOptions) {
		opts.dialects = dialects
	}
}

// WithNumThread sets the number of threads for initialization.
// A value of 0 tells Kiwi to automatically use all available CPU cores.
func WithNumThread(threads int) Option {
	return func(opts *kiwiOptions) {
		opts.numThread = threads
	}
}

// AnalyzeOptionFunc represents a configuration function for Analyze.
type AnalyzeOptionFunc func(*AnalyzeOptions)

// WithMatchOption sets the MatchOption for Analyze.
func WithMatchOption(options MatchOption) AnalyzeOptionFunc {
	return func(opts *AnalyzeOptions) {
		opts.MatchOptions = options
	}
}

// WithDialectCost sets the dialect cost for Analyze.
func WithDialectCost(cost float32) AnalyzeOptionFunc {
	return func(opts *AnalyzeOptions) {
		opts.DialectCost = cost
	}
}

// WithTypoThreshold sets the typo threshold for Analyze.
func WithTypoThreshold(threshold float32) AnalyzeOptionFunc {
	return func(opts *AnalyzeOptions) {
		opts.TypoThreshold = threshold
	}
}

// WithTopN sets the maximum number of results to return from Analyze.
func WithTopN(n int) AnalyzeOptionFunc {
	return func(opts *AnalyzeOptions) {
		opts.TopN = n
	}
}

// AnalyzeOptions provides configuration for the Analyze function.
type AnalyzeOptions struct {
	MatchOptions  MatchOption
	DialectCost   float32
	TypoThreshold float32
	TopN          int
}

// DefaultAnalyzeOptions returns the default AnalyzeOptions recommended by Kiwi.
func DefaultAnalyzeOptions() AnalyzeOptions {
	return AnalyzeOptions{
		MatchOptions:  KIWI_MATCH_ALL,
		DialectCost:   DefaultDialectCost,
		TypoThreshold: DefaultTypoThreshold,
		TopN:          DefaultTopN,
	}
}

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
func New(modelPath string, opts ...Option) *Kiwi {
	options := kiwiOptions{
		buildOptions: KIWI_BUILD_DEFAULT,
		dialects:     DialectStandard,
		numThread:    DefaultNumThread,
	}
	for _, opt := range opts {
		opt(&options)
	}

	return &Kiwi{
		handler: C.kiwi_init(C.CString(modelPath), C.int(options.numThread), C.int(options.buildOptions), C.int(options.dialects)),
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
func (k *Kiwi) Analyze(text string, opts ...AnalyzeOptionFunc) ([]TokenResult, error) {
	var (
		pretokenized C.kiwi_pretokenized_h
		cText        = C.CString(text)
	)

	options := DefaultAnalyzeOptions()
	for _, opt := range opts {
		opt(&options)
	}

	defer C.free(unsafe.Pointer(cText))

	cOptions := C.kiwi_analyze_option_t{
		match_options:  C.int(options.MatchOptions),
		dialect_cost:   C.float(options.DialectCost),
		typo_threshold: C.float(options.TypoThreshold),
	}

	kiwiResH := C.kiwi_analyze(k.handler, cText, C.int(options.TopN), cOptions, pretokenized)
	if kiwiResH == nil {
		return nil, fmt.Errorf("failed to analyze text")
	}
	defer C.kiwi_res_close(kiwiResH)

	resSize := int(C.kiwi_res_size(kiwiResH))
	if resSize < 0 {
		return nil, fmt.Errorf("invalid result size: %d", resSize)
	}

	res := make([]TokenResult, resSize)

	for i := 0; i < resSize; i++ {
		wordNum := int(C.kiwi_res_word_num(kiwiResH, C.int(i)))
		if wordNum < 0 {
			return nil, fmt.Errorf("invalid word number: %d", wordNum)
		}

		tokens := make([]TokenInfo, wordNum)

		for j := 0; j < wordNum; j++ {
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

// SplitResult returns the Sentences.
type SplitResult struct {
	Text  string
	Begin int
	End   int
}

// SplitSentence returns the line of sentences.
func (k *Kiwi) SplitSentence(text string, options MatchOption) ([]SplitResult, error) {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	kiwiSsH := C.kiwi_split_into_sents(k.handler, cText, C.int(options), nil)
	if kiwiSsH == nil {
		return nil, fmt.Errorf("failed to split sentences")
	}
	defer C.kiwi_ss_close(kiwiSsH)

	resSize := int(C.kiwi_ss_size(kiwiSsH))
	if resSize < 0 {
		return nil, fmt.Errorf("invalid result size: %d", resSize)
	}

	res := make([]SplitResult, resSize)

	for i := 0; i < resSize; i++ {
		begin := int(C.kiwi_ss_begin_position(kiwiSsH, C.int(i)))
		end := int(C.kiwi_ss_end_position(kiwiSsH, C.int(i)))

		if begin < 0 || end < begin || end > len(text) {
			return nil, fmt.Errorf("invalid position range: begin=%d, end=%d", begin, end)
		}

		res[i] = SplitResult{
			Text:  text[begin:end],
			Begin: begin,
			End:   end,
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
func NewBuilder(modelPath string, opts ...Option) *KiwiBuilder {
	options := kiwiOptions{
		buildOptions: KIWI_BUILD_DEFAULT,
		dialects:     DialectStandard,
		numThread:    DefaultNumThread,
	}
	for _, opt := range opts {
		opt(&options)
	}

	return &KiwiBuilder{
		handler: C.kiwi_builder_init(C.CString(modelPath), C.int(options.numThread), C.int(options.buildOptions), C.int(options.dialects)),
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
	var (
		typos             C.kiwi_typo_h
		typoCostThreshold = C.float(1.0)
	)

	h := C.kiwi_builder_build(kb.handler, typos, typoCostThreshold)
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
		C.int(minCnt), C.int(maxWordLen), C.float(minScore), C.float(posThreshold),
	)
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
