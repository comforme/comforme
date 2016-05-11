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

func ResetPassword(email, baseURL string) error {
	hash, date, err := GenerateResetCode(email)
	if err != nil {
		return err
	}
	return common.SendResetEmail(email, date, hash, baseURL)
}

func CreatePage(userID int, title, description, address, website string, category int) (categorySlug, pageSlug string, err error) {
	// TODO: Resolve location from address and update lower-level function to accept point
	slug := common.GenSlug(title)
	if len(slug) <= 1 {
		err = common.InvalidTitle
		return
	}

	pageID, err := db.NewPage(userID, title, slug, description, address, website, category)
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

func CreatePost(user_id int, post string, page common.Page) (err error) {
	err = db.NewPost(user_id, page.Id, post)
	if err != nil {
		return
	}

	return
}

func ChangePassword(email, oldPassword, newPassword string) (err error) {
	_, err = db.GetUserID(email, oldPassword)
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

func GetUserInfo(sessionid string) (common.UserInfo, error) {
	return db.GetUserInfo(sessionid)
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

func Login(email, password string) (sessionid string, err error) {
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

func ChangeUsername(email string, newUsername, password string) (err error) {
	if len(newUsername) < minUsernameLength {
		err = UsernameTooShort
		return
	}

	userid, err := db.GetUserID(email, password)
	if err != nil {
		return
	}

	err = db.ChangeUsername(userid, newUsername)

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

	userid, err := db.GetSessionUserID(sessionid)
	if err != nil {
		return
	}

	// Make new users lazy :)
	err = SetCommunityMembership(userid, 1, true)
	return
}

func Register1(email, baseURL string) (err error) {
	if !common.ValidEmail(email) {
		err = common.InvalidEmail
		return
	}

	err = db.CheckEmailInUse(email)
	if err != nil {
		return
	}

	err = common.SendRegEmail(email, baseURL)
	if err != nil {
		return
	}

	return
}

func SetCommunityMembership(userid int, community_id int, value bool) (err error) {
	if value {
		err = db.AddCommunityMembership(userid, community_id)
		if err != nil {
			return
		}
	} else {
		err = db.DeleteCommunityMembership(userid, community_id)
		if err != nil {
			return
		}
	}

	return
}

func OtherSessions(userid int) (num int, err error) {
	num, err = db.OpenSessions(userid)
	num--
	return
}

func LogoutOtherSessions(sessionid string, userid int) (loggedOut int, err error) {
	loggedOut, err = db.DeleteOtherSessions(userid, sessionid)
	if err != nil {
		log.Printf(
			"Error deleting other sessions for userid (%d) with sessionid (%s): %s\n",
			userid,
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

func GetPages() ([]common.Page, error) {
	return db.GetPages()
}

func GetPage(categorySlug, pageSlug string) (page common.Page, err error) {
	page, err = db.GetPage(categorySlug, pageSlug)
	if err != nil {
		log.Printf("Error looking up page with category (%s) and slug (%s): %s\n", categorySlug, pageSlug, err.Error())
		return
	}

	return
}

func GetPosts(userid int, page common.Page) (posts []common.Post, err error) {
	posts, err = db.GetPostsForPage(userid, page.Id)
	if err != nil {
		log.Printf("Error looking up posts for page (%d): %s\n", page.Id, err.Error())
		return
	}

	return
}

func GetCommunityColumns(userid int) ([][]common.Community, error) {
	communities, err := db.ListCommunities(userid)
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

func GenerateResetCode(email string) (hash string, date string, err error) {
	password, err := db.GetPasswordHash(email)
	if err != nil {
		return
	}

	return common.GenerateSecret(email + password)
}

func CheckRegisterLink(code, email, date string) bool {
	err := db.CheckEmailInUse(email)
	if err != nil {
		return false
	}

	return common.CheckSecret(code, email, date)
}

func GetTopPages() (pages []common.PagePostCount, err error) {
	return db.GetTopPages()
}
