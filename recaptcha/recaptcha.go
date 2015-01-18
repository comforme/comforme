package recaptcha

// reCaptcha Documentation:
// https://developers.google.com/recaptcha/docs/verify

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var secret string

var recaptchaError = errors.New("Invalid ReCaptcha")

type recaptchaResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error-codes"`
}

func Init(newSecret string) {
	secret = newSecret
}

func Check(response, remoteip string) error {
	apiEndpoint := "https://www.google.com/recaptcha/api/siteverify"
	params := fmt.Sprintf("?secret=%s&response=%s&remoteip=%s",
		secret,
		response,
		remoteip)
	resp, err := http.Get(apiEndpoint + params)
	defer resp.Body.Close()
	if err == nil {
		body, err := ioutil.ReadAll(resp.Body)
		log.Println("reCaptcha result:", body)
		if err == nil {
			var data recaptchaResult
			json.Unmarshal(body, &data)
			if data.Success {
				return nil
			}
			if data.Error != "" {
				err = errors.New(data.Error)
			} else {
				err = recaptchaError
			}
			return err
		}
		return err
	}
	return err
}
