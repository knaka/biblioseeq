package fts

import (
	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
	"log"
	"strings"
)

func wordsJapanese(s string) []string {
	// BOS: Begin Of Sentence
	// EOS: End Of Sentence
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		log.Panicf("panic 42afe51 (%v)", err)
	}
	// これは t.Analyze(s, tokenizer.Normal)
	// tokens := t.Tokenize(s)
	// 検索用の tokenizer ？ 何が違うんだろう
	tokens := t.Analyze(s, tokenizer.Search)
	var a []string
	for _, token := range tokens {
		// 何も削らずに、ただ ZWSP を足すだけにすれば、後で元の文書の復元が可能
		a = append(a, token.Surface)
	}
	return a
}

// Zero-Width Separator
const zwsp = "\u200B"

// 日本語とは、ZWSP で分かち書きされた言語であると
func divideJapaneseToWordsWithZwsp(text string) string {
	return strings.Join(wordsJapanese(text), zwsp)
}

// todo: 初期呼び出しの、良い方法は？
func init() {
	divideJapaneseToWords("")
}

// クエリとして渡すにはスペースでも良いか
func divideJapaneseToWords(text string) string {
	return strings.Join(wordsJapanese(text), " ")
}

func removeZwsp(s string) string {
	return strings.ReplaceAll(s, zwsp, "")
}
