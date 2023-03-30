/*
 * @Author: cloudyi.li
 * @Date: 2023-03-29 10:33:02
 * @LastEditTime: 2023-03-29 10:33:08
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/utils/security/encrypt.go
 */
package security

import "golang.org/x/crypto/bcrypt"

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
	if err != nil {
		return false
	}
	return true
}
