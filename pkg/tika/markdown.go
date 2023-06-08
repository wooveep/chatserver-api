/*
 * @Author: cloudyi.li
 * @Date: 2023-05-22 13:41:03
 * @LastEditTime: 2023-05-29 10:23:03
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/tika/markdown.go
 */
package tika

import (
	"os"
	"regexp"
)

type Document struct {
	Title    string
	Document string
	Children []Document
}

func ProcessMarkDown(filename string) ([]string, error) {
	title1 := regexp.MustCompile(`(?m)^#\s([\s\S]*?)\n`)
	// title2 := regexp.MustCompile(`(?m)^##\s([\s\S]*?)\n`)
	// title3 := regexp.MustCompile(`(?m)^###\s([\s\S]*?)\n`)
	// title4 := regexp.MustCompile(`(?m)^####\s([\s\S]*?)\n`)
	// title5 := regexp.MustCompile(`(?m)^#####\s([\s\S]*?)\n`)

	// info := regexp.MustCompile(`(?m)^#\s(\S*)`)
	// content := regexp.MustCompile(`(?m)(^第\S+条[\s\S]*?)[第|#]`)
	text, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var textlist []string
	// matches := title.FindAllString(string(text), -1)
	// textlist = matches
	// matches := content.FindAllStringSubmatch(string(text), -1)
	matche_title1 := title1.FindAllStringSubmatch(string(text), 1)
	// matche_title2 := title2.FindAllStringSubmatch(string(text), -1)

	// matche_info := info.FindAllStringSubmatch(string(text), 1)

	textlist = append(textlist, matche_title1[0][1])

	// for _, v := range matche_title2 {
	// 	matche_title3 :=
	// 	textlist = append(textlist, v[1])
	// }
	return textlist, nil
}
