package recaptcha

import(

)

var secret string

type struct recaptchaResult {
	Success bool `json:"success"`
	Error string `json:"error-codes"`
}

func Init(newSecret string) {
	secret = newSecret
}

func Check(response, remoteip string) error {
	return nil
}
