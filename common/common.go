package common

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/keighl/mandrill"
)

const (
	alphaNumeric            = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	sessionIdLength         = 25
	generatedPasswordLength = 15
	fromEmail               = "donotreply@comfor.me"
	fromName                = "ComFor.Me"
)

// Errors
var EmailFailed = errors.New("Sending email failed.")

var mandrillKey = os.Getenv("MANDRILL_APIKEY")
var emailRegex *regexp.Regexp

func init() {
	rand.Seed(time.Now().Unix() ^ int64(time.Now().Nanosecond()))
	emailRegex = regexp.MustCompile("^.+@.+\\..+$")
}

func RandSeq(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphaNumeric[rand.Intn(len(alphaNumeric))]
	}
	return string(b)
}

func NewSessionID() string {
	return RandSeq(sessionIdLength)
}

func GenPassword() string {
	return RandSeq(generatedPasswordLength)
}

func ExecTemplate(tmpl *template.Template, w http.ResponseWriter, pc map[string]interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := tmpl.Execute(w, pc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func SendRegEmail(email, password string) error {
	emailText := fmt.Sprintf(`Thank you for registerering with ComFor.Me.

Your temporary password is: %s

Please change your password after logging in.
`, password)
	return sendEmail(email, "Welcome to ComFor.Me!", "", emailText)
}

func SendResetEmail(email, password string) error {
	emailText := fmt.Sprintf(`We received a password reset request for your account on ComFor.Me.

Your new temporary password is: %s

Please change your password after logging in.

If you did not request this password reset please contact support.
`, password)
	return sendEmail(email, "ComFor.Me Password Reset", "", emailText)
}

func sendEmail(recipient, subject, html, text string) error {
	log.Printf("Sending email to: %s\n", recipient)
	log.Printf("Subject: %s\nText:\n%s\n", subject, text)

	client := mandrill.ClientWithKey(mandrillKey)

	message := &mandrill.Message{}
	message.AddRecipient(recipient, recipient, "to")
	message.FromEmail = fromEmail
	message.FromName = fromName
	message.Subject = subject
	message.HTML = html
	message.Text = text

	responses, apiError, err := client.MessagesSend(message)
	if err != nil || apiError != nil {
		if err != nil {
			log.Printf("Error: %s\n", err.Error())
		}
		if apiError != nil {
			log.Printf("Mandrill API Error: %+v\n", apiError)
		}
		return EmailFailed
	}
	log.Printf("Mandrill responses: %+v\n", responses)
	return nil
}

func ValidEmail(email string) bool {
	return emailRegex.Match([]byte(email))
}

func SetSessionCookie(res http.ResponseWriter, sessionid string) {
	http.SetCookie(res, &http.Cookie{Name: "sessionid", Value: sessionid, Expires: time.Now().AddDate(10, 0, 0)})
}
