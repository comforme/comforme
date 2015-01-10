package databaseActions

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/keighl/mandrill"

	"github.com/comforme/comforme/common"
	"github.com/comforme/comforme/database"
)

// Errors
var InvalidEmail = errors.New("The provided email address is not valid.")
var PasswordTooShort = errors.New(fmp.Sprintf("The supplied password is too short. Minimum password length is %d characters.", minPasswordLength))
var EmailFailed = errors.New("Sending email failed.")
var InvalidSessionID = errors.New("Invalid sessionid.")
var InvalidPassword = errors.New("Invalid password.")
var DatabaseError = errors.New("Unknown database error.")

const (
	minPasswordLength = 6
)

var emailRegex *regexp.Regexp
var db database.DB

func init() {
	db = database.NewDB(os.Getenv("DATABASE_URL"))
	emailRegex = regexp.MustCompile("^.+@.+\\..+$")
}

func ResetPassword(username string) error {
	password, err := db.ResetPassword(username)
	if err != nil {
		return err
	}
	return sendResetEmail(email, password)
}

func ChangePassword(sessionid, oldPassword, newPassword string) error {
	log.Printf("Looking up email with sessionid: %s\n", sessionid)

	// Get email from session
	email, err := db.GetEmail(sessionid)
	if err != nil {
		log.Printf("Error retrieving email from sessionid (%s): %s\n", sessionid, err.Error())
		return InvalidSessionID
	}
	log.Printf("Sessionid: %s is associated with the email: %s\n", sessionid, email)

	// Check old password
	_, err = db.GetUserID(email, oldPassword)
	if err != nil {
		log.Printf("Error validating old password while changing password for user (%s): %s\n", email, err.Error())
		return InvalidPassword
	}

	// Check new password meets requirements
	if len(newPassword) < minPasswordLength {
		log.Printf(
			"New password for user %s of length %d is too short. %d required.\n",
			email,
			len(newPassword),
			minPasswordLength,
		)
		return InvalidPassword
	}

	return db.ChangePassword(email, newPassword)
}

func Logout(sessionid string) error {
	return db.Logout(sessionid)
}

func GetEmail(sessionid string) (email string, err error) {
	return db.GetEmail(sessionid)
}

func Login(email string, password string) (sessionid string, err error) {
	userid, err := db.GetUserID(username, password)
	if err != nil {
		log.Printf("Error while logging in user (%s): %s\n", username, err.Error())
		err = InvalidUsernameOrPassword
		return
	}

	sessionid, err = db.NewSession(userid)
	if err != nil {
		log.Printf("Error while creating session for user (%s): %s\n", username, err.Error())
		err = InvalidUsernameOrPassword
		return
	}

	return
}

func Register(username, email string) (sessionid string, err error) {
	password, err := db.RegisterUser(username, email)
	if err != nil {
		return
	}

	err = common.SendRegEmail(email, password)
	if err != nil {
		return
	}

	return Login(email, password)
}
