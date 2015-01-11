package static

import(
	"fmt"
	"net/http"
)

func Style(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Fprintln(res, `/* ---------- Sign Up/Log In ---------- */

.sign-up-and-log-in.content {
	margin-top: 4rem;
}

.login-tabs {
	border: 1px solid #CCCCCC;
}

.login-tabs dd {
	width: 50%;
	text-align: center;
}
.login-tabs dd.active a {
	background-color: #C8C8C8;
}

.tabs-content {
	border-top: 1px solid #CCCCCC;
	padding: 0 1rem 0;
}
.tabs-content .content {
	float: none;
	padding-bottom: 0;
}
.tabs-content, .tabs-content form {
	margin: 0;
}


/* ---------- Communities Settings ---------- */
.communities-settings section {
	margin-bottom: 1rem;
}

.communities-settings section form {
	margin-bottom: 0;
}`)
}
