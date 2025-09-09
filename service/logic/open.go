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

var ServerPort = -1

var sessionActive = make(map[string]chan<- struct{})

func OpenAction(actionType string, params url.Values) (<-chan struct{}, func(), error) {
	return OpenUI("/action/"+actionType, params, true)
}

func OpenUI(path string, params url.Values, auth bool) (<-chan struct{}, func(), error) {
	if auth {
		allUsers, err := store.Users.List()
		if err != nil {
			return nil, nil, err
		}

		var userName string
		if len(allUsers) == 0 {
			userName = constants.UserName
		} else {
			userName = allUsers[0].Username
		}

		token, err := GenerateJWT(userName)
		if err != nil {
			return nil, nil, err
		}
		params.Set("username", userName)
		params.Set("token", token)
	}

	sessionId := uuid.New().String()
	params.Set("session", sessionId)
	params.Set("port", strconv.Itoa(ServerPort))

	openBrowser.OpenBrowser(constants.UserName, constants.AppBaseUrl+path+"?"+params.Encode())

	channel := make(chan struct{})
	sessionActive[sessionId] = channel
	cleanup := func() {
		delete(sessionActive, sessionId)
	}

	select {
	case <-channel:
		return channel, cleanup, nil
	case <-time.After(5 * time.Second):
		cleanup()
		return nil, nil, fmt.Errorf("failed to open UI: timeout")
	}
}

func OnUIActive(sessionId string) bool {
	channel, ok := sessionActive[sessionId]
	if ok {
		channel <- struct{}{}
	}
	return ok
}
