package helper

import (
	"errors"
	"github.com/loukaspe/solar-panel-data-crud/internal/core/domain"
	apierrors "github.com/loukaspe/solar-panel-data-crud/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestSolarPanelDataEventExtractor_ExtractEventsPerParameterIdToCsvForm(t *testing.T) {
	type args struct {
		solarPanelData *domain.SolarPanelData
	}
	tests := []struct {
		name          string
		args          args
		expected      [][]string
		expectError   bool
		expectedError error
	}{
		{
			name: "valid one parameter one event",
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
			expected: [][]string{
				{"Events"},
				{"event1"},
			},
			expectError: false,
		},
		{
			name: "valid one parameter multiple events",
			args: args{
				solarPanelData: &domain.SolarPanelData{
					Solar: map[string][][]string{
						"uuid1": [][]string{
							{"timestamp1", "event1"},
							{"timestamp2", "event2"},
							{"timestamp3", "event3"},
							{"timestamp4", "event4"},
						},
					},
					Wind: nil,
				},
			},
			expected: [][]string{
				{"Events"},
				{"event1"},
				{"event2"},
				{"event3"},
				{"event4"},
			},
			expectError: false,
		},
		{
			name: "valid multiple parameter one event",
			args: args{
				solarPanelData: &domain.SolarPanelData{
					Solar: map[string][][]string{
						"uuid1": [][]string{
							{"timestamp1", "event1"},
						},
						"uuid2": [][]string{
							{"timestamp1", "event1"},
						},
						"uuid3": [][]string{
							{"timestamp1", "event1"},
						},
						"uuid4": [][]string{
							{"timestamp1", "event1"},
						},
					},
					Wind: nil,
				},
			},
			expected: [][]string{
				{"Events"},
				{"event1"},
				{"event1"},
				{"event1"},
				{"event1"},
			},
			expectError: false,
		},
		{
			name: "valid multiple parameters multiple events",
			args: args{
				solarPanelData: &domain.SolarPanelData{
					Solar: map[string][][]string{
						"uuid1": [][]string{
							{"timestamp1", "event1"},
							{"timestamp2", "event2"},
							{"timestamp3", "event3"},
						},
						"uuid2": [][]string{
							{"timestamp1", "event1"},
							{"timestamp2", "event2"},
							{"timestamp3", "event3"},
						},
						"uuid3": [][]string{
							{"timestamp1", "event1"},
							{"timestamp2", "event2"},
							{"timestamp3", "event3"},
						},
					},
					Wind: nil,
				},
			},
			expected: [][]string{
				{"Events"},
				{"event1"},
				{"event2"},
				{"event3"},
				{"event1"},
				{"event2"},
				{"event3"},
				{"event1"},
				{"event2"},
				{"event3"},
			},
			expectError: false,
		},
		{
			name: "invalid missing event",
			args: args{
				solarPanelData: &domain.SolarPanelData{
					Solar: map[string][][]string{
						"uuid1": [][]string{
							{"timestamp1", "event1"},
							{"timestamp2", "event2"},
							{"timestamp3", "event3"},
						},
						"uuid2": [][]string{
							{"timestamp1", "event1"},
							{"timestamp2"},
							{"timestamp3", "event3"},
						},
						"uuid3": [][]string{
							{"timestamp1", "event1"},
							{"timestamp2", "event2"},
							{"timestamp3", "event3"},
						},
					},
					Wind: nil,
				},
			},
			expected:    [][]string{},
			expectError: true,
			expectedError: apierrors.MalformedEventDataError{
				ReturnedStatusCode:   http.StatusInternalServerError,
				MalformedParameterId: "uuid2",
				OriginalError:        errors.New("parameterId uuid2 contains events with no values"),
			},
		},
		{
			name: "invalid empty string event",
			args: args{
				solarPanelData: &domain.SolarPanelData{
					Solar: map[string][][]string{
						"uuid1": [][]string{
							{"timestamp1", "event1"},
							{"timestamp2", "event2"},
							{"timestamp3", "event3"},
						},
						"uuid2": [][]string{
							{"timestamp1", "event1"},
							{"timestamp2", ""},
							{"timestamp3", "event3"},
						},
						"uuid3": [][]string{
							{"timestamp1", "event1"},
							{"timestamp2", "event2"},
							{"timestamp3", "event3"},
						},
					},
					Wind: nil,
				},
			},
			expected:    [][]string{},
			expectError: true,
			expectedError: apierrors.MalformedEventDataError{
				ReturnedStatusCode:   http.StatusInternalServerError,
				MalformedParameterId: "uuid2",
				OriginalError:        errors.New("parameterId uuid2 contains events with no values"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extractor := SolarPanelDataEventExtractor{}
			actual, actualError := extractor.ExtractEventsPerParameterIdToCsvForm(tt.args.solarPanelData)
			if (actualError != nil) != tt.expectError {
				t.Errorf("ExtractEventsPerParameterIdToCsvForm() error = %v, expectedError %v", actualError, tt.expectedError)
				return
			}

			assert.EqualValues(t, tt.expected, actual)
			if tt.expectError {
				assert.Equal(t, tt.expectedError, actualError)
			}
		})
	}
}
