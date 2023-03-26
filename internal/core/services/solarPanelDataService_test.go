package services

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/loukaspe/solar-panel-data-crud/internal/core/domain"
	mock_ports "github.com/loukaspe/solar-panel-data-crud/mocks/mock_internal/core/ports"
	apierrors "github.com/loukaspe/solar-panel-data-crud/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestSolarPanelDataService_CreateSolarPanelData(t *testing.T) {
	type args struct {
		solarPanelData *domain.SolarPanelData
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepository := mock_ports.NewMockSolarPanelDataRepositoryInterface(mockCtrl)

	tests := []struct {
		name                      string
		args                      args
		shouldMockRepositoryRun   bool
		mockRepositoryReturnError error
		expectedUuid              string
		expectedErrorMessage      string
		expectError               bool
	}{
		{
			name: "creation ok",
			args: args{
				solarPanelData: &domain.SolarPanelData{
					Solar: map[string][][]string{
						"uuid1": [][]string{
							{"timestamp1", "event1"},
						},
					},
					Wind: nil,
				},
			},
			expectedUuid:              "newUuid",
			shouldMockRepositoryRun:   true,
			mockRepositoryReturnError: nil,
			expectedErrorMessage:      "",
			expectError:               false,
		},
		{
			name: "empty solar data",
			args: args{
				solarPanelData: &domain.SolarPanelData{
					Wind: nil,
				},
			},
			expectedUuid:              "",
			shouldMockRepositoryRun:   false,
			mockRepositoryReturnError: nil,
			expectedErrorMessage:      "solar data is empty on request",
			expectError:               true,
		},
		{
			name: "repo returns error",
			args: args{
				solarPanelData: &domain.SolarPanelData{
					Solar: map[string][][]string{
						"uuid1": [][]string{
							{"timestamp1", "event1"},
						},
					},
					Wind: nil,
				},
			},
			shouldMockRepositoryRun:   true,
			mockRepositoryReturnError: errors.New("random error"),
			expectedErrorMessage:      "random error",
			expectError:               true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := SolarPanelDataService{
				repository: mockRepository,
			}

			if tt.shouldMockRepositoryRun {
				mockRepository.EXPECT().
					CreateSolarPanelData(tt.args.solarPanelData).
					Return(tt.expectedUuid, tt.mockRepositoryReturnError)
			}

			insertedUuid, actualError := service.CreateSolarPanelData(tt.args.solarPanelData)
			if (actualError != nil) != tt.expectError {
				t.Errorf("CreateSolarPanelData() error = %v, expectError %v", actualError, tt.expectError)
				return
			}

			assert.Equal(t, insertedUuid, tt.expectedUuid)

			if tt.expectError {
				assert.Equal(t, tt.expectedErrorMessage, actualError.Error())
			}
		})
	}
}

func TestSolarPanelDataService_GetSolarPanelData(t *testing.T) {
	type args struct {
		uuid string
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepository := mock_ports.NewMockSolarPanelDataRepositoryInterface(mockCtrl)

	tests := []struct {
		name                      string
		args                      args
		mockRepositoryReturnData  *domain.SolarPanelData
		mockRepositoryReturnError error
		expected                  *domain.SolarPanelData
		expectedErrorMessage      string
		expectError               bool
	}{
		{
			name: "get ok",
			args: args{
				uuid: "uuid",
			},
			mockRepositoryReturnData: &domain.SolarPanelData{
				Solar: map[string][][]string{
					"uuid1": [][]string{
						{"timestamp1", "event1"},
					},
				},
				Wind: nil,
			},
			expected: &domain.SolarPanelData{
				Solar: map[string][][]string{
					"uuid1": [][]string{
						{"timestamp1", "event1"},
					},
				},
				Wind: nil,
			},
			expectError: false,
		},
		{
			name: "repo random error",
			args: args{
				uuid: "uuid",
			},
			mockRepositoryReturnData:  &domain.SolarPanelData{},
			mockRepositoryReturnError: errors.New("random error"),
			expected:                  &domain.SolarPanelData{},
			expectedErrorMessage:      "random error",
			expectError:               true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := SolarPanelDataService{
				repository: mockRepository,
			}

			mockRepository.EXPECT().
				GetSolarPanelData(tt.args.uuid).
				Return(tt.mockRepositoryReturnData, tt.mockRepositoryReturnError)

			actual, actualError := service.GetSolarPanelData(tt.args.uuid)
			if (actualError != nil) != tt.expectError {
				t.Errorf("GetSolarPanelData() error = %v, expectError %v", actualError, tt.expectError)
				return
			}

			assert.Equal(t, tt.expected, actual)

			if tt.expectError {
				assert.Equal(t, tt.expectedErrorMessage, actualError.Error())
			}
		})
	}
}

func TestSolarPanelDataService_UpdateSolarPanelData(t *testing.T) {
	type args struct {
		uuid           string
		solarPanelData *domain.SolarPanelData
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepository := mock_ports.NewMockSolarPanelDataRepositoryInterface(mockCtrl)

	tests := []struct {
		name                      string
		args                      args
		shouldMockRepositoryRun   bool
		mockRepositoryReturnError error
		expected                  error
	}{
		{
			name: "update ok",
			args: args{
				solarPanelData: &domain.SolarPanelData{
					Solar: map[string][][]string{
						"uuid1": [][]string{
							{"timestamp1", "event1"},
						},
					},
					Wind: nil,
				},
			},
			shouldMockRepositoryRun:   true,
			mockRepositoryReturnError: nil,
			expected:                  nil,
		},
		{
			name: "invalid empty request",
			args: args{
				solarPanelData: &domain.SolarPanelData{
					Wind: nil,
				},
			},
			shouldMockRepositoryRun: false,
			expected: apierrors.EmptySolarDataError{
				ReturnedStatusCode: http.StatusBadRequest,
			},
		},
		{
			name: "repo returns error",
			args: args{
				solarPanelData: &domain.SolarPanelData{
					Solar: map[string][][]string{
						"uuid1": [][]string{
							{"timestamp1", "event1"},
						},
					},
					Wind: nil,
				},
			},
			shouldMockRepositoryRun:   true,
			mockRepositoryReturnError: errors.New("random error"),
			expected:                  errors.New("random error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := SolarPanelDataService{
				repository: mockRepository,
			}

			if tt.shouldMockRepositoryRun {
				mockRepository.EXPECT().
					UpdateSolarPanelData(tt.args.uuid, tt.args.solarPanelData).
					Return(tt.mockRepositoryReturnError)
			}

			actual := service.UpdateSolarPanelData(tt.args.uuid, tt.args.solarPanelData)

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestSolarPanelDataService_DeleteSolarPanelData(t *testing.T) {
	type args struct {
		uuid string
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepository := mock_ports.NewMockSolarPanelDataRepositoryInterface(mockCtrl)

	tests := []struct {
		name                      string
		args                      args
		mockRepositoryReturnError error
		expected                  error
	}{
		{
			name: "delete ok",
			args: args{
				uuid: "uuid",
			},
			mockRepositoryReturnError: nil,
			expected:                  nil,
		},
		{
			name: "repo random error",
			args: args{
				uuid: "uuid",
			},
			mockRepositoryReturnError: errors.New("random error"),
			expected:                  errors.New("random error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := SolarPanelDataService{
				repository: mockRepository,
			}

			mockRepository.EXPECT().
				DeleteSolarPanelData(tt.args.uuid).
				Return(tt.mockRepositoryReturnError)

			actual := service.DeleteSolarPanelData(tt.args.uuid)

			assert.Equal(t, tt.expected, actual)
		})
	}
}
