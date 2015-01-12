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

type AxaxResult struct {
	Message string `json:"message"`
}

type AxaxError struct {
	Message string `json:"error"`
}

func HandleAction(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	cookie, err := req.Cookie("sessionid")
	if err == nil {
		fmt.Fprintln(res, JSONLoginError)
		return
	}

	action := bone.GetValue(req, "action")
	var result interface{}

	if action == "addCommunity" || action == "removeCommunity" {
		community_id, err := strconv.ParseInt(req.PostFormValue("communityid"), 10, 0)
		if err != nil {
			log.Println("Error parsing communityid:", err)
			result = AxaxError{"Invalid communityid."}
		} else {
			err = databaseActions.SetCommunityMembership(cookie.Value, int(community_id), action == "addCommunity")
			if err != nil {
				result = AxaxError{err.Error()}
			}

			result = AxaxResult{fmt.Sprintf("Successfully set membership in community %d to %t.", community_id, action == "addCommunity")}
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
