package healthcheck_handler_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	healthcheckhdl "agnos-backend/internal/handlers/http/health_check"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type HealthCheckHandlerTestSuite struct {
	suite.Suite
}

func TestHealthCheckHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HealthCheckHandlerTestSuite))
}

func (s *HealthCheckHandlerTestSuite) TestHealthCheck() {
	s.T().Run("success - HealthCheck return ok", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/health", nil)

		healthCheckServ := healthcheckhdl.New()
		healthCheckServ.HealthCheck(c)

		res := w.Result()
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		s.Assert().NoError(err)

		s.Assert().Equal(http.StatusOK, res.StatusCode)
		s.Assert().Equal("ok", string(data))
	})
}

func (s *HealthCheckHandlerTestSuite) TestReadinessCheck() {
	s.T().Run("success - Readiness return ok", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/ready", nil)

		healthCheckServ := healthcheckhdl.New()
		healthCheckServ.ReadinessCheck(c)

		res := w.Result()
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		s.Assert().NoError(err)

		s.Assert().Equal(http.StatusOK, res.StatusCode)
		s.Assert().Equal("ok", string(data))
	})
}
