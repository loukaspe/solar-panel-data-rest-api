package repositories

import (
	"errors"
	"github.com/loukaspe/solar-panel-data-crud/internal/core/domain"
	apierrors "github.com/loukaspe/solar-panel-data-crud/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestSolarPanelDataRepository_CreateSolarPanelData(t *testing.T) {
	type args struct {
		solarPanelData *domain.SolarPanelData
	}

	mockDb := SolarPanelDataDB{}

	tests := []struct {
		name                   string
		args                   args
		expectedSolarPanelData *domain.SolarPanelData
		expectError            bool
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
			expectedSolarPanelData: &domain.SolarPanelData{
				Solar: map[string][][]string{
					"uuid1": [][]string{
						{"timestamp1", "event1"},
					},
				},
				Wind: nil,
			},
			expectError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &SolarPanelDataRepository{
				db: mockDb,
			}

			actualInsertedId, err := repo.CreateSolarPanelData(tt.args.solarPanelData)
			if (err != nil) != tt.expectError {
				t.Errorf("CreateSolarPanelData() error = %v, expectError %v", err, tt.expectError)
				return
			}

			actual, ok := mockDb[actualInsertedId]
			if !ok {
				t.Errorf("data with uuid %s not found", actualInsertedId)
				return
			}

			assert.EqualValues(t, tt.expectedSolarPanelData, actual)
		})
	}
}

func TestSolarPanelDataRepository_GetSolarPanelData(t *testing.T) {
	type args struct {
		uuid string
	}

	type fields struct {
		db SolarPanelDataDB
	}

	tests := []struct {
		name          string
		args          args
		fields        fields
		expected      *domain.SolarPanelData
		expectError   bool
		expectedError *apierrors.DataNotFoundErrorWrapper
	}{
		{
			name: "get ok",
			args: args{
				uuid: "uuid",
			},
			fields: fields{
				db: SolarPanelDataDB{
					"uuid": {
						Solar: map[string][][]string{
							"uuid1": [][]string{
								{"timestamp1", "event1"},
							},
						},
						Wind: nil,
					},
				},
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
			name: "data not found",
			args: args{
				uuid: "uuidNotExisting",
			},
			fields: fields{
				db: SolarPanelDataDB{
					"uuid": {
						Solar: map[string][][]string{
							"uuid1": [][]string{
								{"timestamp1", "event1"},
							},
						},
						Wind: nil,
					},
				},
			},
			expected: &domain.SolarPanelData{},
			expectedError: &apierrors.DataNotFoundErrorWrapper{
				ReturnedStatusCode: http.StatusNoContent,
				OriginalError:      errors.New("uuid uuidNotExisting not found"),
			},
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &SolarPanelDataRepository{
				db: tt.fields.db,
			}
			actual, actualError := repo.GetSolarPanelData(tt.args.uuid)
			if (actualError != nil) != tt.expectError {
				t.Errorf("GetSolarPanelData() error = %v, expectError %v", actualError, tt.expectError)
				return
			}

			assert.Equal(t, tt.expected, actual)

			if tt.expectError {
				assert.Equal(t, tt.expectedError, actualError)
			}
		})
	}
}

func TestSolarPanelDataRepository_UpdateSolarPanelData(t *testing.T) {
	type args struct {
		uuid           string
		solarPanelData *domain.SolarPanelData
	}

	mockDb := SolarPanelDataDB{
		"uuid1": &SolarPanelData{
			Solar: map[string][][]string{
				"uuid1": [][]string{
					{"timestamp1", "event1"},
				},
			},
			Wind: nil,
		},
	}

	tests := []struct {
		name                   string
		args                   args
		expectedSolarPanelData *domain.SolarPanelData
		expectError            bool
		expectedErrorMessage   string
	}{
		{
			name: "update ok",
			args: args{
				uuid: "uuid1",
				solarPanelData: &domain.SolarPanelData{
					Solar: map[string][][]string{
						"uuid1": [][]string{
							{"timestamp2", "event2"},
						},
					},
					Wind: nil,
				},
			},
			expectedSolarPanelData: &domain.SolarPanelData{
				Solar: map[string][][]string{
					"uuid1": [][]string{
						{"timestamp2", "event2"},
					},
				},
				Wind: nil,
			},
			expectError: false,
		},
		{
			name: "error data not found",
			args: args{
				uuid: "uuidNotExisting",
				solarPanelData: &domain.SolarPanelData{
					Solar: map[string][][]string{
						"uuid1": [][]string{
							{"timestamp1", "event1"},
						},
					},
					Wind: nil,
				},
			},
			expectError: true,
			// data not found error does not have an error message
			expectedErrorMessage: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &SolarPanelDataRepository{
				db: mockDb,
			}

			err := repo.UpdateSolarPanelData(tt.args.uuid, tt.args.solarPanelData)
			if (err != nil) != tt.expectError {
				t.Errorf("CreateSolarPanelData() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if tt.expectError {
				assert.Equal(t, tt.expectedErrorMessage, err.Error())
				return
			}

			actual, ok := mockDb[tt.args.uuid]
			if !ok {
				t.Errorf("data with uuid %s not found", tt.args.uuid)
				return
			}

			assert.EqualValues(t, tt.expectedSolarPanelData, actual)
		})
	}
}

func TestSolarPanelDataRepository_DeleteSolarPanelData(t *testing.T) {
	type args struct {
		uuid string
	}

	type fields struct {
		db SolarPanelDataDB
	}

	tests := []struct {
		name     string
		args     args
		fields   fields
		expected error
	}{
		{
			name: "delete ok",
			args: args{
				uuid: "uuid",
			},
			fields: fields{
				db: SolarPanelDataDB{
					"uuid": {
						Solar: map[string][][]string{
							"uuid1": [][]string{
								{"timestamp1", "event1"},
							},
						},
						Wind: nil,
					},
				},
			},
			expected: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &SolarPanelDataRepository{
				db: tt.fields.db,
			}
			actual := repo.DeleteSolarPanelData(tt.args.uuid)

			assert.Equal(t, tt.expected, actual)
		})
	}
}
