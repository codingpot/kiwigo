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
	KIWI_BUILD_INTEGRATE_ALLOMORPH BuildOption = C.KIWI_BUILD_INTEGRATE_ALLOMORPH
	KIWI_BUILD_LOAD_DEFAULT_DICT   BuildOption = C.KIWI_BUILD_LOAD_DEFAULT_DICT
	KIWI_BUILD_LOAD_TYPO_DICT      BuildOption = C.KIWI_BUILD_LOAD_TYPO_DICT
	KIWI_BUILD_LOAD_MULTI_DICT     BuildOption = C.KIWI_BUILD_LOAD_MULTI_DICT
	KIWI_BUILD_DEFAULT             BuildOption = C.KIWI_BUILD_DEFAULT

	// Model type is a single-select field, not a bitmask.
	// Select exactly one of the following mutually exclusive model types:
	//
	//   - MODEL_TYPE_DEFAULT: default model (0x0000)
	//   - MODEL_TYPE_LARGEST: largest available model
	//   - MODEL_TYPE_KNLM: KNLM model (deprecated)
	//   - MODEL_TYPE_SBG: SBG model (deprecated)
	//   - MODEL_TYPE_CONG: CoNg model
	//   - MODEL_TYPE_CONG_GLOBAL: CoNg global model
	//
	// Do NOT combine these with bitwise OR.
	KIWI_BUILD_MODEL_TYPE_DEFAULT     BuildOption = C.KIWI_BUILD_MODEL_TYPE_DEFAULT
	KIWI_BUILD_MODEL_TYPE_LARGEST     BuildOption = C.KIWI_BUILD_MODEL_TYPE_LARGEST
	KIWI_BUILD_MODEL_TYPE_KNLM        BuildOption = C.KIWI_BUILD_MODEL_TYPE_KNLM
	KIWI_BUILD_MODEL_TYPE_SBG         BuildOption = C.KIWI_BUILD_MODEL_TYPE_SBG
	KIWI_BUILD_MODEL_TYPE_CONG        BuildOption = C.KIWI_BUILD_MODEL_TYPE_CONG
	KIWI_BUILD_MODEL_TYPE_CONG_GLOBAL BuildOption = C.KIWI_BUILD_MODEL_TYPE_CONG_GLOBAL
)

// MatchOption is a bitwise OR of the KiwiMatchOption values.
type MatchOption int

