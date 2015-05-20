package common

import (
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/keighl/mandrill"
	"golang.org/x/crypto/scrypt"
)

var DebugMode = os.Getenv("DEBUG_MODE") == "true"
var secret = []byte(os.Getenv("SECRET"))
var linkAgeLimit = time.Hour * 24 * 14

const (
	alphaNumeric            = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	sessionIdLength         = 25
	generatedPasswordLength = 15
	fromEmail               = "donotreply@comfor.me"
	Protocol                = "https"
	Domain                  = "comfor.me"
	fromName                = "ComFor.Me"
	slugChars               = "A-Za-z0-9"
	slugRemoveChars         = "'\""
	slugLength              = 25
	MinDescriptionLength    = 3
)

type UserInfo struct {
	SessionID string
	Email     string
	Username  string
	UserID    int
}

// Database row types
type Community struct {
	Id       int
	Name     string
	IsMember bool
}

type Page struct {
	Id           int
	Title        string
	PageSlug     string
	Category     string
	CategorySlug string
	Description  string
	Address      string
	Website      string
	DateCreated  time.Time
}

type Post struct {
	Author           string
	Body             string
	CommonCategories int
	Date             string
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
var InvalidTitle = errors.New("Invalid page title.")
var PageAlreadyExists = errors.New("A page with this category and title already exists.")
var PageNotFound = errors.New("Page not found.")

var mandrillKey = os.Getenv("MANDRILL_APIKEY")
var emailRegex *regexp.Regexp
var ipAddressRegex *regexp.Regexp
var slugFrontCap *regexp.Regexp
var slugEndCap *regexp.Regexp
var slugRemove *regexp.Regexp
var slugMiddle *regexp.Regexp

func init() {
	rand.Seed(time.Now().Unix() ^ int64(time.Now().Nanosecond()))
	emailRegex = regexp.MustCompile("^.+@.+\\..+$")
	ipAddressRegex = regexp.MustCompile("(.+):\\d+$")
	slugFrontCap = regexp.MustCompile("[^" + slugChars + "]+$")
	slugEndCap = regexp.MustCompile("^[^" + slugChars + "]+")
	slugMiddle = regexp.MustCompile("[^" + slugChars + "]+")
	slugRemove = regexp.MustCompile("[" + slugRemoveChars + "]")
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

func GenSlug(title string) string {
	titleBytes := []byte(title)
	slug := slugFrontCap.ReplaceAll(titleBytes, []byte(""))
	slug = slugEndCap.ReplaceAll(titleBytes, []byte(""))
	slug = slugRemove.ReplaceAll(titleBytes, []byte(""))
	slug = slugMiddle.ReplaceAll(titleBytes, []byte("-"))
	return strings.ToLower(string(slug))
}

func ExecTemplate(tmpl *template.Template, w http.ResponseWriter, pc map[string]interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := tmpl.Execute(w, pc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func SendRegEmail(email string) error {
	hash, date, err := GenerateSecret(email)
	if err != nil {
		return err
	}

	emailText := fmt.Sprintf(`Thank you for registering with ComFor.Me.

To complete your registration, please copy and paste the following link into your web browser:
https://comfor.me/wizard?action=register&email=%s&date=%s&code=%s

This link will be valid for 14 days.

Hope to see you soon,
The ComFor.Me team
`, email, date, hash)
	return sendEmail(email, "Welcome to ComFor.Me!", "", emailText)
}

func SendResetEmail(email, date, hash string) error {
	emailText := fmt.Sprintf(`We received a password reset request for your account on ComFor.Me.

To complete your password reset, please copy and paste the following link into your web browser:
https://comfor.me/wizard?action=reset&email=%s&date=%s&code=%s

If you did not request this password reset you can safely ignore this email.

This link will be valid for 14 days.

Hope to see you soon,
The ComFor.Me team
`, email, date, hash)
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

	responses, err := client.MessagesSend(message)
	if err != nil {
		log.Printf("Error: %s\n", err.Error())
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

func GetIpAddress(req *http.Request) string {
	if ipProxy := req.Header.Get("X-FORWARDED-FOR"); len(ipProxy) > 0 {
		ips := strings.Split(ipProxy, ", ") // Check for double forwarding ex. CloudFlare
		if len(ips) > 1 {
			return ips[0]
		} else {
			return ipProxy
		}
	}
	ip, _, _ := net.SplitHostPort(req.RemoteAddr)
	return ip
}

func GetSessionId(res http.ResponseWriter, req *http.Request) (sessionid string, err error) {
	cookie, err := req.Cookie("sessionid")
	if err != nil {
		log.Println("Failed to retrieve sessionid", err)
		Logout(res, req)
		return
	}
	sessionid = cookie.Value
	return
}

func generateSecret(password string) (hash string, err error) {
	hashBytes, err := scrypt.Key([]byte(password), secret, 16384, 8, 1, 32)
	hash = hex.EncodeToString(hashBytes)
	return
}

func GenerateSecret(email string) (hash string, date string, err error) {
	date = time.Now().Format(time.RFC3339)
	hash, err = generateSecret(email + date)
	if err != nil {
		log.Println("Error generating secret:", err)
	}
	return
}

func CheckSecret(hash, email, date string) bool {
	log.Printf("Checking secret: (%s, %s) against hash: %s\n", email, date, hash)
	checkDate, err := time.Parse(time.RFC3339, date)

	if err != nil {
		log.Println("Error parsing date:", err)
	}
	if time.Now().Sub(checkDate) > linkAgeLimit {
		log.Printf("Error date (%s) is more than 14 days old.\n", date)
	}
	if err != nil || time.Now().Sub(checkDate) > linkAgeLimit {
		return false
	}

	newHash, err := generateSecret(email + date)
	if err != nil {
		log.Println("Error generating secret:", err)
	}
	if newHash != hash {
		log.Printf("Error: string (%s%s) does not match hash (%s).\n", email, date, hash)
	}
	return err == nil && newHash == hash
}

func CheckParam(values url.Values, key string) bool {
	if values[key] != nil || len(values[key]) == 1 {
		return true
	}
	return false
}
