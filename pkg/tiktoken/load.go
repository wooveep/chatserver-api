/*
 * @Author: cloudyi.li
 * @Date: 2023-04-23 11:35:31
 * @LastEditTime: 2023-04-23 11:39:43
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/tiktoken/load.go
 */
package tiktoken

import (
	"chatserver-api/pkg/tiktoken/assets"
	"encoding/base64"
	"strconv"
	"strings"
)

func loadTiktokenBpe(filename string) (map[string]int, error) {
	contents, err := assets.Assets.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	bpeRanks := make(map[string]int)
	for _, line := range strings.Split(string(contents), "\n") {
		if line == "" {
			continue
		}
		parts := strings.Split(line, " ")
		token, err := base64.StdEncoding.DecodeString(parts[0])
		if err != nil {
			return nil, err
		}
		rank, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
		bpeRanks[string(token)] = rank
	}
	return bpeRanks, nil
}
