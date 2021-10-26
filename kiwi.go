// Package kiwi is a Go binding for Kiwi (https://github.com/bab2min/Kiwi) project.
package kiwi

/*
#cgo LDFLAGS: -l kiwi
#include <kiwi/capi.h>
*/
import "C"
import "fmt"

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

type POSType string

const (
	POS_UNKNOWN POSType = "UNKNOWN"

	POS_NNG POSType = "NNG"
	POS_NNP POSType = "NNP"
	POS_NNB POSType = "NNB"

	POS_VV POSType = "VV"
	POS_VA POSType = "VA"

	POS_MAG POSType = "MAG"

	POS_NR POSType = "NR"
	POS_NP POSType = "NP"

	POS_VX POSType = "VX"

	POS_MM  POSType = "MM"
	POS_MAJ POSType = "MAJ"

	POS_IC POSType = "IC"

	POS_XPN POSType = "XPN"
	POS_XSN POSType = "XSN"
	POS_XSV POSType = "XSV"
	POS_XSA POSType = "XSA"
	POS_XR  POSType = "XR"

	POS_VCP POSType = "VCP"
	POS_VCN POSType = "VCN"

	POS_SF POSType = "SF"
	POS_SP POSType = "SP"
	POS_SS POSType = "SS"
	POS_SE POSType = "SE"
	POS_SO POSType = "SO"
	POS_SW POSType = "SW"

	POS_SL POSType = "SL"
	POS_SH POSType = "SH"
	POS_SN POSType = "SN"

	POS_W_URL     POSType = "W_URL"
	POS_W_EMAIL   POSType = "W_EMAIL"
	POS_W_MENTION POSType = "W_MENTION"
	POS_W_HASHTAG POSType = "W_HASHTAG"

	POS_JKS POSType = "JKS"
	POS_JKC POSType = "JKC"
	POS_JKG POSType = "JKG"
	POS_JKO POSType = "JKO"
	POS_JKB POSType = "JKB"
	POS_JKV POSType = "JKV"
	POS_JKQ POSType = "JKQ"
	POS_JX  POSType = "JX"
	POS_JC  POSType = "JC"

	POS_EP  POSType = "EP"
	POS_EF  POSType = "EF"
	POS_EC  POSType = "EC"
	POS_ETN POSType = "ETN"
	POS_ETM POSType = "ETM"

	POS_V POSType = "V"

	POS_MAX POSType = "MAX"
)

// ParsePOSType return POS Tag for resault
func ParsePOSType(t string) (POSType, error) {
	switch t {
	case "NNG":
		return POS_NNG, nil
	case "NNP":
		return POS_NNP, nil
	case "NNB":
		return POS_NNB, nil

	case "VV":
		return POS_VV, nil
	case "VA":
		return POS_VA, nil

	case "MAG":
		return POS_MAG, nil

	case "NR":
		return POS_NR, nil
	case "NP":
		return POS_NP, nil

	case "VX":
		return POS_VX, nil

	case "MM":
		return POS_MM, nil
	case "MAJ":
		return POS_MAJ, nil

	case "IC":
		return POS_IC, nil

	case "XPN":
		return POS_XPN, nil
	case "XSN":
		return POS_XSN, nil
	case "XSV":
		return POS_XSV, nil
	case "XSA":
		return POS_XSA, nil
	case "XR":
		return POS_XR, nil

	case "VCP":
		return POS_VCP, nil
	case "VCN":
		return POS_VCN, nil

	case "SF":
		return POS_SF, nil
	case "SP":
		return POS_SP, nil
	case "SS":
		return POS_SS, nil
	case "SE":
		return POS_SE, nil
	case "SO":
		return POS_SO, nil
	case "SW":
		return POS_SW, nil

	case "SL":
		return POS_SL, nil
	case "SH":
		return POS_SH, nil
	case "SN":
		return POS_SN, nil

	case "W_URL":
		return POS_W_URL, nil
	case "W_EMAIL":
		return POS_W_EMAIL, nil
	case "W_MENTION":
		return POS_W_MENTION, nil
	case "W_HASHTAG":
		return POS_W_HASHTAG, nil

	case "JKS":
		return POS_JKS, nil
	case "JKC":
		return POS_JKC, nil
	case "JKG":
		return POS_JKG, nil
	case "JKO":
		return POS_JKO, nil
	case "JKB":
		return POS_JKB, nil
	case "JKV":
		return POS_JKV, nil
	case "JKQ":
		return POS_JKQ, nil
	case "JX":
		return POS_JX, nil
	case "JC":
		return POS_JC, nil

	case "EP":
		return POS_EP, nil
	case "EF":
		return POS_EF, nil
	case "EC":
		return POS_EC, nil
	case "ETN":
		return POS_ETN, nil
	case "ETM":
		return POS_ETM, nil

	case "V":
		return POS_V, nil

	case "MAX":
		return POS_MAX, nil

	default:
		return POS_UNKNOWN, fmt.Errorf("POS type parse err. input type: %s", t)
	}
}

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
	Tag POSType

	// Form is the actual string of this token.
	Form string
}

// TokenResult is a result for Analyze.
type TokenResult struct {
	Tokens []TokenInfo
	Score float32
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
			Score: float32(C.kiwi_res_prob(kiwiResH, C.int(i))),
		}
	}

	return res, nil
}

// Close frees the resource allocated for Kiwi and returns the exit status.
// This must be called after New.
//
// Returns 0 if successful.
func (k *Kiwi) Close() int {
	return int(C.kiwi_close(k.handler))
}
