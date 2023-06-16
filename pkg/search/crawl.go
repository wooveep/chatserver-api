/*
 * @Author: cloudyi.li
 * @Date: 2023-06-07 09:50:35
 * @LastEditTime: 2023-06-16 07:37:15
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/search/crawl.go
 */
package search

import (
	"chatserver-api/pkg/logger"
	"crypto/tls"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	readability "github.com/go-shiori/go-readability"
)

// var (
// 	urls = []string{
// 		"https://zh.wikipedia.org/wiki/%E5%8D%97%E4%BA%AC%E5%B8%82",
// 	}
// )

func delete_extra_space(s string) string {
	//删除字符串中的多余空格，有多个空格时，仅保留一个空格
	s0 := strings.Replace(s, "	", " ", -1) //替换tab为空格
	s1 := strings.Replace(s0, "\n", " ", -1)
	regstr := "\\s{2,}"                          //两个及两个以上空格的正则表达式
	reg, _ := regexp.Compile(regstr)             //编译正则表达式
	s2 := make([]byte, len(s1))                  //定义字符数组切片
	copy(s2, s1)                                 //将字符串复制到切片
	spc_index := reg.FindStringIndex(string(s2)) //在字符串中搜索
	for len(spc_index) > 0 {                     //找到适配项
		s2 = append(s2[:spc_index[0]+1], s2[spc_index[1]:]...) //删除多余空格
		spc_index = reg.FindStringIndex(string(s2))            //继续在字符串中搜索
	}
	return string(s2)
}

func crawlPage(u string) string {
	// dialer, err := proxy.SOCKS5("tcp", "192.168.10.253:1080", nil, proxy.Direct)
	// if err != nil {
	// 	logger.Errorf("failed to download %s: %v\n", u, err)
	// }
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		// Dial:            dialer.Dial,
	}
	c := &http.Client{
		Transport: tr,
		Timeout:   9 * time.Second,
	}
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Mobile/15E148 Safari/604.1")

	resp, err := c.Do(req)
	if err != nil {
		logger.Errorf("failed to download %s: %v\n", u, err)
		return ""
	}
	defer resp.Body.Close()

	ur, _ := url.Parse(u)
	article, err := readability.FromReader(resp.Body, ur)
	if err != nil {
		logger.Errorf("failed to parse %s: %v\n", u, err)
		return ""
	}
	textcontent := delete_extra_space(article.TextContent)
	textlen := len(textcontent)
	if textlen == 0 {
		return ""
	}
	if textlen > 5000 {
		textlen = 5000
	}
	return article.Title + "\n" + textcontent[:textlen-1]
}
