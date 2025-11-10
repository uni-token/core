package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
	"uni-token-service/logic"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupAPIServer() (int, error) {
	logic.InitJWTSecret()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:        []string{"http://localhost:*", "https://uni-token.app"},
		AllowWildcard:       true,
		AllowMethods:        []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:        []string{"*"},
		ExposeHeaders:       []string{"*"},
		AllowCredentials:    true,
		AllowPrivateNetwork: true,
		AllowWebSockets:     true,
		MaxAge:              240 * time.Hour,
	}))

	setupRoutes(router)

	port := findAvailablePort()
	logic.ServerPort = port

	go router.Run(":" + strconv.Itoa(port))
	return port, nil
}

func setupRoutes(router *gin.Engine) {
	SetupActionAPI(router)
	SetupGatewayAPI(router)
	SetupAppAPI(router)
	SetupKeysAPI(router)
	SetupPresetsAPI(router)
	SetupUsageAPI(router)
	SetupAuthAPI(router)
	SetupProxyAPI(router)
	SetupStoreAPI(router)
}

func isPortAvailable(port int) bool {
	_, err := net.Dial("tcp", "localhost:"+strconv.Itoa(port))
	return err != nil
}

const portRangeStart = 18760

func findAvailablePort() int {
	for port := portRangeStart; port < portRangeStart+10; port++ {
		if isPortAvailable(port) {
			return port
		}
		log.Printf("Port %d is not available", port)
	}
	panic(fmt.Sprintf("No available port found in range %d-%d", portRangeStart, portRangeStart+10))
}
