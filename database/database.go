package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"github.com/comforme/comforme/common"
)

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

func (db DB) GetPasswordHash(email string) (hash string, err error) {
	log.Printf("Looking up user: %s\n", email)

	// Get hashed password
	err = db.conn.QueryRow("SELECT password FROM users WHERE email = $1", email).Scan(&hash)
	if err != nil {
		common.LogError(err)
		log.Printf("Error retrieving hashed password for email (%s): %s\n", email, err.Error())
		err = common.InvalidUsernameOrPassword
		return
	}
	return
}

func (db DB) GetUserID(email string, password string) (userid int, err error) {
	log.Printf("Looking up user: %s\n", email)

	// Get hashed password
	var hashed string
	err = db.conn.QueryRow("SELECT id, password FROM users WHERE email = $1", email).Scan(&userid, &hashed)
	if err != nil {
		common.LogError(err)
		log.Printf("Error retrieving userid and hashed password for email (%s): %s\n", email, err.Error())
		err = common.InvalidUsernameOrPassword
		return
	}

	// Check hashed password
	err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		common.LogError(err)
		log.Printf("Error checking password for user with email (%s): %s\n", email, err.Error())
		err = common.InvalidUsernameOrPassword
		return
	}
	return
}

func (db DB) GetSessionUserID(sessionid string) (userid int, err error) {
	log.Printf("Looking up userid for sessionid: %s\n", sessionid)

	err = db.conn.QueryRow("SELECT user_id FROM sessions WHERE sessions.id = $1", sessionid).Scan(&userid)
	if err != nil {
		common.LogError(err)
		err = common.InvalidSessionID
	}
	return
}

func (db DB) NewSession(userid int) (sessionid string, err error) {
	// Create a new unique sessionid
	for numRows := 0; ; {
		sessionid = common.NewSessionID()
		err = db.conn.QueryRow("SELECT count(*) FROM sessions WHERE sessions.id = $1", sessionid).Scan(&numRows)
		if err != nil {
			log.Printf("Error while creating session for userid (%d): %s\n", userid, err.Error())
			err = common.DatabaseError
			return
		}
		if numRows == 0 {
			break
		}
	}

	// Insert new sessionid
	_, err = db.conn.Exec(
		"INSERT INTO sessions (user_id, id) VALUES ($1, $2)",
		userid,
		sessionid,
	)
	if err != nil {
		log.Printf("Error while creating session for userid (%d): %s\n", userid, err.Error())
		err = common.DatabaseError
		return
	}

	return
}

func (db DB) NewPage(userID int, title, slug, description, address, website string, category int) (pageID int, err error) {
	// Insert new page
	err = db.conn.QueryRow(`
		INSERT INTO
			pages (
				title,
				description,
				address,
				category,
				slug,
				user_id,
				website
			)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
		`,
		title,
		description,
		address,
		category,
		slug,
		userID,
		website,
	).Scan(&pageID)
	if err != nil {
		log.Println("Failed to insert page: ", err)
		err = common.PageAlreadyExists
		return
	}

	return
}

func (db DB) GetSlugs(pageID int) (categorySlug, pageSlug string, err error) {
	err = db.conn.QueryRow(`
		SELECT
			categories.slug,
			pages.slug
		FROM
			pages,
			categories
		WHERE
			pages.category = categories.id AND
			pages.id = $1;
		`,
		pageID,
	).Scan(&categorySlug, &pageSlug)
	if err != nil {
		log.Println("Failed to insert page: ", err)
		err = common.DatabaseError
		return
	}

	return
}

func (db DB) CheckEmailInUse(email string) (err error) {
	var numRows int
	err = db.conn.QueryRow("SELECT count(*) FROM users WHERE email = $1", email).Scan(&numRows)
	if err != nil {
		log.Printf("Error checking if email (%s) already exists: %s\n", email, err.Error())
		err = common.DatabaseError
		return
	}
	if numRows != 0 {
		err = common.EmailInUse
		return
	}

	return
}

func (db DB) checkUsernameInUse(username string) (err error) {
	var numRows int
	err = db.conn.QueryRow("SELECT count(*) FROM users WHERE username = $1", username).Scan(&numRows)
	if err != nil {
		log.Printf("Error checking if username (%s) already exists: %s\n", username, err.Error())
		err = common.DatabaseError
		return
	}
	if numRows != 0 {
		err = common.UsernameInUse
		return
	}

	return
}

