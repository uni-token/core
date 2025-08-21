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

	// Setup routes
	SetupActionAPI(router)
	SetupGatewayAPI(router)
	SetupAppAPI(router)
	SetupProvidersAPI(router)
	SetupPresetsAPI(router)
	SetupUsageAPI(router)
	SetupSiliconFlowAPI(router)
	SetupAuthAPI(router)

	// Get a random available port
	port := getPort()

	// Start server in a goroutine
	go func() {
		router.Run(":" + strconv.Itoa(port))
	}()

	return port, nil
}

func checkPortAvailability(port int) bool {
	_, err := net.Dial("tcp", "localhost:"+strconv.Itoa(port))
	return err != nil
}

const PORT_RANGE_START = 18760

func getPort() int {
	for i := PORT_RANGE_START; i < PORT_RANGE_START+10; i++ {
		if checkPortAvailability(i) {
			return i
		} else {
			log.Printf("Port %d is not available", i)
		}
	}
	panic(fmt.Sprintf("No available port found in the range %d-%d", PORT_RANGE_START, PORT_RANGE_START+10))
}
