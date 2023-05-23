/*
 * @Author: cloudyi.li
 * @Date: 2023-05-22 13:41:03
 * @LastEditTime: 2023-05-22 16:49:03
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/tika/markdown.go
 */
package tika

import (
	"os"
	"regexp"
)

func ProcessMarkDown(filename string) ([]string, error) {
	// title := regexp.MustCompile(`(?m)^#\s\S+.*`)
	content := regexp.MustCompile(`(?m)(^第\S+条[\s\S]*?)[第|#]`)
	text, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var textlist []string
	// matches := title.FindAllString(string(text), -1)
	// textlist = matches
	matches := content.FindAllStringSubmatch(string(text), -1)

	for _, v := range matches {
		textlist = append(textlist, v[1])
	}
	return textlist, nil
}
