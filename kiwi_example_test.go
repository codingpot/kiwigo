package kiwi_test

import (
	"fmt"

	kiwi "github.com/codingpot/kiwigo"
)

func Example() {
	kb := kiwi.NewBuilder("./base", 1 /*=numThread*/, kiwi.KIWI_BUILD_INTEGRATE_ALLOMORPH /*=options*/)
	kb.AddWord("코딩냄비", "NNP", 0)

	k := kb.Build()
	defer k.Close() // don't forget to Close()!

	results, _ := k.Analyze("안녕하세요 코딩냄비입니다. 부글부글.", 1 /*=topN*/, kiwi.KIWI_MATCH_ALL)
	fmt.Println(results)
	// Output:
	// [{[{0 NNG 안녕} {2 XSA 하} {3 EF 세요} {6 NNP 코딩냄비} {10 VCP 이} {10 EF ᆸ니다} {13 SF .} {15 MAG 부글부글} {19 SF .}] -55.869953}]
}
