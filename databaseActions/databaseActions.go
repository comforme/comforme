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
var UsernameTooLong = errors.New(fmt.Sprintf("The supplied username is too long. Maximum username length is %d characters.", maxUsernameLength))
var EmailFailed = errors.New("Sending email failed.")
var IncorrectPassword = errors.New("Incorrect password.")
var ShortPassword = errors.New("Password too short.")

const (
	minPasswordLength = 6
	minUsernameLength = 3
	maxUsernameLength = 20
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

func CreatePage(sessionId, title, description, address, website string, category int) (categorySlug, pageSlug string, err error) {
	// TODO: Resolve location from address and update lower-level function to accept point
	slug := common.GenSlug(title)
	if len(slug) <= 1 {
		err = common.InvalidTitle
		return
	}

	pageID, err := db.NewPage(sessionId, title, slug, description, address, website, category)
	if err != nil {
		log.Println("Failed to create page", title)
		return
	}

	categorySlug, pageSlug, err = db.GetSlugs(pageID)
	if err != nil {
		return
	}

	return
}

func CreatePost(sessionId, post string, page common.Page) (err error) {
	user_id, err := db.GetSessionUserID(sessionId)
	if err != nil {
		log.Printf("Error getting userid from sessionid %s: %s\n", sessionId, err.Error())
		return
	}

	err = db.NewPost(user_id, page.Id, post)
	if err != nil {
		return
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
		return ShortPassword
	}

	return db.ChangePassword(email, newPassword)
}

func SetPassword(email, newPassword string) (err error) {
	// Check new password meets requirements
	if len(newPassword) < minPasswordLength {
		log.Printf(
			"New password for user %s of length %d is too short. %d required.\n",
			email,
			len(newPassword),
			minPasswordLength,
		)
		return ShortPassword
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
		err = IncorrectPassword
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

func Register2(username, email, password string) (sessionid string, err error) {
	if !common.ValidEmail(email) {
		err = common.InvalidEmail
		return
	}

	if len(username) < minUsernameLength {
		err = UsernameTooShort
		return
	}

	if len(username) > maxUsernameLength {
		err = UsernameTooLong
		return
	}

	// Check new password meets requirements
	if len(password) < minPasswordLength {
		log.Printf(
			"New password for user %s of length %d is too short. %d required.\n",
			email,
			len(password),
			minPasswordLength,
		)
		err = ShortPassword
		return
	}

	err = db.RegisterUser(username, email, password)
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

func Register1(email string) (err error) {
	if !common.ValidEmail(email) {
		err = common.InvalidEmail
		return
	}

	err = db.CheckEmailInUse(email)
	if err != nil {
		return
	}

	err = common.SendRegEmail(email)
	if err != nil {
		return
	}

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

func GetPage(categorySlug, pageSlug string) (page common.Page, err error) {
	page, err = db.GetPage(categorySlug, pageSlug)
	if err != nil {
		log.Printf("Error looking up page with category (%s) and slug (%s): %s\n", categorySlug, pageSlug, err.Error())
		return
	}

	return
}

func GetPosts(sessionid string, page common.Page) (posts []common.Post, err error) {
	user_id, err := db.GetSessionUserID(sessionid)
	if err != nil {
		log.Printf("Error getting userid from sessionid %s: %s\n", sessionid, err.Error())
		return
	}

	posts, err = db.GetPostsForPage(user_id, page.Id)
	if err != nil {
		log.Printf("Error looking up posts for page (%d): %s\n", page.Id, err.Error())
		return
	}

	return
}

func GetCommunityColumns(sessionid string) ([][]common.Community, error) {
	communities, err := ListCommunities(sessionid)
	if err != nil {
		return [][]common.Community{}, err
	}

	perCol := len(communities) / 4
	extra := len(communities) % 4
	cut1 := perCol
	if extra >= 1 {
		cut1++
	}
	cut2 := cut1 + perCol
	if extra >= 2 {
		cut2++
	}
	cut3 := cut2 + perCol
	if extra >= 3 {
		cut3++
	}
	return [][]common.Community{
		communities[0:cut1],
		communities[cut1:cut2],
		communities[cut2:cut3],
		communities[cut3:],
	}, nil
}

func CheckResetLink(code, email, date string) bool {
	password, err := db.GetPasswordHash(email)
	if err != nil {
		return false
	}

	return common.CheckSecret(code, email+password, date)
}

func CheckRegisterLink(code, email, date string) bool {
	err := db.CheckEmailInUse(email)
	if err != nil {
		return false
	}

	return common.CheckSecret(code, email, date)
}
