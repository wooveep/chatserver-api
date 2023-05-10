/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 10:33:02
 * @LastEditTime: 2023-05-10 12:17:14
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/utils/security/encrypt.go
 */
package security

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Encrypt use golang.org/x/crypto/bcrypt generate password
func Encrypt(src string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(src), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ValidatePassword use golang.org/x/crypto/bcrypt compare passwords for equality
func ValidatePassword(plaintext, ciphertext string) bool {
	if len(ciphertext) <= 0 {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(ciphertext), []byte(plaintext))
	return err == nil
}

func PasswordDecryption(Cipher, Key string) string {
	cipherByte, _ := base64.StdEncoding.DecodeString(Cipher)
	// if err != nil {
	// 	return nil, err
	// }
	plainText, _ := AesDecrypt([]byte(cipherByte), []byte(Key))
	return string(plainText)
}

func PasswordEncrypt(Plain, Key string) string {
	cipherText, _ := AesEncrypt([]byte(Plain), []byte(Key))
	return base64.StdEncoding.EncodeToString(cipherText)
}

func AesEncrypt(data []byte, key []byte) ([]byte, error) {
	//创建加密实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//判断加密快的大小
	blockSize := block.BlockSize()
	//填充
	encryptBytes := pkcs7Padding(data, blockSize)
	//初始化加密数据接收切片
	crypted := make([]byte, len(encryptBytes))
	//使用cbc加密模式
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	//执行加密
	blockMode.CryptBlocks(crypted, encryptBytes)
	return crypted, nil
}

// AesDecrypt 解密
func AesDecrypt(data []byte, key []byte) ([]byte, error) {
	//创建实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//获取块的大小
	blockSize := block.BlockSize()
	//使用cbc
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	//初始化解密数据接收切片
	crypted := make([]byte, len(data))
	//执行解密
	blockMode.CryptBlocks(crypted, data)
	//去除填充
	crypted, err = pkcs7UnPadding(crypted)
	if err != nil {
		return nil, err
	}
	return crypted, nil
}

// pkcs7Padding 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	//判断缺少几位长度。最少1，最多 blockSize
	padding := blockSize - len(data)%blockSize
	//补足位数。把切片[]byte{byte(padding)}复制padding个
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// pkcs7UnPadding 填充的反向操作
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	}
	//获取填充的个数
	unPadding := int(data[length-1])
	return data[:(length - unPadding)], nil
}
