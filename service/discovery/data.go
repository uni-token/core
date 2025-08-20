package discovery

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type ServiceInfo struct {
	Command   []string `json:"command"`
	PID       int      `json:"pid"`
	URL       string   `json:"url"`
	Timestamp int64    `json:"timestamp"`
}

func GetServiceInfo(port *int) string {
	var url string
	if port != nil {
		url = fmt.Sprintf("http://localhost:%d/", *port)
	}

	service := ServiceInfo{
		Command:   os.Args,
		PID:       os.Getpid(),
		URL:       url,
		Timestamp: time.Now().UnixMilli(),
	}

	jsonData, _ := json.MarshalIndent(service, "", "  ")
	return string(jsonData)
}
