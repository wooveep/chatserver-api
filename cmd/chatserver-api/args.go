/*
 * @Author: cloudyi.li
 * @Date: 2023-02-15 11:49:33
 * @LastEditTime: 2023-03-29 11:28:48
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/cmd/chatserver-api/args.go
 */
package chatserverapi

import (
	"chatserver-api/utils/version"

	"github.com/alexflint/go-arg"
)

type Args struct {
	Config string `arg:"-c,--config" help:"config file" default:"configs/config.yml"`
}

func (Args) Version() string {
	return version.PrintVersion()
}

func LoadArgsValid() Args {
	var args Args
	arg.MustParse(&args)
	return args
}
