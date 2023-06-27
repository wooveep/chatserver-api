/*
 * @Author: cloudyi.li
 * @Date: 2023-06-07 09:50:35
 * @LastEditTime: 2023-06-27 15:07:58
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/search/crawl.go
 */
package search

import (
	"bufio"
	"chatserver-api/pkg/logger"
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	readability "github.com/go-shiori/go-readability"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/htmlindex"
)

// var (
//
//	urls = []string{
//		"https://zh.wikipedia.org/wiki/%E5%8D%97%E4%BA%AC%E5%B8%82",
//	}
//
// )
func detectContentCharset(body io.Reader) string {
	r := bufio.NewReader(body)
	if data, err := r.Peek(1024); err == nil {
		if _, name, _ := charset.DetermineEncoding(data, ""); len(name) != 0 {
			return name
		}
	}

	return "utf-8"
}

func DecodeHTMLBody(body io.Reader, charset string) (io.Reader, error) {
	if charset == "" {
		charset = detectContentCharset(body)
	}

	e, err := htmlindex.Get(charset)
	if err != nil {
		return nil, err
	}

	if name, _ := htmlindex.Name(e); name != "utf-8" {
		body = e.NewDecoder().Reader(body)
	}

	return body, nil
}

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
	var article readability.Article
	// dialer, err := proxy.SOCKS5("tcp", "192.168.10.253:1080", nil, proxy.Direct)
	// if err != nil {
	// 	logger.Errorf("failed to download %s: %v\n", u, err)
	// }
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		// Dial:                dialer.Dial,
		MaxIdleConns:        10,
		IdleConnTimeout:     time.Second * 3,
		DisableCompression:  true,
		DisableKeepAlives:   true,
		TLSHandshakeTimeout: time.Second * 10,
	}
	c := &http.Client{
		Transport: tr,
		Timeout:   15 * time.Second,
	}
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Mobile/15E148 Safari/604.1")
	retryCount := 3
	for i := 0; i < retryCount; i++ {
		resp, err := c.Do(req)
		if err != nil {
			logger.Errorf("failed to download %s: %v\n", u, err)
			// return ""
			continue
		}
		defer resp.Body.Close()

		ur, _ := url.Parse(u)
		body, _ := DecodeHTMLBody(resp.Body, "")
		article, err = readability.FromReader(body, ur)
		if err != nil {
			logger.Errorf("failed to parse %s: %v\n", u, err)
			// return ""
			continue
		}
		if resp.StatusCode == http.StatusOK {
			// 处理响应数据
			break
		}
		// fmt.Printf("Request failed with status code %d\n", resp.StatusCode)
		// return ""
		if i < retryCount-1 {
			// fmt.Printf("Retrying request in 5 seconds...\n")
			time.Sleep(time.Second * 2)
		}
		if i == 2 {
			logger.Errorf("重试3次失败")
			return ""
		}
	}

	textcontent := []rune(delete_extra_space(article.TextContent))
	textlen := len(textcontent)
	if textlen == 0 {
		return ""
	}
	if textlen > 2200 {
		textlen = 2200
	}
	return article.Title + "\n" + string(textcontent[:textlen-1])
}
