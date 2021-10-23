package kiwi

/*
#cgo LDFLAGS: -L libs/kiwi/build -l kiwi
#include "libs/kiwi/include/kiwi/capi.h"
*/
import "C"

func KiwiVersion() string {
	return C.GoString(C.kiwi_version())
}

type Kiwi struct {
	handler C.kiwi_h
}

func New(modelPath string, numThread int, options int) *Kiwi {
	return &Kiwi{
		handler: C.kiwi_init(C.CString(modelPath), C.int(numThread), C.int(options)),
	}
}

type TokenInfo struct {
	Str         string
	Position    int
	Length      int
	WorPosition int
}

type KiwiTokenResult struct {
	Tokens []TokenInfo
	Score  float64
}

func (k *Kiwi) Analyze(text string, topN int, options int) KiwiTokenResult {
	kiwiResH := C.kiwi_analyze(k.handler, C.CString(text), C.int(topN), C.int(options))

	score := float64(C.kiwi_res_prob(kiwiResH, 0))

	return KiwiTokenResult{
		Score: score,
	}
}
