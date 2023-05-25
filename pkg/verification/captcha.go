/*
 * @Author: cloudyi.li
 * @Date: 2023-05-24 09:52:31
 * @LastEditTime: 2023-05-25 21:43:19
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/verification/captcha.go
 */
package verification

import (
	"bytes"
	"chatserver-api/internal/consts"
	"chatserver-api/pkg/cache"
	"chatserver-api/pkg/logger"
	"chatserver-api/utils/security"
	"context"
	"encoding/base64"
	"image/color"
	"image/png"
	"time"

	afcap "github.com/afocus/captcha"
	"github.com/go-redis/redis/v8"
)

func getCaptchaCodeKey(code string) string {
	return consts.CaptchaPrefix + security.Md5(code)
}

func GenerateCaptcha(ctx context.Context) (imgbase string, err error) {
	cap := afcap.New()
	rc := cache.GetRedisClient()
	timer := 600 * time.Second
	cap.SetFont("./fonts/comic.ttf")
	// 设置验证码大小
	cap.SetSize(128, 64)
	// 设置干扰强度
	cap.SetDisturbance(afcap.MEDIUM)
	// 设置前景色 可以多个 随机替换文字颜色 默认黑色
	cap.SetFrontColor(color.RGBA{255, 255, 255, 255})
	// 设置背景色 可以多个 随机替换背景色 默认白色
	cap.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})
	img, code := cap.Create(4, afcap.NUM)
	buffer := new(bytes.Buffer)

	err = png.Encode(buffer, img)
	if err != nil {
		return
	}
	err = rc.SetNX(ctx, getCaptchaCodeKey(code), code, timer).Err()
	if err != nil {
		return
	}
	imgbase = base64.StdEncoding.EncodeToString(buffer.Bytes())
	return
}

func VerifyCaptcha(ctx context.Context, code string) bool {
	rc := cache.GetRedisClient()
	code_re, err := rc.Get(ctx, getCaptchaCodeKey(code)).Result()
	if err != nil {
		if err != redis.Nil {
			logger.Errorf("Redis连接异常:%v", err.Error())
		}
		return false
	}
	rc.Del(ctx, getCaptchaCodeKey(code))
	if code == code_re {
		return true
	} else {
		return false
	}
}
