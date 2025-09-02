package logic

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"uni-token-service/constants"
	openBrowser "uni-token-service/logic/open_browser"
	"uni-token-service/store"

	"github.com/google/uuid"
)

var openingUI = make(map[string](chan struct{}))
var ServerPort = -1

func OpenUI(params url.Values, auth bool) error {
	if auth {
		allUsers, err := store.Users.List()
		if err != nil {
			return err
		}

		var userName string
		if len(allUsers) == 0 {
			userName = constants.UserName
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

	sessionId := uuid.New().String()
	params.Set("session", sessionId)
	params.Set("port", strconv.Itoa(ServerPort))

	err := openBrowser.OpenBrowser(constants.UserName, constants.AppBaseUrl+"?"+params.Encode())
	if err != nil {
		return err
	}

	channel := make(chan struct{})
	openingUI[sessionId] = channel
	defer delete(openingUI, sessionId)

	select {
	case <-channel:
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("failed to open UI: timeout")
	}
}

func OnUIOpened(sessionId string) {
	channel, ok := openingUI[sessionId]
	if !ok {
		return
	}
	channel <- struct{}{}
	close(channel)
}
