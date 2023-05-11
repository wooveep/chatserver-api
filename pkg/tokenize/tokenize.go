/*
 * @Author: cloudyi.li
 * @Date: 2023-05-08 14:04:14
 * @LastEditTime: 2023-05-11 23:38:58
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/tokenize/tokenize.go
 */
package tokenize

import (
	"path"

	"github.com/yanyiwu/gojieba"
)

var _ Tokenizer = (*tokenizer)(nil)

type Tokenizer interface {
	GetKeyword(s string) (keyword string)
}

type tokenizer struct {
	jieba *gojieba.Jieba
}

func NewTokenizer() *tokenizer {
	dictDir := "./dict"
	jiebaPath := path.Join(dictDir, "jieba.dict.utf8")
	hmmPath := path.Join(dictDir, "hmm_model.utf8")
	userPath := path.Join(dictDir, "user.dict.utf8")
	idfPath := path.Join(dictDir, "idf.utf8")
	stopPath := path.Join(dictDir, "stop_words.utf8")
	tokenzier := gojieba.NewJieba(jiebaPath, hmmPath, userPath, idfPath, stopPath)

	return &tokenizer{
		jieba: tokenzier,
	}
}

func (t *tokenizer) GetKeyword(s string) (keyword string) {
	// tokenzier.Free()
	words := t.jieba.Extract(s, 35)
	for _, v := range words {
		keyword += v + " "
	}
	return
}
