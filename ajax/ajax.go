package ajax

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-zoo/bone"

	"github.com/comforme/comforme/databaseActions"
)

const (
	JSONLoginError  = "{ \"error\": \"Not logged in.\" }"
	JSONActionError = "{ \"error\": \"Invalid action.\" }"
	JSONError       = "{ \"error\": \"Unknown error.\" }"
)

type AjaxResultNum struct {
	Message string `json:"message"`
	Number  int    `json:"number"`
}

type AjaxResult struct {
	Message string `json:"message"`
}

type AjaxError struct {
	Message string `json:"error"`
}

func HandleAction(res http.ResponseWriter, req *http.Request, sessionid, email, username string, userID int) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	action := bone.GetValue(req, "action")
	var result interface{}

	if action == "addCommunity" || action == "removeCommunity" {
		community_id, err := strconv.ParseInt(req.PostFormValue("communityid"), 10, 0)
		if err != nil {
			log.Println("Error parsing communityid:", err)
			result = AjaxError{"Invalid communityid."}
		} else {
			err = databaseActions.SetCommunityMembership(userID, int(community_id), action == "addCommunity")
			if err != nil {
				result = AjaxError{err.Error()}
			} else {
				result = AjaxResult{fmt.Sprintf("Successfully set membership in community %d to %t.", community_id, action == "addCommunity")}
			}
		}
	} else if action == "logoutOtherSessions" {
		loggedOut, err := databaseActions.LogoutOtherSessions(sessionid, userID)
		if err != nil {
			result = AjaxError{err.Error()}
		} else {
			result = AjaxResultNum{
				fmt.Sprintf(
					"Successfully logged out %d other sessions for user with sessionid %s.",
					loggedOut,
					sessionid,
				),
				loggedOut,
			}
		}
	} else {
		fmt.Fprintln(res, JSONActionError)
		return
	}

	encoded, err := json.Marshal(result)
	if err != nil {
		log.Println("Error marshaling response:", err)
		fmt.Fprintln(res, JSONError)
		return
	}

	log.Println("AJAX result:", string(encoded))

	fmt.Fprintln(res, string(encoded))
	return
}
