package staff_handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	staff_handler "agnos-backend/internal/handlers/http/staff"
	"agnos-backend/internal/models"
	"agnos-backend/internal/port/mocks"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var (
	staffID    = uuid.MustParse("b0000000-0000-0000-0000-000000000001")
	hospitalID = uuid.MustParse("a0000000-0000-0000-0000-000000000001")
	mockStaff  = &models.Staff{
		ID:         staffID,
		Username:   "john",
		HospitalID: hospitalID,
	}
)

type StaffHandlerTestSuite struct {
	suite.Suite
	ctrl *gomock.Controller
	svc  *mocks.MockIStaffService
	r    *gin.Engine
}

func TestStaffHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(StaffHandlerTestSuite))
}

func (s *StaffHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.ctrl = gomock.NewController(s.T())
	s.svc = mocks.NewMockIStaffService(s.ctrl)
	h := staff_handler.New(s.svc)
	s.r = gin.New()
	s.r.POST("/staff/create", h.Create)
	s.r.POST("/staff/login", h.Login)
}

func (s *StaffHandlerTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *StaffHandlerTestSuite) makeRequest(method, path string, body interface{}) *httptest.ResponseRecorder {
	var b []byte
	switch v := body.(type) {
	case string:
		b = []byte(v)
	default:
		var err error
		b, err = json.Marshal(v)
		s.Require().NoError(err)
	}
	req := httptest.NewRequest(method, path, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.r.ServeHTTP(w, req)
	return w
}

func (s *StaffHandlerTestSuite) TestCreate() {
	type mockSetup struct {
		expectCall  bool
		username    string
		password    string
		hospital    string
		returnStaff *models.Staff
		returnErr   error
	}

	type expected struct {
		statusCode int
		message    string
		staffID    string
		hospitalID string
	}

	tests := []struct {
		name     string
		body     interface{}
		mock     mockSetup
		expected expected
	}{
		{
			name: "success",
			body: models.CreateStaffRequest{Username: "john", Password: "password123", Hospital: "Hospital A"},
			mock: mockSetup{
				expectCall:  true,
				username:    "john",
				password:    "password123",
				hospital:    "Hospital A",
				returnStaff: mockStaff,
			},
			expected: expected{
				statusCode: http.StatusOK,
				message:    "success",
				staffID:    staffID.String(),
				hospitalID: hospitalID.String(),
			},
		},
		{
			name: "fail - hospital not found",
			body: models.CreateStaffRequest{Username: "john", Password: "password123", Hospital: "Unknown"},
			mock: mockSetup{
				expectCall: true,
				username:   "john",
				password:   "password123",
				hospital:   "Unknown",
				returnErr:  errors.New("hospital not found"),
			},
			expected: expected{statusCode: http.StatusInternalServerError, message: "hospital not found"},
		},
		{
			name:     "fail - missing username",
			body:     models.CreateStaffRequest{Password: "password123", Hospital: "Hospital A"},
			mock:     mockSetup{expectCall: false},
			expected: expected{statusCode: http.StatusBadRequest},
		},
		{
			name:     "fail - missing password",
			body:     models.CreateStaffRequest{Username: "john", Hospital: "Hospital A"},
			mock:     mockSetup{expectCall: false},
			expected: expected{statusCode: http.StatusBadRequest},
		},
		{
			name:     "fail - missing hospital",
			body:     models.CreateStaffRequest{Username: "john", Password: "password123"},
			mock:     mockSetup{expectCall: false},
			expected: expected{statusCode: http.StatusBadRequest},
		},
		{
			name:     "fail - invalid body",
			body:     "not-json",
			mock:     mockSetup{expectCall: false},
			expected: expected{statusCode: http.StatusBadRequest},
		},
	}

	for _, tc := range tests {
		s.T().Run(tc.name, func(t *testing.T) {
			if tc.mock.expectCall {
				s.svc.EXPECT().
					Create(tc.mock.username, tc.mock.password, tc.mock.hospital).
					Return(tc.mock.returnStaff, tc.mock.returnErr).
					Times(1)
			}

			w := s.makeRequest(http.MethodPost, "/staff/create", tc.body)
			s.Require().Equal(tc.expected.statusCode, w.Code)

			if tc.expected.message != "" {
				var resp map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &resp)
				s.Require().Equal(tc.expected.message, resp["message"])
			}

			if tc.expected.staffID != "" {
				var resp map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &resp)
				data := resp["data"].(map[string]interface{})
				s.Require().Equal(tc.expected.staffID, data["staff_id"])
				s.Require().Equal(tc.expected.hospitalID, data["hospital_id"])
			}
		})
	}
}

func (s *StaffHandlerTestSuite) TestLogin() {
	type mockSetup struct {
		expectCall  bool
		username    string
		password    string
		hospital    string
		returnToken string
		returnErr   error
	}

	type expected struct {
		statusCode int
		message    string
		token      string
	}

	tests := []struct {
		name     string
		body     interface{}
		mock     mockSetup
		expected expected
	}{
		{
			name: "success",
			body: models.LoginRequest{Username: "john", Password: "password123", Hospital: "Hospital A"},
			mock: mockSetup{
				expectCall:  true,
				username:    "john",
				password:    "password123",
				hospital:    "Hospital A",
				returnToken: "jwt-token",
			},
			expected: expected{statusCode: http.StatusOK, message: "success", token: "jwt-token"},
		},
		{
			name: "fail - invalid credentials",
			body: models.LoginRequest{Username: "john", Password: "wrongpassword", Hospital: "Hospital A"},
			mock: mockSetup{
				expectCall: true,
				username:   "john",
				password:   "wrongpassword",
				hospital:   "Hospital A",
				returnErr:  errors.New("invalid credentials"),
			},
			expected: expected{statusCode: http.StatusUnauthorized, message: "invalid credentials"},
		},
		{
			name: "fail - hospital not found",
			body: models.LoginRequest{Username: "john", Password: "password123", Hospital: "Unknown"},
			mock: mockSetup{
				expectCall: true,
				username:   "john",
				password:   "password123",
				hospital:   "Unknown",
				returnErr:  errors.New("hospital not found"),
			},
			expected: expected{statusCode: http.StatusUnauthorized, message: "hospital not found"},
		},
		{
			name:     "fail - missing fields",
			body:     models.LoginRequest{Username: "john"},
			mock:     mockSetup{expectCall: false},
			expected: expected{statusCode: http.StatusBadRequest},
		},
		{
			name:     "fail - invalid body",
			body:     "not-json",
			mock:     mockSetup{expectCall: false},
			expected: expected{statusCode: http.StatusBadRequest},
		},
	}

	for _, tc := range tests {
		s.T().Run(tc.name, func(t *testing.T) {
			if tc.mock.expectCall {
				s.svc.EXPECT().
					Login(tc.mock.username, tc.mock.password, tc.mock.hospital).
					Return(tc.mock.returnToken, tc.mock.returnErr).
					Times(1)
			}

			w := s.makeRequest(http.MethodPost, "/staff/login", tc.body)
			s.Require().Equal(tc.expected.statusCode, w.Code)

			if tc.expected.message != "" {
				var resp map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &resp)
				s.Require().Equal(tc.expected.message, resp["message"])
			}

			if tc.expected.token != "" {
				var resp map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &resp)
				data := resp["data"].(map[string]interface{})
				s.Require().Equal(tc.expected.token, data["token"])
			}
		})
	}
}