func (db DB) NewPost(userID, pageID int, post string) (err error) {
	_, err = db.conn.Exec(
		"INSERT INTO posts (user_id, page_id, body) VALUES ($1, $2, $3)",
		userID,
		pageID,
		post,
	)
	if err != nil {
		log.Printf("Error post (%s) to page (%d) with user (%d): %s\n", post, pageID, userID, err.Error())
		err = common.DatabaseError
		return
	}

	return
}

func (db DB) RegisterUser(username, email, password string) (err error) {
	err = db.CheckEmailInUse(email)
	if err != nil {
		return
	}

	err = db.checkUsernameInUse(username)
	if err != nil {
		return
	}

	// Add new user
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
		"INSERT INTO users (email, username, password, reset_required) VALUES ($1, $2, $3, false)",
		email,
		username,
		hashed,
	)
	if err != nil {
		log.Println("Error adding user:", err)
		err = common.DatabaseError
		return
	}

	return
}

func (db DB) ChangeUsername(user_id int, newUsername string) error {
	err := db.checkUsernameInUse(newUsername)
	if err != nil {
		return err
	}

	result, err := db.conn.Exec(
		"UPDATE users SET username = $2 WHERE id = $1;",
		user_id,
		newUsername,
	)

	if err != nil {
		common.LogError(err)
		return common.DatabaseError
	}

	return checkSingleRow(result, common.DatabaseError)
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return "", common.DatabaseError
	}
	return string(hashed), nil
}

func (db DB) GetEmail(sessionid string) (email string, err error) {
	err = db.conn.QueryRow(
		"SELECT email FROM sessions, users WHERE sessions.id = $1 AND sessions.user_id = users.id",
		sessionid,
	).Scan(&email)
	if err != nil {
		log.Printf("Error looking up email associated with sessionid  (%s): %s\n", sessionid, err.Error())
		err = common.InvalidSessionID
	}
	return
}

func (db DB) GetUserInfo(sessionid string) (userInfo common.UserInfo, err error) {
	userInfo.SessionID = sessionid
	err = db.conn.QueryRow(
		"SELECT email, username, user_id FROM sessions, users WHERE sessions.id = $1 AND sessions.user_id = users.id",
		sessionid,
	).Scan(&userInfo.Email, &userInfo.Username, &userInfo.UserID)
	if err != nil {
		log.Printf("Error looking up email and ID associated with sessionid  (%s): %s\n", sessionid, err.Error())
		err = common.InvalidSessionID
	}
	return
}

func (db DB) GetUsername(sessionid string) (username string, err error) {
	err = db.conn.QueryRow(
		"SELECT username FROM sessions, users WHERE sessions.id = $1 AND sessions.user_id = users.id",
		sessionid,
	).Scan(&username)
	if err != nil {
		log.Printf("Error looking up username associated with sessionid  (%s): %s\n", sessionid, err.Error())
		err = common.InvalidSessionID
	}
	return
}

func (db DB) PasswordChangeRequired(sessionid string) (isRequired bool, err error) {
	err = db.conn.QueryRow(
		"SELECT reset_required FROM sessions, users WHERE sessions.id = $1 AND sessions.user_id = users.id",
		sessionid,
	).Scan(&isRequired)
	if err != nil {
		log.Printf("Error checking if user with sessionid  (%s) needs to change their password: %s\n", sessionid, err.Error())
		err = common.InvalidSessionID
	}
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

	if err != nil {
		common.LogError(err)
		err = common.InvalidEmail
		return
	}

	err = db.requirePasswordReset(email)

	return
}

func (db DB) changePassword(email, hashed string) error {
	result, err := db.conn.Exec(
		"UPDATE users SET password = $2, reset_required = false WHERE email = $1;",
		email,
		hashed,
	)

	if err != nil {
		common.LogError(err)
		return common.DatabaseError
	}

	return checkSingleRow(result, common.InvalidEmail)
}

func (db DB) requirePasswordReset(email string) error {
	result, err := db.conn.Exec(
		"UPDATE users SET reset_required = true WHERE email = $1;",
		email,
	)

	if err != nil {
		common.LogError(err)
		return common.DatabaseError
	}

	return checkSingleRow(result, common.InvalidEmail)
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
		common.LogError(err)
		return common.DatabaseError
	}

	return checkSingleRow(result, common.InvalidSessionID)
}

func checkSingleRow(result sql.Result, otherwise error) error {
	log.Println("Checking single row...")
	rows, err := result.RowsAffected()
	if err != nil {
		common.LogErrorSkipLevels(err, 1)
		return common.DatabaseError
	}

	log.Printf("Query affected %d rows.", rows)
	if rows != 1 {
		return otherwise
	}

	return nil
}

