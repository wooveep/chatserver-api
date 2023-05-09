/*
 * @Author: cloudyi.li
 * @Date: 2023-04-28 11:37:30
 * @LastEditTime: 2023-04-29 22:27:56
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/tika/tika.go
 */
package tika

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	gotika "github.com/google/go-tika/tika"
	"golang.org/x/net/html"
)

// import "os"
func ParseDocFile() (tmpfilename string) {
	f, err := os.Open("sinovatio2021.PDF")
	if err != nil {
		// logger.Fatal(err)
		fmt.Print(err)
	}
	defer f.Close()

	client := gotika.NewClient(nil, "http://192.168.10.251:18081/")
	// body, err := client.MetaRecursive(context.Background(), f)
	body, err := client.MetaRecursiveType(context.Background(), f, "html")
	text := body[0]["X-TIKA:content"][0]
	tmpfile, err := os.Create("tmpfile.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer tmpfile.Close()
	_, err = tmpfile.WriteString(text)
	if err != nil {
		return
	}
	// list := strings.Split(text, "\n")

	// type textbodys struct {
	// 	textbody string
	// 	tokens   int
	// }
	// var textbodylist []textbodys
	// var textbody string
	// var tokens int
	// for _, v := range list {
	// 	if tokens < 300 {
	// 		textbody += v
	// 		tokens += tiktoken.NumTokensSingleString(v)
	// 	} else {
	// 		textbodylist = append(textbodylist, textbodys{textbody, tokens})
	// 		textbody = ""
	// 		tokens = 0
	// 	}
	// }
	// file, err := os.Create("output.txt")
	// if err != nil {
	// 	fmt.Println("Error creating file:", err)
	// 	return
	// }
	// defer file.Close()

	// writer := bufio.NewWriter(file)

	// for _, v := range textbodylist {
	// 	writer.WriteString(v.textbody)
	// 	writer.WriteString("\n-------" + strconv.Itoa(v.tokens) + "-----------\n")
	// }

	// writer.Flush()
	return "tmpfile.txt"
}

func ParseHtml(tmpfilename string) (textList []string) {
	data, err := ioutil.ReadFile(tmpfilename)
	if err != nil {
		// logger.Fatal(err)
		fmt.Print(err)
	}
	// re := regexp.MustCompile(`\s+`)
	pContentre1 := regexp.MustCompile("\\n")
	text := string(data)
	// // re := regexp.MustCompile(`(?<![0-9])\n[。！？.!?](?!\d)`)
	// re := regexp.MustCompile(`[。！？；.!?;]$`)
	// re2 := regexp.MustCompile(`\d+\.$`)
	// re3 := regexp.MustCompile(`^\s+[:：]`)
	doc, err := html.Parse(strings.NewReader(text))
	if err != nil {
		fmt.Print(err)
	}

	var divContents string

	var traverseNode func(*html.Node, []string)
	traverseNode = func(n *html.Node, path []string) {
		if n.Type == html.ElementNode && n.Data == "div" {
			divPath := append(path, "div")
			var divContent string
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "div" {
					traverseNode(c, divPath)
				} else if c.Type == html.ElementNode && c.Data == "p" {
					var pContent []string
					for cc := c.FirstChild; cc != nil; cc = cc.NextSibling {
						if cc.Type == html.TextNode && cc.Data != "" {
							a := pContentre1.ReplaceAllString(cc.Data, "")
							// if re.MatchString(a) && !re2.MatchString(a) {
							// 	pContent = append(pContent, v, "\n")
							// } else if re3.MatchString(v) {
							// 	pContent = append(pContent, "\n", v)
							// } else {
							pContent = append(pContent, a)
							// }

						}
					}
					if len(pContent) > 0 {
						divContent += strings.Join(pContent, " ")
					}
				}
			}
			divContents += divContent
		} else {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				traverseNode(c, path)
			}
		}
	}

	traverseNode(doc, []string{})

	textList = strings.Split(divContents, " ")
	fmt.Printf("ttttintitnt")
	// textList = divContents
	return
	// var textbodylist []string
	// var textbody string

	// var tokens int
	// for _, v := range textList {

	// 	tokens += tiktoken.NumTokensSingleString(v)
	// 	if tokens > 300 {
	// 		textbodylist = append(textbodylist, textbody)
	// 		textbody = ""
	// 	} else {
	// 		textbody += v
	// 	}
	// }
	// fmt.Println(textList)

	// file, err := os.Create("output.txt")
	// if err != nil {
	// 	fmt.Println("Error creating file:", err)
	// 	return
	// }
	// defer file.Close()

	// writer := bufio.NewWriter(file)

	// for _, v := range textList {
	// 	writer.WriteString(v)
	// 	writer.WriteString("\n-------------------\n")
	// }

	// writer.Flush()

	// fmt.Println("Output saved to file successfully!")
}
