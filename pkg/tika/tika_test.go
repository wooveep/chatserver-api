/*
 * @Author: cloudyi.li
 * @Date: 2023-04-28 11:37:30
 * @LastEditTime: 2023-04-29 22:34:32
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/tika/tika_test.go
 */
package tika

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestParseDocFile(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// name := ParseDocFile()
			// fmt.Print(name)
			name := "tmpfile.txt"
			list := ParseHtml(name)
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
			var list2 []string
			var text string
			// var len int
			for _, v := range list {

				if len(text) < 400 {
					text += v
				} else {
					list2 = append(list2, text)
					text = ""
				}

			}
			file, err := os.Create("output.txt")
			if err != nil {
				fmt.Println("Error creating file:", err)
				return
			}
			defer file.Close()

			writer := bufio.NewWriter(file)

			// for _, v := range textbodylist {
			// 	writer.WriteString(v.textbody)
			// 	writer.WriteString("\n-------" + strconv.Itoa(v.tokens) + "-----------\n")
			// }
			for _, v := range list2 {
				writer.WriteString(v)
				writer.WriteString("\n-------------\n")
			}
			writer.Flush()

			fmt.Println("Output saved to file successfully!")
		})
	}
}
