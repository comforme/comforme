package recaptcha

import(

)

var secret string

func Init(newSecret string) {
	secret = newSecret
}

func Check(response, remoteip string) error {
	return nil
}
