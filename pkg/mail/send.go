/*
 * @Author: cloudyi.li
 * @Date: 2023-05-09 14:27:26
 * @LastEditTime: 2023-05-16 21:12:47
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/mail/send.go
 */
package mail

import (
	"chatserver-api/pkg/config"
	"fmt"
	"net"
	"strconv"
	"time"

	gomail "github.com/go-mail/mail"
	"golang.org/x/net/proxy"
)

func send(mailTo []string, subject, body string) error {
	mailcfg := config.AppConfig.EmailCofig
	mailConn := map[string]string{
		"user": mailcfg.Username,
		"pass": mailcfg.Password,
		"host": mailcfg.Host,
		"port": mailcfg.Port,
	}

	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int

	m := gomail.NewMessage()

	m.SetHeader("From", m.FormatAddress(mailcfg.Sender, "Chatserver"))
	m.SetHeader("To", mailTo...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])
	if mailcfg.ProxyMode == "socks5" {
		gomail.NetDialTimeout = func(network string, address string, timeout time.Duration) (net.Conn, error) {
			proxyAdd := fmt.Sprintf("%s:%s", mailcfg.ProxyIP, mailcfg.ProxyPort)
			dialer, err := proxy.SOCKS5("tcp", proxyAdd, nil, proxy.Direct)
			if err != nil {
				return nil, err
			}
			return dialer.Dial("tcp", d.Host+":"+strconv.Itoa(d.Port))
		}
	}
	d.Timeout = 20 * time.Second

	for i := 1; i <= 3; i++ {
		err := d.DialAndSend(m)
		if err == nil {
			return nil
		}
		if i == 3 {
			return err
		}
		time.Sleep(5)
	}
	return nil
}
