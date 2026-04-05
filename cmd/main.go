package main

import (
	"log"
	"net/http"
	"time"

	"agnos-backend/internal/adapter/his"
	"agnos-backend/internal/config"
	hospital_repository "agnos-backend/internal/repositories/hospital"
	patient_repository "agnos-backend/internal/repositories/patient"
	staff_repository "agnos-backend/internal/repositories/staff"
	patient_service "agnos-backend/internal/services/patient"
	staff_service "agnos-backend/internal/services/staff"

	handler "agnos-backend/internal/handlers/http"
	healthCheckHandler "agnos-backend/internal/handlers/http/health_check"
	patientHandler "agnos-backend/internal/handlers/http/patient"
	staffHandler "agnos-backend/internal/handlers/http/staff"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := config.Connect(&config.DBConnectConfig{
		Name:                cfg.DBConfig.Name,
		SSLMode:             cfg.DBConfig.SSLMode,
		MaxOpenConns:        cfg.DBConfig.MaxOpenConns,
		MaxIdleConns:        cfg.DBConfig.MaxIdleConns,
		ConnMaxLifetimeHour: cfg.DBConfig.ConnMaxLifetimeHour,
		Host:                cfg.DBConfig.Host,
		Port:                cfg.DBConfig.Port,
		User:                cfg.DBConfig.User,
		Password:            cfg.Secrets.DBPassword,
		TimeZone:            cfg.Server.TimeZone,
	})
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	staffRepo := staff_repository.New(db)
	patientRepo := patient_repository.New(db)
	hospitalRepo := hospital_repository.New(db)
	hisClient := his.New(&http.Client{
		Timeout: 30 * time.Second,
	})

	staffSvc := staff_service.New(staffRepo, hospitalRepo, &cfg)
	patientSvc := patient_service.New(patientRepo, hospitalRepo, hisClient)

	staffHdlr := staffHandler.New(staffSvc)
	patientHdlr := patientHandler.New(patientSvc)

	healthCheckHdlr := healthCheckHandler.New()

	server := gin.Default()

	httpServer := handler.NewHttpServer(&cfg, server, healthCheckHdlr, staffHdlr, patientHdlr)
	if err := httpServer.Start(cfg.Server.Address); err != nil {
		log.Fatal("http server error")
	}
}
