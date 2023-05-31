/*
 * @Author: cloudyi.li
 * @Date: 2023-05-31 17:32:37
 * @LastEditTime: 2023-05-31 22:27:05
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/mail/verifymail.go
 */
package mail

import (
	"errors"

	emailverifier "github.com/AfterShip/email-verifier"
)

type Verifier struct {
	verifier *emailverifier.Verifier
}

func NewVerifier() *Verifier {
	return &Verifier{
		verifier: emailverifier.NewVerifier().EnableSMTPCheck().DisableCatchAllCheck().FromEmail("user@wooveep.net"),
	}
}
func (v *Verifier) VerifierEmail(email string) error {

	ret, err := v.verifier.Verify(email)
	if err != nil {
		// logger.Warnf("verify email address failed, error is: ", err)
		return err
	}
	if !ret.Syntax.Valid {
		// logger.Warnf("email address syntax is invalid")
		return errors.New("email address syntax is invalid")
	}
	if !ret.SMTP.Deliverable {
		return errors.New("email address not deliverable")
	}
	return nil
}
