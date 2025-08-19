package logic

import "strconv"

// -X 'logic.version=$(date +'%Y%m%d')'
var version string
// -X 'logic.appBaseUrl=http://uni-token.app'
var appBaseUrl string

func GetVersion() int {
	if version == "" {
		return 0
	}
	v, err := strconv.Atoi(version)
	if err != nil {
		return 0
	}
	return v
}

func GetAppBaseUrl() string {
	if appBaseUrl == "" {
		return "http://localhost:5173"
	}
	return appBaseUrl
}
