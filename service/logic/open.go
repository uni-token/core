package logic

import (
	"net/url"

	"github.com/pkg/browser"

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
			userName = "guest"
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

	return browser.OpenURL(GetAppBaseUrl() + "?" + params.Encode())
}
