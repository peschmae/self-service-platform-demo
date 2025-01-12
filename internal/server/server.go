package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int
}

func getPort() int {
	port, set := os.LookupEnv("PORT")
	if !set {
		return 8080
	}
	p, _ := strconv.Atoi(port)
	return p
}

func NewServer() *http.Server {
	NewServer := &Server{
		port: getPort(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
