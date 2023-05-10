/*
 * @Author: cloudyi.li
 * @Date: 2023-05-09 18:46:24
 * @LastEditTime: 2023-05-10 18:37:20
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/mail/activemail.go
 */
package mail

import (
	"bytes"
	"chatserver-api/pkg/config"
	"text/template"
)

// 发送用户激活链接邮件
var activeTemplate = `
<html lang="en" xmlns:th="http://www.thymeleaf.org">
    <head>
        <meta charset="UTF-8">
        <title>激活邮件</title>
        <style type="text/css">
            * {
                margin: 0;
                padding: 0;
                box-sizing: border-box;
                font-family: Arial, Helvetica, sans-serif;
            }

            body {
                background-color: #ECECEC;
            }

            .container {
                width: 800px;
                margin: 50px auto;
            }

            .header {
                height: 80px;
                background-color: #49bcff;
                border-top-left-radius: 5px;
                border-top-right-radius: 5px;
                padding-left: 30px;
            }

            .header h2 {
                padding-top: 25px;
                color: white;
            }

            .content {
                background-color: #fff;
                padding-left: 30px;
                padding-bottom: 30px;
                border-bottom: 1px solid #ccc;
            }

            .content h2 {
                padding-top: 20px;
                padding-bottom: 20px;
            }

            .content p {
                padding-top: 10px;
            }

            .footer {
                background-color: #fff;
                border-bottom-left-radius: 5px;
                border-bottom-right-radius: 5px;
                padding: 35px;
            }

            .footer p {
                color: #747474;
                padding-top: 10px;
            }
        </style>
    </head>

    <body>
        <div class="container">
            <div class="header">
                <h2>欢迎使用ChatServer AI助手</h2>
            </div>
            <div class="content">
                <h2>亲爱的{{ .NickName }}: 您好!</h2>
                <p>请点击链接激活<b><a href="{{ .CodeLinK }}">{{ .CodeLinK }}<a></b></p>
                <p><strong>如果链接无法点击，请复制链接到浏览器打开</strong></p>
                <p>在使用前请查看使用说明：</p>
                <p>如果后续使用有任何问题可以联系管理员，Email: <b>cloudyi@wooveep.net</b></p>
            </div>
            <div class="footer">
                <p>此为系统邮件，请勿回复</p>
                <p>请保管好您的信息，避免被他人盗用</p>
                <p>©wooveep.net</p>
            </div>
        </div>
    </body>

</html> 
`

func SendActiceCode(email string, nickname string, activeCode string) error {
	var bodyBytes bytes.Buffer
	CodeLink := config.AppConfig.ExternalURL + "#/active/" + activeCode
	tpl := template.Must(template.New("").Parse(activeTemplate))
	err := tpl.Execute(&bodyBytes, map[string]interface{}{"CodeLinK": CodeLink, "NickName": nickname})
	if err != nil {
		return err
	}
	body := bodyBytes.String()
	err = send([]string{email}, "Chatserver用户激活邮件", body)
	if err != nil {
		return err
	}
	return nil
}
