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
	"runtime"
	"strings"
	"time"

	"github.com/keighl/mandrill"
)

var DebugMode = os.Getenv("DEBUG_MODE") == "true"

const (
	alphaNumeric            = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	sessionIdLength         = 25
	generatedPasswordLength = 15
	fromEmail               = "donotreply@comfor.me"
	fromName                = "ComFor.Me"
)

// Database row types
type Community struct {
	Id       int
	Name     string
	IsMember bool
}

type Page struct {
	Id          int
	Title       string
	Slug        string
	Category    string
	Description string
	DateCreated time.Time
}

// Errors
var EmailFailed = errors.New("Sending email failed.")
var EmailInUse = errors.New("You have already registered with this email address.")
var UsernameInUse = errors.New("This username is in use. Please select a different one.")
var InvalidUsernameOrPassword = errors.New("Invalid username or password.")
var DatabaseError = errors.New("Unknown database error.")
var InvalidSessionID = errors.New("Invalid sessionid.")
var InvalidEmail = errors.New("The provided email address is not valid.")
var InvalidIpAddress = errors.New("There is something wrong with your IP address.")

var mandrillKey = os.Getenv("MANDRILL_APIKEY")
var emailRegex *regexp.Regexp
var ipAddressRegex *regexp.Regexp

func init() {
	rand.Seed(time.Now().Unix() ^ int64(time.Now().Nanosecond()))
	emailRegex = regexp.MustCompile("^.+@.+\\..+$")
	ipAddressRegex = regexp.MustCompile("(.+):\\d+$")
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

func GenSlug(seed string) string {
	slug := strings.ToLower(seed)
	return slug
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

func Logout(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "/logout", http.StatusFound)
}

func LogError(err error) {
	pc, file, line, ok := runtime.Caller(1)
	if ok {
		function := runtime.FuncForPC(pc)
		log.Printf("Error occurred in function (%s) at (%s:%d): %s\n", function.Name(), file, line, err.Error())
	} else {
		log.Println("Error occurred when trying to lookup caller info for the function that generated the error:", err)
	}
}

func LogErrorSkipLevels(err error, levels int) {
	pc, file, line, ok := runtime.Caller(levels + 1)
	if ok {
		function := runtime.FuncForPC(pc)
		log.Printf("Error occurred in function (%s) at (%s:%d): %s\n", function.Name(), file, line, err.Error())
	} else {
		log.Println("Error occurred when trying to lookup caller info for the function that generated the error:", err)
	}
}

func GetIpAddress(req *http.Request) (string, error) {
	addrParts := ipAddressRegex.FindSubmatch([]byte(req.RemoteAddr))
	if len(addrParts) != 2 {
		log.Println("The following remote address did not allow IP address extraction:", req.RemoteAddr)
		return "", InvalidIpAddress
	}
	return string(addrParts[1]), nil
}
