package http_handlers

import (
	"agnos-backend/internal/config"
	healthcheckhdl "agnos-backend/internal/handlers/http/health_check"
	"agnos-backend/internal/middleware"
	"agnos-backend/internal/port"
	"log"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	config             *config.AppConfig
	server             *gin.Engine
	healthCheckHandler healthcheckhdl.IHealthCheckHandler
	staffHandler       port.IStaffHandler
	patientHandler     port.IPatientHandler
}

func NewHttpServer(
	config *config.AppConfig,
	server *gin.Engine,
	healthCheckHandler healthcheckhdl.IHealthCheckHandler,
	staffHandler port.IStaffHandler,
	patientHandler port.IPatientHandler,
) *HttpServer {
	httpServer := &HttpServer{
		config:             config,
		server:             server,
		healthCheckHandler: healthCheckHandler,
		staffHandler:       staffHandler,
		patientHandler:     patientHandler,
	}
	httpServer.initRoute()
	return httpServer
}

func (s *HttpServer) initRoute() {
	e := s.server

	e.GET("/health", s.healthCheckHandler.HealthCheck)

	// Staff routes
	e.POST("/staff/create", s.staffHandler.Create)
	e.POST("/staff/login", s.staffHandler.Login)

	// Patient routes (protected)
	auth := e.Group("/")
	auth.Use(middleware.JWTAuth(s.config.Secrets.JwtSecretKey))
	auth.POST("/patient/search", s.patientHandler.Search)

}

func (s *HttpServer) Start(address string) error {
	log.Printf("Server running on port %s", address)
	return s.server.Run(address)
}
func (s *HttpServer) Server() *gin.Engine {
	return s.server
}