func (db DB) ListCommunities(userid int) (communities []common.Community, err error) {
	rows, err := db.conn.Query(`
		SELECT
			communities.id,
			communities.name,
			my_memberships.member IS NOT NULL as is_member
		FROM
			communities
		LEFT JOIN
			(
				SELECT
					community_memberships.community_id,
					true as member
				FROM
					community_memberships
				WHERE
					community_memberships.user_id = $1
			) as my_memberships
				ON
					my_memberships.community_id = communities.id
		ORDER BY communities.id ASC;
		`,
		userid,
	)
	if err != nil {
		common.LogError(err)
		err = common.DatabaseError
		return
	}

	defer rows.Close()

	communities = []common.Community{}
	for rows.Next() {
		row := common.Community{}
		if err := rows.Scan(
			&row.Id,
			&row.Name,
			&row.IsMember,
		); err != nil {
			log.Println("Unknown iteration error:", err)
			return nil, err
		}
		communities = append(communities, row)
	}

	if err = rows.Err(); err != nil {
		common.LogError(err)
		err = common.DatabaseError
		return
	}

	// Success
	return
}

func (db DB) GetPages() (pages []common.Page, err error) {
	rows, err := db.conn.Query(`
		SELECT
			pages.id,
			title,
			pages.slug,
			categories.name,
			categories.slug,
			description,
			date_created
		FROM
			pages,
			categories
		WHERE
			categories.id = pages.category
		ORDER BY date_created DESC
		`)
	if err != nil {
		common.LogError(err)
		err = common.DatabaseError
		return
	}

	defer rows.Close()

	pages = []common.Page{}
	for rows.Next() {
		var row common.Page
		if err := rows.Scan(
			&row.Id,
			&row.Title,
			&row.PageSlug,
			&row.Category,
			&row.CategorySlug,
			&row.Description,
			&row.DateCreated,
		); err != nil {
			log.Fatal(err)
		}
		pages = append(pages, row)
	}

	if err = rows.Err(); err != nil {
		common.LogError(err)
		err = common.DatabaseError
		return
	}

	// Success
	return
}

func (db DB) SearchPages(query string) (pages []common.Page, err error) {
	rows, err := db.conn.Query(`
			SELECT
				pages.id,
				title,
				pages.slug,
				categories.name,
				categories.slug,
				description,
				date_created
			FROM
				pages,
				categories
			WHERE
				categories.id = pages.category AND
				to_tsvector('english', title) @@ plainto_tsquery($1) -- Full text search
			ORDER BY date_created DESC
			`,
		query,
	)
	if err != nil {
		common.LogError(err)
		err = common.DatabaseError
		return
	}

	defer rows.Close()

	pages = []common.Page{}
	for rows.Next() {
		var row common.Page
		if err := rows.Scan(
			&row.Id,
			&row.Title,
			&row.PageSlug,
			&row.Category,
			&row.CategorySlug,
			&row.Description,
			&row.DateCreated,
		); err != nil {
			log.Fatal(err)
		}
		pages = append(pages, row)
	}

	if err = rows.Err(); err != nil {
		common.LogError(err)
		err = common.DatabaseError
		return
	}

	// Success
	return
}

func (db DB) ListCategories() (categories map[string]string, err error) {
	rows, err := db.conn.Query(`
		SELECT
			id,
			name
		FROM
			categories
		`,
	)
	if err != nil {
		common.LogError(err)
		err = common.DatabaseError
		return
	}

	defer rows.Close()

	categories = map[string]string{}
	for rows.Next() {
		var category string
		var id string
		if err := rows.Scan(
			&id,
			&category,
		); err != nil {
			log.Fatal(err)
		}
		categories[id] = category
	}

	if err = rows.Err(); err != nil {
		common.LogError(err)
		err = common.DatabaseError
		return
	}

	// Success
	return
}

func (db DB) AddCommunityMembership(user_id, community_id int) (err error) {
	_, err = db.conn.Exec(
		"INSERT INTO community_memberships (user_id, community_id) VALUES ($1, $2)",
		user_id,
		community_id,
	)
	if err != nil {
		log.Println("Error adding community:", err)
		err = common.DatabaseError
		return
	}
	return
}

func (db DB) DeleteCommunityMembership(user_id, community_id int) (err error) {
	_, err = db.conn.Exec(
		"DELETE FROM community_memberships WHERE user_id = $1 AND community_id = $2;",
		user_id,
		community_id,
	)
	if err != nil {
		log.Println("Error deleting community:", err)
		err = common.DatabaseError
		return
	}
	return
}

