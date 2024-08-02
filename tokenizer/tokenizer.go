package tokenizer

import (
	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
	"strings"

	. "github.com/knaka/go-utils"
)

func wordsJapanese(s string) []string {
	// BOS: Begin Of Sentence
	// EOS: End Of Sentence
	t := V(tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos()))
	var a []string
	// 以下だと、 t.Analyze(s, tokenizer.Normal) と同等のよう？
	// for _, token := range t.Tokenize(s) {
	// 検索用に tokenize。何が違うんだろう
	for _, token := range t.Analyze(s, tokenizer.Search) {
		// 何も削らずに、ただ ZWSP を足すだけにする。後で元の文書の復元が可能になる
		a = append(a, token.Surface)
	}
	return a
}

const zeroWidthSpace = "\u200B"

func SeparateJapaneseWithZWSP(text string) string {
	return strings.Join(wordsJapanese(text), zeroWidthSpace)
}

func SeparateJapanese(text string) string {
	return strings.Join(wordsJapanese(text), " ")
}

func RemoveZWSP(s string) string {
	return strings.ReplaceAll(s, zeroWidthSpace, "")
}

func init() {
	SeparateJapanese("")
}
