package recaptcha

// reCaptcha Documentation:
// https://developers.google.com/recaptcha/docs/verify

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var secret string

var recaptchaError = errors.New("Invalid ReCaptcha")

type recaptchaResult struct {
	Success bool     `json:"success"`
	Errors  []string `json:"error-codes"`
}

func Init(newSecret string) {
	secret = newSecret
}

func Check(response, remoteip string) error {
	var Url *url.URL
	Url, err := url.Parse("https://www.google.com/recaptcha/api/siteverify")
	if err != nil {
		return err
	}

	parameters := url.Values{}
	parameters.Add("secret", secret)
	parameters.Add("response", response)
	parameters.Add("remoteip", remoteip)
	Url.RawQuery = parameters.Encode()

	log.Println("Making reCaptcha request:", Url.String())
	resp, err := http.Get(Url.String())
	defer resp.Body.Close()
	if err == nil {
		body, err := ioutil.ReadAll(resp.Body)
		log.Println("reCaptcha result:", string(body))
		if err == nil {
			var data recaptchaResult
			json.Unmarshal(body, &data)
			if data.Success {
				return nil
			}
			if len(data.Errors) >= 1 {
				err = errors.New("reCaptcha error(s): " + fmt.Sprintf("%v", data.Errors))
			} else {
				err = recaptchaError
			}
			return err
		}
		return err
	}
	return err
}
