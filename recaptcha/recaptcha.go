package recaptcha

// reCaptcha Documentation:
// https://developers.google.com/recaptcha/docs/verify

var secret string

type recaptchaResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error-codes"`
}

func Init(newSecret string) {
	secret = newSecret
}

func Check(response, remoteip string) error {
	return nil
}
