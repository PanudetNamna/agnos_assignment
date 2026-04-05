package patient_handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	patient_handler "agnos-backend/internal/handlers/http/patient"
	"agnos-backend/internal/models"
	"agnos-backend/internal/port/mocks"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var (
	hospitalID = uuid.MustParse("a0000000-0000-0000-0000-000000000001")
	dbPatients = []models.Patient{
		{
			HospitalID:  hospitalID,
			NationalID:  "1234567890123",
			FirstNameEN: "Somchai",
			LastNameEN:  "Jaidee",
		},
	}
)

type PatientHandlerTestSuite struct {
	suite.Suite
	ctrl *gomock.Controller
	svc  *mocks.MockIPatientService
	r    *gin.Engine
}

func TestPatientHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(PatientHandlerTestSuite))
}

func (s *PatientHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.ctrl = gomock.NewController(s.T())
	s.svc = mocks.NewMockIPatientService(s.ctrl)
	h := patient_handler.New(s.svc)

	s.r = gin.New()
	s.r.POST("/patient/search", func(c *gin.Context) {
		c.Set("hospital_id", hospitalID.String())
		h.Search(c)
	})
}

func (s *PatientHandlerTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *PatientHandlerTestSuite) makeRequest(body interface{}) *httptest.ResponseRecorder {
	var b []byte
	switch v := body.(type) {
	case string:
		b = []byte(v)
	default:
		var err error
		b, err = json.Marshal(v)
		s.Require().NoError(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/patient/search", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.r.ServeHTTP(w, req)
	return w
}

func (s *PatientHandlerTestSuite) TestSearch() {
	type mockSetup struct {
		expectCall     bool
		filter         models.SearchRequest
		returnPatients []models.Patient
		returnErr      error
	}

	type expected struct {
		statusCode int
		message    string
	}

	tests := []struct {
		name     string
		body     interface{}
		mock     mockSetup
		expected expected
	}{
		{
			name: "success - search by national_id",
			body: models.SearchRequest{NationalID: "1234567890123"},
			mock: mockSetup{
				expectCall:     true,
				filter:         models.SearchRequest{NationalID: "1234567890123"},
				returnPatients: dbPatients,
			},
			expected: expected{statusCode: http.StatusOK, message: "success"},
		},
		{
			name: "success - empty filter",
			body: models.SearchRequest{},
			mock: mockSetup{
				expectCall:     true,
				filter:         models.SearchRequest{},
				returnPatients: dbPatients,
			},
			expected: expected{statusCode: http.StatusOK, message: "success"},
		},
		{
			name: "fail - service error",
			body: models.SearchRequest{NationalID: "1234567890123"},
			mock: mockSetup{
				expectCall: true,
				filter:     models.SearchRequest{NationalID: "1234567890123"},
				returnErr:  errors.New("internal error"),
			},
			expected: expected{statusCode: http.StatusInternalServerError, message: "internal error"},
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
					Search(hospitalID, tc.mock.filter).
					Return(tc.mock.returnPatients, tc.mock.returnErr).
					Times(1)
			}

			w := s.makeRequest(tc.body)
			s.Require().Equal(tc.expected.statusCode, w.Code)

			if tc.expected.message != "" {
				var resp map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &resp)
				s.Require().Equal(tc.expected.message, resp["message"])
			}
		})
	}
}

func (s *PatientHandlerTestSuite) TestSearch_MissingHospitalID() {
	h := patient_handler.New(s.svc)
	r := gin.New()
	r.POST("/patient/search", h.Search)

	req := httptest.NewRequest(http.MethodPost, "/patient/search", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	s.Require().Equal(http.StatusUnauthorized, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Require().Equal("unauthorized", resp["message"])
}

func (s *PatientHandlerTestSuite) TestSearch_InvalidHospitalID() {
	h := patient_handler.New(s.svc)
	r := gin.New()
	r.POST("/patient/search", func(c *gin.Context) {
		c.Set("hospital_id", "not-a-uuid")
		h.Search(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/patient/search", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	s.Require().Equal(http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Require().Equal("invalid hospital id in token", resp["message"])
}
