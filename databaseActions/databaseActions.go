package databaseActions

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/comforme/comforme/common"
	"github.com/comforme/comforme/database"
)

// Errors
var PasswordTooShort = errors.New(fmt.Sprintf("The supplied password is too short. Minimum password length is %d characters.", minPasswordLength))
var UsernameTooShort = errors.New(fmt.Sprintf("The supplied username is too short. Minimum username length is %d characters.", minUsernameLength))
var EmailFailed = errors.New("Sending email failed.")
var InvalidPassword = errors.New("Invalid password.")

const (
	minPasswordLength = 6
	minUsernameLength = 3
)

var db database.DB

func init() {
	var err error
	db, err = database.NewDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
}

func ResetPassword(email string) error {
	password, err := db.ResetPassword(email)
	if err != nil {
		return err
	}
	return common.SendResetEmail(email, password)
}

func CreatePage(sessionId string, title string, description string, address string, category int) (err error) {
	// TODO Resolve location from address and update lower-level function to accept point
	err = db.NewPage(sessionId, title, description, address, category)
	if err != nil {
		log.Println("Failed to create page", title)
	}
	return
}

func ChangePassword(sessionid, oldPassword, newPassword string) (err error) {
	email, err := CheckPassword(sessionid, oldPassword)
	if err != nil {
		return
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

func GetEmail(sessionid string) (string, error) {
	return db.GetEmail(sessionid)
}

func GetUsername(sessionid string) (string, error) {
	return db.GetUsername(sessionid)
}

func PasswordChangeRequired(sessionid string) (bool, error) {
	return db.PasswordChangeRequired(sessionid)
}

func ListCategories() (map[string]string, error) {
	return db.ListCategories()
}

func Login(email string, password string) (sessionid string, err error) {
	userid, err := db.GetUserID(email, password)
	if err != nil {
		log.Printf("Error while logging in user (%s): %s\n", email, err.Error())
		return
	}

	sessionid, err = db.NewSession(userid)
	if err != nil {
		return
	}

	return
}

func CheckPassword(sessionid, password string) (email string, err error) {
	log.Printf("Looking up email with sessionid: %s\n", sessionid)

	// Get email from session
	email, err = db.GetEmail(sessionid)
	if err != nil {
		return
	}
	log.Printf("Sessionid: %s is associated with the email: %s\n", sessionid, email)

	// Check old password
	_, err = db.GetUserID(email, password)
	if err != nil {
		log.Printf("Error validating old password while changing password for user (%s): %s\n", email, err.Error())
		err = InvalidPassword
		return
	}

	return
}

func ChangeUsername(sessionid, newUsername, password string) (err error) {
	if len(newUsername) < minUsernameLength {
		err = UsernameTooShort
		return
	}

	_, err = CheckPassword(sessionid, password)
	if err != nil {
		return
	}

	user_id, err := db.GetSessionUserID(sessionid)
	if err != nil {
		log.Printf("Error getting userid from sessionid %s: %s\n", sessionid, err.Error())
		return
	}

	err = db.ChangeUsername(user_id, newUsername)

	return
}

func Register(username, email string) (sessionid string, err error) {
	if !common.ValidEmail(email) {
		err = common.InvalidEmail
		return
	}

	if len(username) < minUsernameLength {
		err = UsernameTooShort
		return
	}

	password, err := db.RegisterUser(username, email)
	if err != nil {
		return
	}

	err = common.SendRegEmail(email, password)
	if err != nil {
		return
	}

	sessionid, err = Login(email, password)
	if err != nil {
		return
	}

	// Make new users lazy :)
	err = SetCommunityMembership(sessionid, 1, true)

	return
}

func ListCommunities(sessionid string) (communities []common.Community, err error) {
	communities, err = db.ListCommunities(sessionid)
	return
}

func SetCommunityMembership(sessionid string, community_id int, value bool) (err error) {
	user_id, err := db.GetSessionUserID(sessionid)
	if err != nil {
		log.Printf("Error getting userid from sessionid %s: %s\n", sessionid, err.Error())
		return
	}

	if value {
		err = db.AddCommunityMembership(user_id, community_id)
		if err != nil {
			return
		}
	} else {
		err = db.DeleteCommunityMembership(user_id, community_id)
		if err != nil {
			return
		}
	}

	return
}

func OtherSessions(sessionid string) (num int, err error) {
	user_id, err := db.GetSessionUserID(sessionid)
	if err != nil {
		log.Printf("Error getting userid from sessionid %s: %s\n", sessionid, err.Error())
		return
	}

	num, err = db.OpenSessions(user_id)
	num--
	return
}

func LogoutOtherSessions(sessionid string) (loggedOut int, err error) {
	user_id, err := db.GetSessionUserID(sessionid)
	if err != nil {
		log.Printf("Error getting userid from sessionid %s: %s\n", sessionid, err.Error())
		return
	}

	loggedOut, err = db.DeleteOtherSessions(user_id, sessionid)
	if err != nil {
		log.Printf(
			"Error deleting other sessions for userid (%d) with sessionid (%s): %s\n",
			user_id,
			sessionid,
			err.Error(),
		)
		return
	}

	return
}

func SearchPages(sessionid, query string) ([]common.Page, error) {
	return db.SearchPages(query)
}

func GetPage(sessionid, category, slug string) (page common.Page, posts []common.Post, err error) {
	user_id, err := db.GetSessionUserID(sessionid)
	if err != nil {
		log.Printf("Error getting userid from sessionid %s: %s\n", sessionid, err.Error())
		return
	}

	page, err = db.GetPage(category, slug)
	if err != nil {
		log.Printf("Error looking up page with category (%s) and slug (%s): %s\n", category, slug, err.Error())
		return
	}

		posts, err = db.GetPostsForPage(user_id, page.Id)
		if err != nil {
			log.Printf("Error looking up posts for page (%s) with category (%s) and slug (%s): %s\n", page.Title, category, slug, err.Error())
			return
		}

	return
}
