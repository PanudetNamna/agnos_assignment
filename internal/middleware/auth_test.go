package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"agnos-backend/internal/middleware"
	"agnos-backend/internal/utility"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const testSecret = "test-secret-key"

func makeToken(t *testing.T, staffID, hospitalID uuid.UUID, secret string, expiry time.Duration) string {
	claims := jwt.MapClaims{
		"staff_id":    staffID.String(),
		"hospital_id": hospitalID.String(),
		"exp":         time.Now().Add(expiry).Unix(),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	assert.NoError(t, err)
	return token
}

func setupRouter(secret string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/protected", middleware.JWTAuth(secret), func(c *gin.Context) {
		utility.Success(c, gin.H{
			"staff_id":    c.GetString("staff_id"),
			"hospital_id": c.GetString("hospital_id"),
		})
	})
	return r
}

func doRequest(r *gin.Engine, authHeader string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

var (
	staffID    = uuid.MustParse("b0000000-0000-0000-0000-000000000001")
	hospitalID = uuid.MustParse("a0000000-0000-0000-0000-000000000001")
)

func TestJWTAuth_ValidToken(t *testing.T) {
	r := setupRouter(testSecret)
	token := makeToken(t, staffID, hospitalID, testSecret, time.Hour)

	w := doRequest(r, "Bearer "+token)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJWTAuth_ClaimsSetInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/protected", middleware.JWTAuth(testSecret), func(c *gin.Context) {
		assert.Equal(t, staffID.String(), c.GetString("staff_id"))
		assert.Equal(t, hospitalID.String(), c.GetString("hospital_id"))
		c.Status(http.StatusOK)
	})

	token := makeToken(t, staffID, hospitalID, testSecret, time.Hour)
	w := doRequest(r, "Bearer "+token)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJWTAuth_MissingHeader(t *testing.T) {
	r := setupRouter(testSecret)

	w := doRequest(r, "")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTAuth_MissingBearerPrefix(t *testing.T) {
	r := setupRouter(testSecret)
	token := makeToken(t, staffID, hospitalID, testSecret, time.Hour)

	w := doRequest(r, token)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTAuth_ExpiredToken(t *testing.T) {
	r := setupRouter(testSecret)
	token := makeToken(t, staffID, hospitalID, testSecret, -time.Hour) // หมดอายุแล้ว

	w := doRequest(r, "Bearer "+token)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTAuth_WrongSecret(t *testing.T) {
	r := setupRouter(testSecret)
	token := makeToken(t, staffID, hospitalID, "wrong-secret", time.Hour)

	w := doRequest(r, "Bearer "+token)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTAuth_MalformedToken(t *testing.T) {
	r := setupRouter(testSecret)

	w := doRequest(r, "Bearer not.a.valid.token")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTAuth_EmptyToken(t *testing.T) {
	r := setupRouter(testSecret)

	w := doRequest(r, "Bearer ")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
