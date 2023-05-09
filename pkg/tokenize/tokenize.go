/*
 * @Author: cloudyi.li
 * @Date: 2023-05-08 14:04:14
 * @LastEditTime: 2023-05-09 13:45:23
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/tokenize/tokenize.go
 */
package tokenize

import "github.com/yanyiwu/gojieba"

var _ Tokenizer = (*tokenizer)(nil)

type Tokenizer interface {
	GetKeyword(s string) (keyword string)
}

type tokenizer struct {
	jieba *gojieba.Jieba
}

func NewTokenizer() *tokenizer {
	tokenzier := gojieba.NewJieba()

	return &tokenizer{
		jieba: tokenzier,
	}
}

func (t *tokenizer) GetKeyword(s string) (keyword string) {
	// tokenzier.Free()
	words := t.jieba.Extract(s, 10)
	for _, v := range words {
		keyword += v + " "
	}
	return
}
