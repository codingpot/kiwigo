package kiwi_test

import (
	"fmt"

	kiwi "github.com/codingpot/kiwigo"
)

func ExampleAnalyze() {
	k := kiwi.New("./ModelGenerator", 1 /*=numThread*/, kiwi.KIWI_BUILD_DEFAULT /*=options*/)
	defer k.Close()

	results, err := k.Analyze("안녕하세요 코딩냄비입니다. 부글부글.", 1 /*=topN*/, kiwi.KIWI_MATCH_ALL)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(results)

	// Output:
	// [{[{0 NNG 안녕} {2 XSA 하} {4 EP 시} {3 EC 어요} {6 NNG 코딩} {8 NNG 냄비} {10 VCP 이} {11 EF ᆸ니다} {13 SF .} {15 NNP 부글부} {18 NNG 글} {19 SF .}] -92.20007}]
}

func ExempleAddWord() {
	kb := kiwi.NewBuilder("./ModelGenerator", 1 /*=numThread*/, kiwi.KIWI_BUILD_INTEGRATE_ALLOMORPH /*=options*/)
	defer kb.Close()

	kb.AddWord("안녕하세요", "NNG", 0)

	k := kb.Build()
	defer k.Close()
	results, err := k.Analyze("안녕하세요 코딩냄비입니다. 부글부글", 1, kiwi.KIWI_MATCH_ALL)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(results)

	// Output:
	// [{[{0 NNG 안녕하세요} {6 NNP 코딩냄비} {10 VCP 이} {11 EF ᆸ니다} {13 SF .} {15 NNP 부글부} {18 NNG 글}] -83.74562}]
}
