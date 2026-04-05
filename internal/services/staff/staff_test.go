package staff_service_test

import (
	"errors"
	"testing"

	"agnos-backend/internal/config"
	"agnos-backend/internal/models"
	"agnos-backend/internal/port"
	"agnos-backend/internal/port/mocks"
	staff_service "agnos-backend/internal/services/staff"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

var (
	hospitalID = uuid.MustParse("a0000000-0000-0000-0000-000000000001")
	staffID    = uuid.MustParse("b0000000-0000-0000-0000-000000000001")
	hospital   = &models.Hospital{ID: hospitalID, Name: "Hospital A"}
	testCfg    = &config.AppConfig{
		Secrets: config.Secrets{JwtSecretKey: "test-secret-key"},
	}
)

type StaffServiceTestSuite struct {
	suite.Suite
	ctrl         *gomock.Controller
	staffRepo    *mocks.MockIStaffRepository
	hospitalRepo *mocks.MockIHospitalRepository
	svc          port.IStaffService
}

func TestStaffServiceTestSuite(t *testing.T) {
	suite.Run(t, new(StaffServiceTestSuite))
}

func (s *StaffServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.staffRepo = mocks.NewMockIStaffRepository(s.ctrl)
	s.hospitalRepo = mocks.NewMockIHospitalRepository(s.ctrl)
	s.svc = staff_service.New(s.staffRepo, s.hospitalRepo, testCfg)
}

func (s *StaffServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *StaffServiceTestSuite) TestCreate() {
	type mockSetup struct {
		expectHospital bool
		hospitalName   string
		hospitalRes    *models.Hospital
		hospitalErr    error

		expectCreate bool
		createErr    error
	}

	type expected struct {
		username   string
		hospitalID uuid.UUID
		errMsg     string
	}

	tests := []struct {
		name     string
		username string
		password string
		hospital string
		mock     mockSetup
		expected expected
		wantErr  bool
	}{
		{
			name:     "success",
			username: "john",
			password: "password123",
			hospital: "Hospital A",
			mock: mockSetup{
				expectHospital: true,
				hospitalName:   "Hospital A",
				hospitalRes:    hospital,
				expectCreate:   true,
			},
			expected: expected{username: "john", hospitalID: hospitalID},
		},
		{
			name:     "fail - hospital not found",
			username: "john",
			password: "password123",
			hospital: "Unknown",
			mock: mockSetup{
				expectHospital: true,
				hospitalName:   "Unknown",
				hospitalErr:    errors.New("not found"),
			},
			wantErr:  true,
			expected: expected{errMsg: "hospital not found"},
		},
		{
			name:     "fail - db insert error",
			username: "john",
			password: "password123",
			hospital: "Hospital A",
			mock: mockSetup{
				expectHospital: true,
				hospitalName:   "Hospital A",
				hospitalRes:    hospital,
				expectCreate:   true,
				createErr:      errors.New("duplicate username"),
			},
			wantErr:  true,
			expected: expected{errMsg: "duplicate username"},
		},
	}

	for _, tc := range tests {
		s.T().Run(tc.name, func(t *testing.T) {
			if tc.mock.expectHospital {
				s.hospitalRepo.EXPECT().
					FindByName(tc.mock.hospitalName).
					Return(tc.mock.hospitalRes, tc.mock.hospitalErr).
					Times(1)
			}

			if tc.mock.expectCreate {
				s.staffRepo.EXPECT().
					Create(gomock.Any()).
					Return(tc.mock.createErr).
					Times(1)
			}

			result, err := s.svc.Create(tc.username, tc.password, tc.hospital)

			if tc.wantErr {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.expected.errMsg)
				s.Require().Nil(result)
				return
			}

			s.Require().NoError(err)
			s.Require().NotNil(result)
			s.Require().Equal(tc.expected.username, result.Username)
			s.Require().Equal(tc.expected.hospitalID, result.HospitalID)
			s.Require().NotEqual(tc.password, result.Password)
		})
	}
}

func (s *StaffServiceTestSuite) TestLogin() {
	hashedCorrect, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	hashedOther, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	type mockSetup struct {
		expectHospital bool
		hospitalName   string
		hospitalRes    *models.Hospital
		hospitalErr    error

		expectStaff bool
		staffRes    *models.Staff
		staffErr    error
	}

	type expected struct {
		hasToken bool
		errMsg   string
	}

	tests := []struct {
		name     string
		username string
		password string
		hospital string
		mock     mockSetup
		expected expected
		wantErr  bool
	}{
		{
			name:     "success",
			username: "john",
			password: "password123",
			hospital: "Hospital A",
			mock: mockSetup{
				expectHospital: true,
				hospitalName:   "Hospital A",
				hospitalRes:    hospital,
				expectStaff:    true,
				staffRes: &models.Staff{
					ID:         staffID,
					Username:   "john",
					Password:   string(hashedCorrect),
					HospitalID: hospitalID,
				},
			},
			expected: expected{hasToken: true},
		},
		{
			name:     "fail - hospital not found",
			username: "john",
			password: "password123",
			hospital: "Unknown",
			mock: mockSetup{
				expectHospital: true,
				hospitalName:   "Unknown",
				hospitalErr:    errors.New("not found"),
			},
			wantErr:  true,
			expected: expected{errMsg: "hospital not found"},
		},
		{
			name:     "fail - staff not found",
			username: "john",
			password: "password123",
			hospital: "Hospital A",
			mock: mockSetup{
				expectHospital: true,
				hospitalName:   "Hospital A",
				hospitalRes:    hospital,
				expectStaff:    true,
				staffErr:       errors.New("not found"),
			},
			wantErr:  true,
			expected: expected{errMsg: "invalid credentials"},
		},
		{
			name:     "fail - wrong password",
			username: "john",
			password: "wrongpassword",
			hospital: "Hospital A",
			mock: mockSetup{
				expectHospital: true,
				hospitalName:   "Hospital A",
				hospitalRes:    hospital,
				expectStaff:    true,
				staffRes: &models.Staff{
					ID:         staffID,
					Username:   "john",
					Password:   string(hashedOther),
					HospitalID: hospitalID,
				},
			},
			wantErr:  true,
			expected: expected{errMsg: "invalid credentials"},
		},
	}

	for _, tc := range tests {
		s.T().Run(tc.name, func(t *testing.T) {
			if tc.mock.expectHospital {
				s.hospitalRepo.EXPECT().
					FindByName(tc.mock.hospitalName).
					Return(tc.mock.hospitalRes, tc.mock.hospitalErr).
					Times(1)
			}

			if tc.mock.expectStaff {
				s.staffRepo.EXPECT().
					FindByUsernameAndHospital(tc.username, hospitalID).
					Return(tc.mock.staffRes, tc.mock.staffErr).
					Times(1)
			}

			token, err := s.svc.Login(tc.username, tc.password, tc.hospital)

			if tc.wantErr {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.expected.errMsg)
				s.Require().Empty(token)
				return
			}

			s.Require().NoError(err)
			s.Require().NotEmpty(token)
		})
	}
}
