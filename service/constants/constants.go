package constants

import (
	"os"
	"os/user"
	"strconv"
)

// -X 'logic.version=$(date +'%Y%m%d')'
var version string

// -X 'logic.appBaseUrl=http://uni-token.app'
var appBaseUrl string

var Version = func() int {
	v, err := strconv.Atoi(version)
	if err != nil {
		return 0
	}
	return v
}()

var AppBaseUrl = func() string {
	if appBaseUrl == "" {
		return "http://localhost:5173"
	}
	return appBaseUrl
}()

var UserName = func() string {
	if userName := os.Getenv("UNI_TOKEN_SERVICE_USER"); userName != "" {
		return userName
	}
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	return currentUser.Username
}()

var ShouldChangeUser = func() bool {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	return currentUser.Username != UserName
}()