func (db DB) OpenSessions(user_id int) (numSessions int, err error) {
	err = db.conn.QueryRow("SELECT count(*) FROM sessions WHERE user_id = $1;", user_id).Scan(&numSessions)
	if err != nil {
		log.Println("Error counting sessions:", err)
		err = common.DatabaseError
		return
	}
	return
}

func (db DB) DeleteOtherSessions(user_id int, sessionid string) (loggedOut int, err error) {
	result, err := db.conn.Exec(
		"DELETE FROM sessions WHERE user_id = $1 AND id <> $2;",
		user_id,
		sessionid,
	)
	if err != nil {
		log.Println("Error deleting other sessions:", err)
		err = common.InvalidSessionID
		return
	}

	loggedOutNum, err := result.RowsAffected()
	if err != nil {
		log.Println("Error deleted sessions:", err)
		err = common.DatabaseError
		return
	}
	loggedOut = int(loggedOutNum)
	return
}

func (db DB) GetPostsForPage(userid, pageid int) (posts []common.Post, err error) {
	rows, err := db.conn.Query(
		`
			SELECT
				posts.body,
				authors.username AS author,
				to_char(posts.date_created, 'YYYY-MM-DD HH24:MI:SS'),
				(
					SELECT count(*)
					FROM
						(
							SELECT community_id
							FROM
								community_memberships
							WHERE
								user_id = $1
						) my_communities,
						(
							SELECT community_id
							FROM
								community_memberships
							WHERE
								user_id = authors.id
						) author_communities
					WHERE
						my_communities.community_id = author_communities.community_id
				) AS common_communities
			FROM
				users authors,
				posts
			WHERE
				posts.user_id = authors.id
				AND posts.page_id = $2
			ORDER BY common_communities DESC;
		`,
		userid,
		pageid,
	)
	if err != nil {
		common.LogError(err)
		err = common.DatabaseError
		return
	}

	defer rows.Close()

	for rows.Next() {
		var row common.Post
		if err := rows.Scan(
			&row.Body,
			&row.Author,
			&row.Date,
			&row.CommonCategories,
		); err != nil {
			log.Fatal(err)
		}
		posts = append(posts, row)
	}

	if err = rows.Err(); err != nil {
		common.LogError(err)
		err = common.DatabaseError
		return
	}

	// Success
	return
}

func (db DB) GetPage(categorySlug, pageSlug string) (page common.Page, err error) {
	err = db.conn.QueryRow(`
		SELECT
			pages.id,
			title,
			pages.slug,
			categories.name,
			categories.slug,
			description,
			address,
			website,
			date_created
		FROM
			pages,
			categories
		WHERE
			categories.id = pages.category
			AND categories.slug = $1
			AND pages.slug = $2;`,
		categorySlug,
		pageSlug,
	).Scan(
		&page.Id,
		&page.Title,
		&page.PageSlug,
		&page.Category,
		&page.CategorySlug,
		&page.Description,
		&page.Address,
		&page.Website,
		&page.DateCreated,
	)
	if err != nil {
		log.Println("Error getting page:", err)
		err = common.PageNotFound
		return
	}
	return
}

func (db DB) GetTopPages() (pages []common.PagePostCount, err error) {
	rows, err := db.conn.Query(`
		SELECT
			top_pages.title,
			top_pages.slug AS page_slug,
			categories.name AS category_name,
			categories.slug AS category_slug,
			top_pages.count
		FROM
			categories,
			(
				SELECT
					pages.id,
					pages.title,
					pages.slug,
					pages.category,
					COUNT(*) as count
				FROM
					pages,
					posts
				WHERE
					pages.id = posts.page_id
				GROUP BY
					pages.id
				ORDER BY
					COUNT(*) DESC
				LIMIT 5
			) as top_pages
		WHERE
			categories.id = top_pages.category;`,
	)
	if err != nil {
		common.LogError(err)
		err = common.DatabaseError
		return
	}

	defer rows.Close()

	for rows.Next() {
		var row common.PagePostCount
		if err := rows.Scan(
			&row.Title,
			&row.PageSlug,
			&row.Category,
			&row.CategorySlug,
			&row.PostCount,
		); err != nil {
			log.Fatal(err)
		}
		pages = append(pages, row)
	}

	if err = rows.Err(); err != nil {
		common.LogError(err)
		err = common.DatabaseError
		return
	}

	// Success
	return
}
