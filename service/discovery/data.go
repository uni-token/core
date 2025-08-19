package discovery

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type ServiceData struct {
	Command   []string `json:"command"`
	PID       int      `json:"pid"`
	URL       *string  `json:"url"`
	Timestamp int64    `json:"timestamp"`
}

func GetData(port *int) string {
	var url *string
	if port != nil {
		urlStr := fmt.Sprintf("http://localhost:%d/", *port)
		url = &urlStr
	}

	data := ServiceData{
		Command:   os.Args,
		PID:       os.Getpid(),
		URL:       url,
		Timestamp: time.Now().UnixMilli(),
	}

	jsonData, _ := json.MarshalIndent(data, "", "  ")
	return string(jsonData)
}
