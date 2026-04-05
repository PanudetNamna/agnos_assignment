package patient_service_test

import (
	"errors"
	"testing"

	"agnos-backend/internal/models"
	"agnos-backend/internal/port"
	"agnos-backend/internal/port/mocks"
	patient_service "agnos-backend/internal/services/patient"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var (
	hospitalID = uuid.MustParse("a0000000-0000-0000-0000-000000000001")
	hospital   = &models.Hospital{
		ID:      hospitalID,
		Name:    "Hospital A",
		APIBase: "http://host.docker.internal:9090",
	}
	hisPatient = &models.HISPatientResponse{
		FirstNameTH: "สมชาย",
		LastNameTH:  "ใจดี",
		FirstNameEN: "Somchai",
		LastNameEN:  "Jaidee",
		DateOfBirth: "1990-01-01",
		PatientHN:   "HN001",
		NationalID:  "1234567890123",
		Gender:      "M",
	}
	dbPatients = []models.Patient{
		{
			HospitalID:  hospitalID,
			NationalID:  "1234567890123",
			FirstNameEN: "Somchai",
			LastNameEN:  "Jaidee",
		},
	}
)

type PatientServiceTestSuite struct {
	suite.Suite
	ctrl         *gomock.Controller
	patientRepo  *mocks.MockIPatientRepository
	hospitalRepo *mocks.MockIHospitalRepository
	hisClient    *mocks.MockIHisClient
	svc          port.IPatientService
}

func TestPatientServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PatientServiceTestSuite))
}

func (s *PatientServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.patientRepo = mocks.NewMockIPatientRepository(s.ctrl)
	s.hospitalRepo = mocks.NewMockIHospitalRepository(s.ctrl)
	s.hisClient = mocks.NewMockIHisClient(s.ctrl)
	s.svc = patient_service.New(s.patientRepo, s.hospitalRepo, s.hisClient)
}

func (s *PatientServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *PatientServiceTestSuite) TestSearch() {
	type mockSetup struct {
		expectHospital bool
		hospitalRes    *models.Hospital
		hospitalErr    error

		expectHIS bool
		hisID     string
		hisRes    *models.HISPatientResponse
		hisErr    error

		expectUpsert bool
		upsertErr    error

		expectSearch bool
		searchRes    []models.Patient
		searchErr    error
	}

	type expected struct {
		resultLen int
		isEmpty   bool
		err       error
		errMsg    string
	}

	tests := []struct {
		name     string
		filter   models.SearchRequest
		mock     mockSetup
		expected expected
	}{
		{
			name:   "success - search by national_id",
			filter: models.SearchRequest{NationalID: "1234567890123"},
			mock: mockSetup{
				expectHospital: true,
				hospitalRes:    hospital,
				expectHIS:      true,
				hisID:          "1234567890123",
				hisRes:         hisPatient,
				expectUpsert:   true,
				expectSearch:   true,
				searchRes:      dbPatients,
			},
			expected: expected{resultLen: 1},
		},
		{
			name:   "success - search by passport_id",
			filter: models.SearchRequest{PassportID: "AB123456"},
			mock: mockSetup{
				expectHospital: true,
				hospitalRes:    hospital,
				expectHIS:      true,
				hisID:          "AB123456",
				hisRes:         hisPatient,
				expectUpsert:   true,
				expectSearch:   true,
				searchRes:      dbPatients,
			},
			expected: expected{resultLen: 1},
		},
		{
			name:   "success - no id, skip HIS",
			filter: models.SearchRequest{FirstName: "Somchai"},
			mock: mockSetup{
				expectSearch: true,
				searchRes:    dbPatients,
			},
			expected: expected{resultLen: 1},
		},
		{
			name:   "success - empty result",
			filter: models.SearchRequest{NationalID: "9999999999999"},
			mock: mockSetup{
				expectHospital: true,
				hospitalRes:    hospital,
				expectHIS:      true,
				hisID:          "9999999999999",
				hisRes:         hisPatient,
				expectUpsert:   true,
				expectSearch:   true,
				searchRes:      []models.Patient{},
			},
			expected: expected{isEmpty: true},
		},
		{
			name:   "success - upsert error continues search",
			filter: models.SearchRequest{NationalID: "1234567890123"},
			mock: mockSetup{
				expectHospital: true,
				hospitalRes:    hospital,
				expectHIS:      true,
				hisID:          "1234567890123",
				hisRes:         hisPatient,
				expectUpsert:   true,
				upsertErr:      errors.New("upsert failed"),
				expectSearch:   true,
				searchRes:      dbPatients,
			},
			expected: expected{resultLen: 1},
		},
		{
			name:   "fail - hospital not found",
			filter: models.SearchRequest{NationalID: "1234567890123"},
			mock: mockSetup{
				expectHospital: true,
				hospitalErr:    errors.New("hospital not found"),
			},
			expected: expected{err: errors.New("find hospital"), errMsg: "find hospital"},
		},
		{
			name:   "fail - HIS client error",
			filter: models.SearchRequest{NationalID: "1234567890123"},
			mock: mockSetup{
				expectHospital: true,
				hospitalRes:    hospital,
				expectHIS:      true,
				hisID:          "1234567890123",
				hisErr:         errors.New("connection refused"),
			},
			expected: expected{err: errors.New("fetch patient"), errMsg: "fetch patient"},
		},
		{
			name:   "fail - db search error",
			filter: models.SearchRequest{NationalID: "1234567890123"},
			mock: mockSetup{
				expectHospital: true,
				hospitalRes:    hospital,
				expectHIS:      true,
				hisID:          "1234567890123",
				hisRes:         hisPatient,
				expectUpsert:   true,
				expectSearch:   true,
				searchErr:      errors.New("db error"),
			},
			expected: expected{err: errors.New("db error")},
		},
	}

	for _, tc := range tests {
		s.T().Run(tc.name, func(t *testing.T) {
			if tc.mock.expectHospital {
				s.hospitalRepo.EXPECT().
					FindByID(hospitalID).
					Return(tc.mock.hospitalRes, tc.mock.hospitalErr).
					Times(1)
			}

			if tc.mock.expectHIS {
				s.hisClient.EXPECT().
					FetchPatient(hospital.APIBase, tc.mock.hisID).
					Return(tc.mock.hisRes, tc.mock.hisErr).
					Times(1)
			}

			if tc.mock.expectUpsert {
				s.patientRepo.EXPECT().
					Upsert(gomock.Any()).
					Return(tc.mock.upsertErr).
					Times(1)
			}

			if tc.mock.expectSearch {
				s.patientRepo.EXPECT().
					Search(hospitalID, tc.filter).
					Return(tc.mock.searchRes, tc.mock.searchErr).
					Times(1)
			}

			result, err := s.svc.Search(hospitalID, tc.filter)

			if tc.expected.err != nil {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expected.errMsg)
				s.Require().Nil(result)
				return
			}

			s.Require().NoError(err)
			if tc.expected.isEmpty {
				s.Require().Empty(result)
			} else {
				s.Require().Len(result, tc.expected.resultLen)
			}
		})
	}
}
