package handler_test

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SawitProRecruitment/UserService/config"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/repository"
	apperror "github.com/SawitProRecruitment/UserService/utils/app_error"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestPostEstate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	server := &handler.Server{Repository: mockRepo, Validator: validator.New(), Config: &config.Config{}}
	e := echo.New()

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		setup          func(*repository.MockRepositoryInterface)
	}{
		{
			name: "Success",
			requestBody: `{
				"width": 50,
				"length": 100
			}`,
			expectedStatus: http.StatusCreated,
			setup: func(mockRepo *repository.MockRepositoryInterface) {
				mockRepo.EXPECT().
					CreateEstate(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:           "Invalid JSON (Bind Error)",
			requestBody:    `{invalid_json}`,
			expectedStatus: http.StatusBadRequest,
			setup:          func(mockRepo *repository.MockRepositoryInterface) {},
		},
		{
			name: "Payload Validation Error",
			requestBody: `{
				"width": 0,
				"length": 0
			}`,
			expectedStatus: http.StatusBadRequest,
			setup:          func(mri *repository.MockRepositoryInterface) {},
		},
		{
			name: "Repository Error",
			requestBody: `{
				"width": 50,
				"length": 100
			}`,
			expectedStatus: http.StatusInternalServerError,
			setup: func(mockRepo *repository.MockRepositoryInterface) {
				mockRepo.EXPECT().
					CreateEstate(gomock.Any(), gomock.Any()).
					Return(apperror.WrapWithCode(errors.New("database error"), http.StatusInternalServerError)).
					Times(1)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodPost, "/estate", bytes.NewReader([]byte(tc.requestBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			tc.setup(mockRepo)

			err := server.PostEstate(c)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}

func TestGetEstateIdDronePlan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	server := &handler.Server{Repository: mockRepo, Validator: validator.New(), Config: &config.Config{}}
	e := echo.New()

	validID := openapi_types.UUID(uuid.New())
	validParams := generated.GetEstateIdDronePlanParams{MaxDistance: ptr(1000)} // MaxDistance = 1000 meters

	mockTree := []repository.Tree{
		{Id: openapi_types.UUID(uuid.New()), X: 2, Y: 1, Height: 5},
		{Id: openapi_types.UUID(uuid.New()), X: 3, Y: 1, Height: 3},
		{Id: openapi_types.UUID(uuid.New()), X: 4, Y: 1, Height: 4},
	}

	mockEstate := &repository.Estate{
		Id:     validID,
		Width:  1,
		Length: 5,
		Trees:  mockTree,
	}
	// mockDistance := 500

	tests := []struct {
		name           string
		id             openapi_types.UUID
		params         generated.GetEstateIdDronePlanParams
		mockReturnData *repository.Estate
		expectedStatus int
		setup          func(*repository.MockRepositoryInterface)
	}{
		{
			name:           "Success",
			id:             validID,
			params:         validParams,
			mockReturnData: mockEstate,
			expectedStatus: http.StatusOK,
			setup: func(mockRepo *repository.MockRepositoryInterface) {
				mockRepo.EXPECT().
					GetEstateWithAllDetails(gomock.Any(), validID).
					Return(mockEstate, nil).
					Times(1)
			},
		},
		// {
		// 	name:           "Payload Validation Error",
		// 	id:             validID,
		// 	params:         invalidParams, // Invalid max_distance
		// 	expectedStatus: http.StatusBadRequest,
		// 	setup:          func(mockRepo *repository.MockRepositoryInterface) {},
		// },
		// {
		// 	name:           "Estate Not Found",
		// 	id:             validID,
		// 	params:         validParams,
		// 	mockReturnErr:  apperror.WrapWithCode(errors.New("estate not found"), http.StatusNotFound),
		// 	mockReturnData: nil,
		// 	expectedStatus: http.StatusNotFound,
		// 	setup: func(mockRepo *repository.MockRepositoryInterface) {
		// 		mockRepo.EXPECT().
		// 			GetEstateByIdWithTrees(gomock.Any(), validID).
		// 			Return(nil, apperror.WrapWithCode(errors.New("estate not found"), http.StatusNotFound)).
		// 			Times(1)
		// 	},
		// },
		{
			name:           "Repository Error",
			id:             validID,
			params:         validParams,
			mockReturnData: nil,
			expectedStatus: http.StatusInternalServerError,
			setup: func(mockRepo *repository.MockRepositoryInterface) {
				mockRepo.EXPECT().
					GetEstateWithAllDetails(gomock.Any(), validID).
					Return(nil, apperror.WrapWithCode(errors.New("database error"), http.StatusInternalServerError)).
					Times(1)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/estate/%s/drone-plan", tc.id), nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			tc.setup(mockRepo)

			err := server.GetEstateIdDronePlan(c, tc.id, tc.params)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}

func TestGetEstateIdStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	server := &handler.Server{
		Repository: mockRepo,
	}
	e := echo.New()

	// Sample UUID
	testID := uuid.New()

	tests := []struct {
		name           string
		estateID       uuid.UUID
		expectedStatus int
		expectedBody   generated.EstateStats
		setup          func(*repository.MockRepositoryInterface)
	}{
		{
			name:           "Success",
			estateID:       testID,
			expectedStatus: http.StatusOK,
			expectedBody: generated.EstateStats{
				Count:  ptr(int64(10)),
				Max:    ptr(100),
				Median: ptr(50),
				Min:    ptr(5),
			},
			setup: func(mockRepo *repository.MockRepositoryInterface) {
				mockRepo.EXPECT().
					GetEstateWithAllDetails(gomock.Any(), testID, repository.RELATION_TREES).
					Return(&repository.Estate{
						Stats: &repository.EstateStats{
							TreeCount:    10,
							MaxHeight:    10,
							MedianHeight: 50,
							MinHeight:    5,
						},
					}, nil).
					Times(1)
			},
		},
		{
			name:           "Estate Not Found",
			estateID:       testID,
			expectedStatus: http.StatusNotFound,
			setup: func(mockRepo *repository.MockRepositoryInterface) {
				mockRepo.EXPECT().
					GetEstateWithAllDetails(gomock.Any(), testID, repository.RELATION_TREES).
					Return(&repository.Estate{}, apperror.WrapWithCode(errors.New("estate not found"), http.StatusNotFound)).
					Times(1)
			},
		},
		{
			name:           "Database Error",
			estateID:       testID,
			expectedStatus: http.StatusInternalServerError,
			setup: func(mockRepo *repository.MockRepositoryInterface) {
				mockRepo.EXPECT().
					GetEstateWithAllDetails(gomock.Any(), testID, repository.RELATION_TREES).
					Return(&repository.Estate{}, apperror.WrapWithCode(errors.New("database error"), http.StatusInternalServerError)).
					Times(1)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/estate/"+tc.estateID.String()+"/stats", bytes.NewReader([]byte{}))
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tc.estateID.String())

			tc.setup(mockRepo)

			err := server.GetEstateIdStats(c, testID)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}

func TestPostEstateIdTree(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	server := &handler.Server{
		Repository: mockRepo,
		Validator:  validator.New(),
	}
	e := echo.New()

	// Sample UUID
	testEstateID := uuid.New()

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		setup          func(*repository.MockRepositoryInterface)
	}{
		{
			name: "Success",
			requestBody: `{
				"x": 10,
				"y": 20,
				"height": 15
			}`,
			expectedStatus: http.StatusCreated,
			setup: func(mockRepo *repository.MockRepositoryInterface) {
				mockRepo.EXPECT().
					CreateTree(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)

				mockRepo.EXPECT().
					GetEstateWithAllDetails(gomock.Any(), gomock.Any()).
					Return(&repository.Estate{
						Id:     uuid.New(),
						Width:  1,
						Length: 1,
						Trees: []repository.Tree{
							{X: 1, Y: 1, Height: 10},
						},
						Stats: &repository.EstateStats{
							TreeCount:    1,
							MaxHeight:    10,
							MedianHeight: 10,
							MinHeight:    5,
						},
					}, nil).
					Times(1)

				mockRepo.EXPECT().
					GetCalculatedEstateStats(gomock.Any(), gomock.Any()).Return(
					&repository.EstateStats{
						Id:           uuid.New(),
						TreeCount:    1,
						MaxHeight:    10,
						MedianHeight: 10,
						MinHeight:    5,
					}, nil,
				)

				mockRepo.EXPECT().UpsertEstateStats(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:           "Invalid JSON (Bind Error)",
			requestBody:    `{invalid_json}`,
			expectedStatus: http.StatusBadRequest,
			setup:          func(mockRepo *repository.MockRepositoryInterface) {},
		},
		{
			name: "Payload Validation Error",
			requestBody: `{
				"x": -10,
				"y": -20,
				"height": 0
			}`,
			expectedStatus: http.StatusBadRequest,
			setup:          func(mockRepo *repository.MockRepositoryInterface) {},
		},
		{
			name: "Repository Create Tree Error",
			requestBody: `{
				"x": 10,
				"y": 20,
				"height": 15
			}`,
			expectedStatus: http.StatusInternalServerError,
			setup: func(mockRepo *repository.MockRepositoryInterface) {
				mockRepo.EXPECT().
					CreateTree(gomock.Any(), gomock.Any()).
					Return(apperror.WrapWithCode(errors.New("database error"), http.StatusInternalServerError)).
					Times(1)
			},
		},
		{
			name: "Repository GetCalculatedEstateStats Error",
			requestBody: `{
				"x": 10,
				"y": 20,
				"height": 15
			}`,
			expectedStatus: http.StatusInternalServerError,
			setup: func(mockRepo *repository.MockRepositoryInterface) {
				mockRepo.EXPECT().
					CreateTree(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)

				mockRepo.EXPECT().
					GetEstateWithAllDetails(gomock.Any(), gomock.Any()).
					Return(&repository.Estate{
						Id:     uuid.New(),
						Width:  1,
						Length: 1,
						Trees: []repository.Tree{
							{X: 1, Y: 1, Height: 10},
						},
						Stats: &repository.EstateStats{
							TreeCount:    1,
							MaxHeight:    10,
							MedianHeight: 10,
							MinHeight:    5,
						},
					}, nil).
					Times(1)

				mockRepo.EXPECT().
					GetCalculatedEstateStats(gomock.Any(), gomock.Any()).Return(
					nil, errors.New("database error"),
				)
			},
		},
		{
			name: "Repository GetCalculatedEstateStats Error",
			requestBody: `{
				"x": 10,
				"y": 20,
				"height": 15
			}`,
			expectedStatus: http.StatusInternalServerError,
			setup: func(mockRepo *repository.MockRepositoryInterface) {
				mockRepo.EXPECT().
					CreateTree(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)

				mockRepo.EXPECT().
					GetEstateWithAllDetails(gomock.Any(), gomock.Any()).
					Return(&repository.Estate{
						Id:     uuid.New(),
						Width:  1,
						Length: 1,
						Trees: []repository.Tree{
							{X: 1, Y: 1, Height: 10},
						},
						Stats: &repository.EstateStats{
							TreeCount:    1,
							MaxHeight:    10,
							MedianHeight: 10,
							MinHeight:    5,
						},
					}, nil).
					Times(1)

				mockRepo.EXPECT().
					GetCalculatedEstateStats(gomock.Any(), gomock.Any()).Return(
					nil, errors.New("database error"),
				)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/estate/%s/tree", testEstateID.String()), bytes.NewReader([]byte(tc.requestBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(testEstateID.String())

			tc.setup(mockRepo)

			err := server.PostEstateIdTree(c, testEstateID)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}

// Helper function to create pointers
func ptr[T any](v T) *T {
	return &v
}
