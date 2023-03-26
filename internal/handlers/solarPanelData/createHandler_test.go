package solarPanelData

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/loukaspe/solar-panel-data-crud/internal/core/domain"
	mock_services "github.com/loukaspe/solar-panel-data-crud/mocks/mock_internal/core/services"
	apierrors "github.com/loukaspe/solar-panel-data-crud/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateSolarPanelDataHandler_CreateSolarPanelDataController(t *testing.T) {
	logger := logrus.New()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockService := mock_services.NewMockSolarPanelDataServiceInterface(mockCtrl)

	tests := []struct {
		name            string
		requestBody     []byte
		mockRequestData *domain.SolarPanelData
		// variable to check if the handler returns error before the mock service runs
		shouldMockServiceRun     bool
		mockServiceResponseUuid  string
		mockServiceResponseError error
		expected                 []byte
		expectedStatusCode       int
	}{
		{
			name: "valid",
			requestBody: json.RawMessage(`{
  "solar": {
    "38d503e5-dc1c-4549-8172-09d9c29070f7": [
      [
        "20211231T221500Z",
        "0.0"
      ]
    ]
  },
  "wind": null
}`),
			mockRequestData: &domain.SolarPanelData{
				Solar: map[string][][]string{
					"38d503e5-dc1c-4549-8172-09d9c29070f7": [][]string{
						{"20211231T221500Z", "0.0"},
					},
				},
				Wind: nil,
			},
			shouldMockServiceRun:     true,
			mockServiceResponseUuid:  "newUuid",
			mockServiceResponseError: nil,
			expected: json.RawMessage(`{"id":"newUuid","dataSubmitted":{"solar":{"38d503e5-dc1c-4549-8172-09d9c29070f7":[["20211231T221500Z","0.0"]]},"wind":null}}
`),
			expectedStatusCode: 201,
		},
		{
			name: "invalid bad request",
			requestBody: json.RawMessage(`{
  "solar": 
    "38d503e5-dc1c-4549-8172-09d9c29070f7": [
      [
        "20211231T221500Z",
        "0.0"
      ]
    ]
  },
  "wind": null
}`),
			shouldMockServiceRun: false,
			expected: json.RawMessage(`{"errorMessage":"malformed solar panel data request"}
`),
			expectedStatusCode: 400,
		},
		{
			name:                 "invalid empty request",
			requestBody:          json.RawMessage(``),
			shouldMockServiceRun: false,
			expected: json.RawMessage(`{"errorMessage":"malformed solar panel data request"}
`),
			expectedStatusCode: 400,
		},
		{
			name: "invalid service empty solar data error",
			requestBody: json.RawMessage(`{
  "wind": null
}`),
			mockRequestData: &domain.SolarPanelData{
				Wind: nil,
			},
			shouldMockServiceRun:    true,
			mockServiceResponseUuid: "",
			mockServiceResponseError: apierrors.EmptySolarDataError{
				ReturnedStatusCode: http.StatusBadRequest,
			},
			expected: json.RawMessage(`{"errorMessage":"solar data is empty on request"}
`),
			expectedStatusCode: 400,
		},
		{
			name: "invalid service error",
			requestBody: json.RawMessage(`{
  "solar": {
    "38d503e5-dc1c-4549-8172-09d9c29070f7": [
      [
        "20211231T221500Z",
        "0.0"
      ]
    ]
  },
  "wind": null
}`),
			mockRequestData: &domain.SolarPanelData{
				Solar: map[string][][]string{
					"38d503e5-dc1c-4549-8172-09d9c29070f7": [][]string{
						{"20211231T221500Z", "0.0"},
					},
				},
				Wind: nil,
			},
			shouldMockServiceRun:     true,
			mockServiceResponseUuid:  "",
			mockServiceResponseError: errors.New("random error"),
			expected: json.RawMessage(`{"errorMessage":"random error"}
`),
			expectedStatusCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBodyReader := bytes.NewBuffer(tt.requestBody)

			mockRequest := httptest.NewRequest("POST", "/solarPanelData", requestBodyReader)
			mockRequest.Header.Set("Content-Type", "application/json")
			mockResponseRecorder := httptest.NewRecorder()

			if tt.shouldMockServiceRun {
				mockService.EXPECT().
					CreateSolarPanelData(tt.mockRequestData).
					Return(tt.mockServiceResponseUuid, tt.mockServiceResponseError)
			}

			handler := &CreateSolarPanelDataHandler{
				SolarPanelDataService: mockService,
				logger:                logger,
			}
			sut := handler.CreateSolarPanelDataController

			sut(mockResponseRecorder, mockRequest)

			mockResponse := mockResponseRecorder.Result()
			actual, err := io.ReadAll(mockResponse.Body)
			if err != nil {
				t.Errorf("error with response reading: %v", err)
				return
			}
			actualStatusCode := mockResponse.StatusCode

			assert.Equal(t, string(tt.expected), string(actual))
			assert.Equal(t, tt.expectedStatusCode, actualStatusCode)
		})
	}
}
