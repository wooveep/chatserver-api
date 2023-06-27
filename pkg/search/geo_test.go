/*
 * @Author: cloudyi.li
 * @Date: 2023-06-15 10:49:51
 * @LastEditTime: 2023-06-15 13:36:12
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/search/geo_test.go
 */
package search

import (
	"chatserver-api/pkg/config"
	"chatserver-api/pkg/logger"
	"fmt"
	"testing"
)

func Test_geoSearch(t *testing.T) {
	c := config.Load("../../configs/config.yml")
	logger.InitLogger(&c.LogConfig, c.AppName)
	type args struct {
		address string
	}
	tests := []struct {
		name string
		args args
		// want string
	}{
		{
			name: "test3",
			args: args{
				address: "南京",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(GeoSearch(tt.args.address))
		})
	}
}
