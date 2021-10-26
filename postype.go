package kiwi

import "fmt"

type POSType string

// NOTE: update isValid function when adding a new POSType.
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

func (p POSType) isValid() bool {
	switch p {
	case POS_UNKNOWN,
		POS_NNG,
		POS_NNP,
		POS_NNB,
		POS_VV,
		POS_VA,
		POS_MAG,
		POS_NR,
		POS_NP,
		POS_VX,
		POS_MM,
		POS_MAJ,
		POS_IC,
		POS_XPN,
		POS_XSN,
		POS_XSV,
		POS_XSA,
		POS_XR,
		POS_VCP,
		POS_VCN,
		POS_SF,
		POS_SP,
		POS_SS,
		POS_SE,
		POS_SO,
		POS_SW,
		POS_SL,
		POS_SH,
		POS_SN,
		POS_W_URL,
		POS_W_EMAIL,
		POS_W_MENTION,
		POS_W_HASHTAG,
		POS_JKS,
		POS_JKC,
		POS_JKG,
		POS_JKO,
		POS_JKB,
		POS_JKV,
		POS_JKQ,
		POS_JX,
		POS_JC,
		POS_EP,
		POS_EF,
		POS_EC,
		POS_ETN,
		POS_ETM,
		POS_V,
		POS_MAX:
		return true
	default:
		return false
	}
}

// ParsePOSType return POS Tag for resault
func ParsePOSType(t string) (POSType, error) {
	pos := POSType(t)
	if !pos.isValid() {
		return POS_UNKNOWN, fmt.Errorf("POS type parse err. input type: %s", t)
	}
	return pos, nil
}
