/*
 * @Author: cloudyi.li
 * @Date: 2023-05-09 14:27:26
 * @LastEditTime: 2023-05-09 17:17:21
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/sendmail/login.go
 */
package sendmail

import (
	"chatserver-api/pkg/config"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/go-mail/mail"
	"golang.org/x/net/proxy"
)

func LoginMail(mailTo []string, activationCode string) error {
	mailcfg := config.AppConfig.EmailCofig
	mailConn := map[string]string{
		"user": mailcfg.Username,
		"pass": mailcfg.Password,
		"host": mailcfg.Host,
		"port": mailcfg.Port,
	}

	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int

	m := mail.NewMessage()

	// 配置邮件内容
	subject := "你的验证激活码"
	body := fmt.Sprintf("Dear ,\n\nYour activation code is: %s\n\nBest regards,\nYour Name", activationCode)

	m.SetHeader("From", m.FormatAddress(mailcfg.Sender, "Chatserver")) //这种方式可以添加别名，即“XX官方”
	//说明：如果是用网易邮箱账号发送，以下方法别名可以是中文，如果是qq企业邮箱，以下方法用中文别名，会报错，需要用上面此方法转码
	//m.SetHeader("From", "FB Sample"+"<"+mailConn["user"]+">") //这种方式可以添加别名，即“FB Sample”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	//m.SetHeader("From", mailConn["user"])
	m.SetHeader("To", mailTo...)    //发送给多个用户
	m.SetHeader("Subject", subject) //设置邮件主题
	m.SetBody("text/html", body)    //设置邮件正文

	d := mail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])
	if mailcfg.ProxyMode == "socks5" {
		mail.NetDialTimeout = func(network string, address string, timeout time.Duration) (net.Conn, error) {
			proxyAdd := fmt.Sprintf("%s:%s", mailcfg.ProxyIP, mailcfg.ProxyPort)
			dialer, err := proxy.SOCKS5("tcp", proxyAdd, nil, proxy.Direct)
			if err != nil {
				return nil, err
			}
			return dialer.Dial("tcp", d.Host+":"+strconv.Itoa(d.Port))
		}
	}
	// Send
	err := d.DialAndSend(m)
	return err

}
