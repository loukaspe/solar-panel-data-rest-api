package solarPanelData

import (
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	mock_services "github.com/loukaspe/solar-panel-data-crud/mocks/mock_internal/core/services"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"testing"
)

func TestDeleteSolarPanelDataHandler_DeleteSolarPanelDataController(t *testing.T) {
	logger := logrus.New()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockService := mock_services.NewMockSolarPanelDataServiceInterface(mockCtrl)

	tests := []struct {
		name                     string
		requestedUuid            string
		shouldMockServiceRun     bool
		mockServiceResponseError error
		expected                 []byte
		expectedStatusCode       int
	}{
		{
			name:                     "valid",
			requestedUuid:            "38d503e5-dc1c-4549-8172-09d9c29070f7",
			shouldMockServiceRun:     true,
			mockServiceResponseError: nil,
			expected:                 json.RawMessage(``),
			expectedStatusCode:       200,
		},
		{
			name:                 "missing id",
			requestedUuid:        "",
			shouldMockServiceRun: false,
			expected: json.RawMessage(`{"errorMessage":"missing solarPanelData id"}
`),
			expectedStatusCode: 400,
		},
		{
			name:                     "invalid service random error",
			requestedUuid:            "aaaaaa",
			shouldMockServiceRun:     true,
			mockServiceResponseError: errors.New("random error"),
			expected: json.RawMessage(`{"errorMessage":"random error"}
`),
			expectedStatusCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRequest := httptest.NewRequest(
				"DELETE",
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
					DeleteSolarPanelData(tt.requestedUuid).
					Return(tt.mockServiceResponseError)
			}

			handler := &DeleteSolarPanelDataHandler{
				SolarPanelDataService: mockService,
				logger:                logger,
			}
			sut := handler.DeleteSolarPanelDataController

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
