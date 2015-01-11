package database

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"github.com/comforme/comforme/common"
)

// Errors
var EmailInUse = errors.New("You have already registered with this email address.")
var UsernameInUse = errors.New("This username is in use. Please select a different one.")
var InvalidUsernameOrPassword = errors.New("Invalid username or password.")

type DB struct {
	conn *sql.DB
}

func NewDB(constr string) (DB, error) {
	conn, err := connect(constr)
	if err != nil {
		return DB{}, err
	}
	return DB{conn}, nil
}

// Establishes connection to Postgres database
func connect(constr string) (*sql.DB, error) {
	conn, err := sql.Open("postgres", constr)
	return conn, err
}

func (db DB) GetUserID(email string, password string) (userid int, err error) {
	log.Printf("Looking up user: %s\n", email)

	// Get hashed password
	var hashed string
	err = db.conn.QueryRow("SELECT id, password FROM users WHERE email = $1", email).Scan(&userid, &hashed)
	if err != nil {
		return
	}

	// Check hashed password
	err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return
}

func (db DB) GetSessionUserID(sessionid string) (userid int, err error) {
	log.Printf("Looking up userid for sessionid: %s\n", sessionid)

	err = db.conn.QueryRow("SELECT userid FROM sessions WHERE sessions.id = $1", sessionid).Scan(&userid)
	return
}

func (db DB) NewSession(userid int) (sessionid string, err error) {
	// Create a new unique sessionid
	for numRows := 0; ; {
		sessionid = common.NewSessionID()
		err = db.conn.QueryRow("SELECT count(*) FROM sessions WHERE sessions.id = $1", sessionid).Scan(&numRows)
		if err != nil {
			return
		}
		if numRows == 0 {
			break
		}
	}

	// Insert new sessionid
	_, err = db.conn.Exec(
		"INSERT INTO sessions (userid, id) VALUES ($1, $2)",
		userid,
		sessionid,
	)

	return
}

func (db DB) RegisterUser(username, email string) (password string, err error) {
	// Check if email is already in use
	var numRows int
	err = db.conn.QueryRow("SELECT count(*) FROM users WHERE email = $1", email).Scan(&numRows)
	if err != nil {
		return
	}
	if numRows != 0 {
		err = EmailInUse
		return
	}

	// Check if username is already in use
	err = db.conn.QueryRow("SELECT count(*) FROM users WHERE username = $1", username).Scan(&numRows)
	if err != nil {
		return
	}
	if numRows != 0 {
		err = UsernameInUse
		return
	}

	// Add new user
	password = common.GenPassword()
	hashed, err := hashPassword(password)
	if err != nil {
		return
	}
	log.Printf(
		"Adding user: %s with hashed password: %s\n",
		email,
		hashed,
	)
	_, err = db.conn.Exec(
		"INSERT INTO users (email, username, password) VALUES ($1, $2, $3)",
		email,
		username,
		hashed,
	)
	return
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

func (db DB) GetEmail(sessionid string) (email string, err error) {
	err = db.conn.QueryRow(
		"SELECT email FROM sessions, users WHERE sessions.id = $1 AND sessions.userid = users.id",
		sessionid,
	).Scan(&email)
	return
}

func (db DB) PasswordChangeRequired(sessionid string) (isRequired bool, err error) {
	err = db.conn.QueryRow(
		"SELECT reset_required FROM sessions, users WHERE sessions.id = $1 AND sessions.userid = users.id",
		sessionid,
	).Scan(&isRequired)
	return
}

func (db DB) ResetPassword(email string) (password string, err error) {
	password = common.GenPassword()
	hashed, err := hashPassword(password)
	if err != nil {
		return
	}

	log.Printf(
		"Resetting password for user: %s with hashed password: %s\n",
		email,
		hashed,
	)
	err = db.changePassword(email, hashed)
	return
}

func (db DB) changePassword(email, hashed string) error {
	result, err := db.conn.Exec(
		"UPDATE users SET password = $2, reset_required = false WHERE email = $1;",
		email,
		hashed,
	)

	if err != nil {
		return err
	}

	return checkSingleRow(result)
}

func (db DB) ChangePassword(email, newPassword string) error {
	hashed, err := hashPassword(newPassword)
	if err != nil {
		return err
	}

	log.Printf(
		"Changing password for user: %s with hashed password: %s\n",
		email,
		hashed,
	)
	return db.changePassword(email, hashed)
}

func (db DB) Logout(sessionid string) error {
	log.Println(
		"Logging out sessionid:",
		sessionid,
	)

	result, err := db.conn.Exec(
		"DELETE FROM sessions WHERE sessions.id = $1;",
		sessionid,
	)
	if err != nil {
		return err
	}

	return checkSingleRow(result)
}

func checkSingleRow(result sql.Result) error {
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return sql.ErrNoRows
	}

	return nil
}
