/*
 * @Author: cloudyi.li
 * @Date: 2023-02-15 11:49:33
 * @LastEditTime: 2023-03-29 09:33:48
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/cmd/chatserver-api/args.go
 */
package chatserverapi

import "github.com/alexflint/go-arg"

type Args struct {
	Names  string `arg:"-n,--names,required" help:"wechat"`
	Config string `arg:"-c,--config" help:"config file" default:"configs/config.toml"`
}

func (Args) Version() string {
	return ""
}

func LoadArgsValid() Args {
	var args Args
	arg.MustParse(&args)
	return args
}
