package kiwi

import "fmt"

type POSType string

// NOTE: update isValid function when adding a new POSType.
// POS type definitions from Kiwi C++ library:
// - Type https://github.com/bab2min/Kiwi/blob/18b4062f98cf3beaedac149fc511a2708891f71d/include/kiwi/Types.h#L188-L219
// - String https://github.com/bab2min/Kiwi/blob/18b4062f98cf3beaedac149fc511a2708891f71d/src/Utils.cpp#L300
const (
	POS_UNKNOWN POSType = "UN"

	POS_NNG POSType = "NNG"
	POS_NNP POSType = "NNP"
	POS_NNB POSType = "NNB"
	POS_NR  POSType = "NR"
	POS_NP  POSType = "NP"

	POS_VV  POSType = "VV"
	POS_VA  POSType = "VA"
	POS_VX  POSType = "VX"
	POS_VCP POSType = "VCP"
	POS_VCN POSType = "VCN"

	POS_MM POSType = "MM"

	POS_MAG POSType = "MAG"
	POS_MAJ POSType = "MAJ"

	POS_IC POSType = "IC"

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

	POS_XPN POSType = "XPN"

	POS_XSN POSType = "XSN"
	POS_XSV POSType = "XSV"
	POS_XSA POSType = "XSA"
	POS_XSM POSType = "XSM"

	POS_XR POSType = "XR"

	POS_SF  POSType = "SF"
	POS_SP  POSType = "SP"
	POS_SS  POSType = "SS"
	POS_SSO POSType = "SSO"
	POS_SSC POSType = "SSC"
	POS_SE  POSType = "SE"
	POS_SO  POSType = "SO"
	POS_SW  POSType = "SW"
	POS_SL  POSType = "SL"
	POS_SH  POSType = "SH"
	POS_SN  POSType = "SN"
	POS_SB  POSType = "SB"

	POS_W_URL     POSType = "W_URL"
	POS_W_EMAIL   POSType = "W_EMAIL"
	POS_W_HASHTAG POSType = "W_HASHTAG"
	POS_W_MENTION POSType = "W_MENTION"
	POS_W_SERIAL  POSType = "W_SERIAL"
	POS_W_EMOJI   POSType = "W_EMOJI"

	POS_Z_CODA POSType = "Z_CODA"
	POS_Z_SIOT POSType = "Z_SIOT"

	POS_USER_0 POSType = "USER0"
	POS_USER_1 POSType = "USER1"
	POS_USER_2 POSType = "USER2"
	POS_USER_3 POSType = "USER3"
	POS_USER_4 POSType = "USER4"

	POS_VV_I  POSType = "VV-I"
	POS_VA_I  POSType = "VA-I"
	POS_VX_I  POSType = "VX-I"
	POS_XSA_I POSType = "XSA-I"

	POS_VV_R  POSType = "VV-R"
	POS_VA_R  POSType = "VA-R"
	POS_VX_R  POSType = "VX-R"
	POS_XSA_R POSType = "XSA-R"

	POS_V POSType = "V"
)

func (p POSType) isValid() bool {
	switch p {
	case
		POS_UNKNOWN,
		POS_NNG,
		POS_NNP,
		POS_NNB,
		POS_NR,
		POS_NP,
		POS_VV,
		POS_VA,
		POS_VX,
		POS_VCP,
		POS_VCN,
		POS_MM,
		POS_MAG,
		POS_MAJ,
		POS_IC,
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
		POS_XPN,
		POS_XSN,
		POS_XSV,
		POS_XSA,
		POS_XSM,
		POS_XR,
		POS_SF,
		POS_SP,
		POS_SS,
		POS_SSO,
		POS_SSC,
		POS_SE,
		POS_SO,
		POS_SW,
		POS_SL,
		POS_SH,
		POS_SN,
		POS_SB,
		POS_W_URL,
		POS_W_EMAIL,
		POS_W_HASHTAG,
		POS_W_MENTION,
		POS_W_SERIAL,
		POS_W_EMOJI,
		POS_Z_CODA,
		POS_Z_SIOT,
		POS_USER_0,
		POS_USER_1,
		POS_USER_2,
		POS_USER_3,
		POS_USER_4,
		POS_VV_I,
		POS_VA_I,
		POS_VX_I,
		POS_XSA_I,
		POS_VV_R,
		POS_VA_R,
		POS_VX_R,
		POS_XSA_R,
		POS_V:
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
