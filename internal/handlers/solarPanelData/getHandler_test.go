package solarPanelData

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/loukaspe/solar-panel-data-crud/internal/core/domain"
	mock_services "github.com/loukaspe/solar-panel-data-crud/mocks/mock_internal/core/services"
	mock_helper "github.com/loukaspe/solar-panel-data-crud/mocks/mock_pkg/helper"
	apierrors "github.com/loukaspe/solar-panel-data-crud/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetSolarPanelDataHandler_GetSolarPanelDataController(t *testing.T) {
	logger := logrus.New()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockService := mock_services.NewMockSolarPanelDataServiceInterface(mockCtrl)
	mockEventExtractor := mock_helper.NewMockSolarPanelDataEventExtractorInterface(mockCtrl)

	tests := []struct {
		name                            string
		requestedUuid                   string
		shouldMockServiceRun            bool
		mockServiceResponseData         *domain.SolarPanelData
		mockServiceResponseError        error
		shouldMockEventExtractorRun     bool
		mockEventExtractorResponseData  [][]string
		mockEventExtractorResponseError error
		expected                        string
		expectedStatusCode              int
	}{
		{
			name:                 "valid",
			requestedUuid:        "38d503e5-dc1c-4549-8172-09d9c29070f7",
			shouldMockServiceRun: true,
			mockServiceResponseData: &domain.SolarPanelData{
				Solar: map[string][][]string{
					"38d503e5-dc1c-4549-8172-09d9c29070f7": [][]string{
						{"20211231T221500Z", "0.0"},
					},
				},
				Wind: nil,
			},
			mockServiceResponseError:    nil,
			shouldMockEventExtractorRun: true,
			mockEventExtractorResponseData: [][]string{
				{"Events"},
				{"0.0"},
			},
			mockEventExtractorResponseError: nil,
			expected: `Events
0.0
`,
			expectedStatusCode: 200,
		},
		{
			name:                        "missing id",
			requestedUuid:               "",
			shouldMockServiceRun:        false,
			shouldMockEventExtractorRun: false,
			expected: `missing solarPanelData id
`,
			expectedStatusCode: 400,
		},
		{
			name:                 "invalid service data not found error",
			requestedUuid:        "aaaaaa",
			shouldMockServiceRun: true,
			mockServiceResponseData: &domain.SolarPanelData{
				Solar: map[string][][]string{
					"38d503e5-dc1c-4549-8172-09d9c29070f7": [][]string{
						{"20211231T221500Z", "0.0"},
					},
				},
				Wind: nil,
			},
			mockServiceResponseError: &apierrors.DataNotFoundErrorWrapper{
				ReturnedStatusCode: 204,
				OriginalError:      errors.New("uuid aaaaaa not found"),
			},
			shouldMockEventExtractorRun: false,
			expected:                    ``,
			expectedStatusCode:          204,
		},
		{
			name:                 "invalid service random error",
			requestedUuid:        "aaaaaa",
			shouldMockServiceRun: true,
			mockServiceResponseData: &domain.SolarPanelData{
				Solar: map[string][][]string{
					"38d503e5-dc1c-4549-8172-09d9c29070f7": [][]string{
						{"20211231T221500Z", "0.0"},
					},
				},
				Wind: nil,
			},
			mockServiceResponseError:    errors.New("random error"),
			shouldMockEventExtractorRun: false,
			expected:                    ``,
			expectedStatusCode:          500,
		},
		{
			name:                 "invalid malformed event data error",
			requestedUuid:        "38d503e5-dc1c-4549-8172-09d9c29070f7",
			shouldMockServiceRun: true,
			mockServiceResponseData: &domain.SolarPanelData{
				Solar: map[string][][]string{
					"38d503e5-dc1c-4549-8172-09d9c29070f7": [][]string{
						{"20211231T221500Z"},
					},
				},
				Wind: nil,
			},
			mockServiceResponseError:       nil,
			shouldMockEventExtractorRun:    true,
			mockEventExtractorResponseData: [][]string{},
			mockEventExtractorResponseError: apierrors.MalformedEventDataError{
				ReturnedStatusCode:   http.StatusInternalServerError,
				MalformedParameterId: "38d503e5-dc1c-4549-8172-09d9c29070f7",
				OriginalError:        errors.New("parameterId 38d503e5-dc1c-4549-8172-09d9c29070f7 contains events with no values"),
			},
			expected:           "\"malformed solar panel data, check parameter 38d503e5-dc1c-4549-8172-09d9c29070f7\"\n",
			expectedStatusCode: 500,
		},
		{
			name:                 "invalid extractor random error",
			requestedUuid:        "38d503e5-dc1c-4549-8172-09d9c29070f7",
			shouldMockServiceRun: true,
			mockServiceResponseData: &domain.SolarPanelData{
				Solar: map[string][][]string{
					"38d503e5-dc1c-4549-8172-09d9c29070f7": [][]string{
						{"20211231T221500Z"},
					},
				},
				Wind: nil,
			},
			mockServiceResponseError:        nil,
			shouldMockEventExtractorRun:     true,
			mockEventExtractorResponseData:  [][]string{},
			mockEventExtractorResponseError: errors.New("random error"),
			expected:                        ``,
			expectedStatusCode:              500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRequest := httptest.NewRequest(
				"GET",
				"/solarPanelData",
				nil,
			)
			vars := map[string]string{
				"id": tt.requestedUuid,
			}
			mockRequest = mux.SetURLVars(mockRequest, vars)
			mockResponseRecorder := httptest.NewRecorder()

			if tt.shouldMockServiceRun {
				mockService.EXPECT().
					GetSolarPanelData(tt.requestedUuid).
					Return(tt.mockServiceResponseData, tt.mockServiceResponseError)
			}

			if tt.shouldMockEventExtractorRun {
				mockEventExtractor.EXPECT().
					ExtractEventsPerParameterIdToCsvForm(tt.mockServiceResponseData).
					Return(tt.mockEventExtractorResponseData, tt.mockEventExtractorResponseError)
			}

			handler := &GetSolarPanelDataHandler{
				SolarPanelDataService:        mockService,
				SolarPanelDataEventExtractor: mockEventExtractor,
				logger:                       logger,
			}
			sut := handler.GetSolarPanelDataController

			sut(mockResponseRecorder, mockRequest)

			mockResponse := mockResponseRecorder.Result()
			actual, err := io.ReadAll(mockResponse.Body)
			if err != nil {
				t.Errorf("error with response reading: %v", err)
				return
			}
			actualStatusCode := mockResponse.StatusCode

			assert.Equal(t, tt.expected, string(actual))
			assert.Equal(t, tt.expectedStatusCode, actualStatusCode)
		})
	}
}