const (
	KIWI_MATCH_URL     MatchOption = C.KIWI_MATCH_URL
	KIWI_MATCH_EMAIL   MatchOption = C.KIWI_MATCH_EMAIL
	KIWI_MATCH_HASHTAG MatchOption = C.KIWI_MATCH_HASHTAG
	KIWI_MATCH_MENTION MatchOption = C.KIWI_MATCH_MENTION
	KIWI_MATCH_SERIAL  MatchOption = C.KIWI_MATCH_SERIAL
	KIWI_MATCH_EMOJI   MatchOption = C.KIWI_MATCH_EMOJI

	// OOV detection mode is a 2-bit field (bits 8-9), not independent flags.
	// Select exactly one of the following mutually exclusive modes:
	//
	//   - OOV_RULE_ONLY: rule-based scoring (default, value 0 << 8)
	//   - OOV_CHR_MODEL: character model-based scoring
	//   - OOV_CHR_FREQ_MODEL: character + frequency model scoring
	//   - OOV_CHR_FREQ_BRANCH_MODEL: character + frequency + branch model scoring
	//
	// Do NOT combine these with bitwise OR. KIWI_MATCH_OOV_MASK is for
	// internal use and should not be selected directly.
	KIWI_MATCH_OOV_RULE_ONLY             MatchOption = C.KIWI_MATCH_OOV_RULE_ONLY
	KIWI_MATCH_OOV_CHR_MODEL             MatchOption = C.KIWI_MATCH_OOV_CHR_MODEL
	KIWI_MATCH_OOV_CHR_FREQ_MODEL        MatchOption = C.KIWI_MATCH_OOV_CHR_FREQ_MODEL
	KIWI_MATCH_OOV_CHR_FREQ_BRANCH_MODEL MatchOption = C.KIWI_MATCH_OOV_CHR_FREQ_BRANCH_MODEL
	KIWI_MATCH_OOV_MASK                  MatchOption = C.KIWI_MATCH_OOV_MASK

	KIWI_MATCH_NORMALIZE_CODA   MatchOption = C.KIWI_MATCH_NORMALIZE_CODA
	KIWI_MATCH_JOIN_NOUN_PREFIX MatchOption = C.KIWI_MATCH_JOIN_NOUN_PREFIX
	KIWI_MATCH_JOIN_NOUN_SUFFIX MatchOption = C.KIWI_MATCH_JOIN_NOUN_SUFFIX
	KIWI_MATCH_JOIN_VERB_SUFFIX MatchOption = C.KIWI_MATCH_JOIN_VERB_SUFFIX
	KIWI_MATCH_JOIN_ADJ_SUFFIX  MatchOption = C.KIWI_MATCH_JOIN_ADJ_SUFFIX
	KIWI_MATCH_JOIN_ADV_SUFFIX  MatchOption = C.KIWI_MATCH_JOIN_ADV_SUFFIX
	KIWI_MATCH_JOIN_V_SUFFIX    MatchOption = C.KIWI_MATCH_JOIN_V_SUFFIX
	KIWI_MATCH_JOIN_AFFIX       MatchOption = C.KIWI_MATCH_JOIN_AFFIX
	KIWI_MATCH_SPLIT_COMPLEX    MatchOption = C.KIWI_MATCH_SPLIT_COMPLEX
	KIWI_MATCH_Z_CODA           MatchOption = C.KIWI_MATCH_Z_CODA
	KIWI_MATCH_COMPATIBLE_JAMO  MatchOption = C.KIWI_MATCH_COMPATIBLE_JAMO
	KIWI_MATCH_SPLIT_SAISIOT    MatchOption = C.KIWI_MATCH_SPLIT_SAISIOT
	KIWI_MATCH_MERGE_SAISIOT    MatchOption = C.KIWI_MATCH_MERGE_SAISIOT
	KIWI_MATCH_JOIN_PARTICLE_YO MatchOption = C.KIWI_MATCH_JOIN_PARTICLE_YO
	KIWI_MATCH_USE_OLD_SPLITTER MatchOption = C.KIWI_MATCH_USE_OLD_SPLITTER

	KIWI_MATCH_ALL                  MatchOption = C.KIWI_MATCH_ALL
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

// WithBlocklist sets the blocklist for Analyze.
func WithBlocklist(blocklist *Morphset) AnalyzeOptionFunc {
	return func(opts *AnalyzeOptions) {
		opts.Blocklist = blocklist
	}
}

// WithOpenEnding sets whether to keep the sentence open after the last morpheme.
func WithOpenEnding(openEnding bool) AnalyzeOptionFunc {
	return func(opts *AnalyzeOptions) {
		opts.OpenEnding = openEnding
	}
}

// WithAllowedDialects sets the allowed dialects for Analyze.
func WithAllowedDialects(dialects Dialect) AnalyzeOptionFunc {
	return func(opts *AnalyzeOptions) {
		opts.AllowedDialects = dialects
	}
}

// AnalyzeOptions provides configuration for the Analyze function.
type AnalyzeOptions struct {
	MatchOptions    MatchOption
	Blocklist       *Morphset
	OpenEnding      bool
	AllowedDialects Dialect
	DialectCost     float32
	TypoThreshold   float32
	TopN            int
}

// Morphset represents a set of morphemes that can be used as a blocklist.
type Morphset struct {
	handler C.kiwi_morphset_h
}

// NewMorphset creates a new morpheme set.
// The Morphset must be closed after use, before the parent Kiwi instance is closed.
func (k *Kiwi) NewMorphset() (*Morphset, error) {
	h := C.kiwi_new_morphset(k.handler)
	if h == nil {
		return nil, fmt.Errorf("failed to create morphset: %s", KiwiError())
	}
	return &Morphset{handler: h}, nil
}

// Add adds a morpheme to the set.
// tag is a POS tag such as "NNG". If tag is empty, all morphemes matching form are added.
// Returns the number of morphemes added, or an error.
func (ms *Morphset) Add(form string, tag string) (int, error) {
	cForm := C.CString(form)
	defer C.free(unsafe.Pointer(cForm))

	var cTag *C.char
	if tag != "" {
		cTag = C.CString(tag)
		defer C.free(unsafe.Pointer(cTag))
	}

	result := int(C.kiwi_morphset_add(ms.handler, cForm, cTag))
	if result < 0 {
		return 0, fmt.Errorf("failed to add morpheme: %s", KiwiError())
	}
	return result, nil
}

// Close frees the resources allocated for the Morphset.
// Must be called before the parent Kiwi instance is closed.
func (ms *Morphset) Close() {
	if ms.handler != nil {
		C.kiwi_morphset_close(ms.handler)
		ms.handler = nil
	}
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
	handler  C.kiwi_h
	dialects Dialect
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
		handler:  C.kiwi_init(C.CString(modelPath), C.int(options.numThread), C.int(options.buildOptions), C.int(options.dialects)),
		dialects: options.dialects,
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

	allowedDialects := options.AllowedDialects
	if allowedDialects == 0 {
		allowedDialects = k.dialects
	}

	var blocklistHandler C.kiwi_morphset_h
	if options.Blocklist != nil {
		blocklistHandler = options.Blocklist.handler
	}

	openEnding := 0
	if options.OpenEnding {
		openEnding = 1
	}

	cOptions := C.kiwi_analyze_option_t{
		match_options:    C.int(options.MatchOptions),
		blocklist:        blocklistHandler,
		open_ending:      C.int(openEnding),
		allowed_dialects: C.int(allowedDialects),
		dialect_cost:     C.float(options.DialectCost),
		typo_threshold:   C.float(options.TypoThreshold),
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

// Config represents the configuration for Kiwi analysis.
type Config struct {
	IntegrateAllomorph             bool
	CutOffThreshold                float32
	OovRuleScale                   float32
	OovRuleBias                    float32
	OovChrBias                     float32
	OovGlobalWeight                float32
	OovLocalWeight                 float32
	OovGlobalMinFreq               float32
	SpacePenalty                   float32
	TypoCostWeight                 float32
	MaxUnkFormSize                 uint32
	MaxUnkFormSizeFollowedByJClass uint32
	SpaceTolerance                 uint32
}

// GetGlobalConfig returns the global configuration of the Kiwi instance.
func (k *Kiwi) GetGlobalConfig() Config {
	cConfig := C.kiwi_get_global_config(k.handler)
	return Config{
		IntegrateAllomorph:             cConfig.integrate_allomorph != 0,
		CutOffThreshold:                float32(cConfig.cut_off_threshold),
		OovRuleScale:                   float32(cConfig.oov_rule_scale),
		OovRuleBias:                    float32(cConfig.oov_rule_bias),
		OovChrBias:                     float32(cConfig.oov_chr_bias),
		OovGlobalWeight:                float32(cConfig.oov_global_weight),
		OovLocalWeight:                 float32(cConfig.oov_local_weight),
		OovGlobalMinFreq:               float32(cConfig.oov_global_min_freq),
		SpacePenalty:                   float32(cConfig.space_penalty),
		TypoCostWeight:                 float32(cConfig.typo_cost_weight),
		MaxUnkFormSize:                 uint32(cConfig.max_unk_form_size),
		MaxUnkFormSizeFollowedByJClass: uint32(cConfig.max_unk_form_size_followed_by_j_class),
		SpaceTolerance:                 uint32(cConfig.space_tolerance),
	}
}

// SetGlobalConfig sets the global configuration of the Kiwi instance.
func (k *Kiwi) SetGlobalConfig(config Config) {
	cConfig := C.kiwi_config_t{
		integrate_allomorph:                   boolToCUint8(config.IntegrateAllomorph),
		cut_off_threshold:                     C.float(config.CutOffThreshold),
		oov_rule_scale:                        C.float(config.OovRuleScale),
		oov_rule_bias:                         C.float(config.OovRuleBias),
		oov_chr_bias:                          C.float(config.OovChrBias),
		oov_global_weight:                     C.float(config.OovGlobalWeight),
		oov_local_weight:                      C.float(config.OovLocalWeight),
		oov_global_min_freq:                   C.float(config.OovGlobalMinFreq),
		space_penalty:                         C.float(config.SpacePenalty),
		typo_cost_weight:                      C.float(config.TypoCostWeight),
		max_unk_form_size:                     C.uint(config.MaxUnkFormSize),
		max_unk_form_size_followed_by_j_class: C.uint(config.MaxUnkFormSizeFollowedByJClass),
		space_tolerance:                       C.uint(config.SpaceTolerance),
	}
	C.kiwi_set_global_config(k.handler, cConfig)
}

func boolToCUint8(b bool) C.uint8_t {
	if b {
		return 1
	}
	return 0
}

// OptionType represents the type of option for kiwi_set_option/kiwi_get_option.
type OptionType int

const (
	KIWI_NUM_THREADS OptionType = C.KIWI_NUM_THREADS
)

// GetOptionF returns the float value of the specified option.
// Note: As of Kiwi v0.23.2, there are no float options available.
// This function is provided for future compatibility.
func (k *Kiwi) GetOptionF(option OptionType) float32 {
	return float32(C.kiwi_get_option_f(k.handler, C.int(option)))
}

// SetOptionF sets the float value of the specified option.
// Note: As of Kiwi v0.23.2, there are no float options available.
// This function is provided for future compatibility.
func (k *Kiwi) SetOptionF(option OptionType, value float32) {
	C.kiwi_set_option_f(k.handler, C.int(option), C.float(value))
}

// GetOption returns the int value of the specified option.
func (k *Kiwi) GetOption(option OptionType) int {
	return int(C.kiwi_get_option(k.handler, C.int(option)))
}

// SetOption sets the int value of the specified option.
func (k *Kiwi) SetOption(option OptionType, value int) {
	C.kiwi_set_option(k.handler, C.int(option), C.int(value))
}
