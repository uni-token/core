package logic

import (
	"net/url"

	"uni-token-service/logic/open_browser"
	"uni-token-service/store"
)

func OpenUI(params url.Values, auth bool) error {
	if auth {
		allUsers, err := store.Users.List()
		if err != nil {
			return err
		}

		var userName string
		if len(allUsers) == 0 {
			userName = GetUserName()
		} else {
			userName = allUsers[0].Username
		}

		token, err := GenerateJWT(userName)
		if err != nil {
			return err
		}
		params.Set("username", userName)
		params.Set("token", token)
	}

	return openBrowser.OpenBrowser(GetUserName(), GetAppBaseUrl()+"?"+params.Encode())
}
