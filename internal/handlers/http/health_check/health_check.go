package healthcheck_handler

import "github.com/gin-gonic/gin"

type IHealthCheckHandler interface {
	HealthCheck(c *gin.Context)
	ReadinessCheck(c *gin.Context)
}

type HealthCheckHandler struct{}

func New() *HealthCheckHandler {
	return &HealthCheckHandler{}
}

func (h *HealthCheckHandler) HealthCheck(c *gin.Context) {
	c.String(200, "ok")
}

func (h *HealthCheckHandler) ReadinessCheck(c *gin.Context) {
	c.String(200, "ok")
}
