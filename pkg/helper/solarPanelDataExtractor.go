package helper

import (
	"errors"
	"github.com/loukaspe/solar-panel-data-crud/internal/core/domain"
	apierrors "github.com/loukaspe/solar-panel-data-crud/pkg/errors"
	"net/http"
)

type SolarPanelDataEventExtractorInterface interface {
	ExtractEventsPerParameterIdToCsvForm(*domain.SolarPanelData) ([][]string, error)
}

type SolarPanelDataEventExtractor struct{}

func NewSolarPanelDataEventExtractor() *SolarPanelDataEventExtractor {
	return &SolarPanelDataEventExtractor{}
}

// ExtractEventsPerParameterId creates a 2-dimensional array of string that holds
// an event string in every row of the array. It's meant to be used for creating a
// CSV
func (extractor SolarPanelDataEventExtractor) ExtractEventsPerParameterIdToCsvForm(
	solarPanelData *domain.SolarPanelData,
) ([][]string, error) {
	const EventArrayElementNormalSize = 2

	var SolarPanelDataEvents [][]string
	var eventValue string

	csvHeaderRow := []string{"Events"}

	SolarPanelDataEvents = append(
		SolarPanelDataEvents,
		csvHeaderRow,
	)

	for parameterId, parameterIdEvents := range solarPanelData.Solar {
		for _, event := range parameterIdEvents {
			if len(event) < EventArrayElementNormalSize {
				return [][]string{}, apierrors.MalformedEventDataError{
					ReturnedStatusCode:   http.StatusInternalServerError,
					MalformedParameterId: parameterId,
					OriginalError:        errors.New("parameterId " + parameterId + " contains events with no values"),
				}
			}

			eventValue = event[1]

			if eventValue == "" {
				return [][]string{}, apierrors.MalformedEventDataError{
					ReturnedStatusCode:   http.StatusInternalServerError,
					MalformedParameterId: parameterId,
					OriginalError:        errors.New("parameterId " + parameterId + " contains events with no values"),
				}
			}

			SolarPanelDataEvents = append(
				SolarPanelDataEvents,
				[]string{eventValue},
			)
		}
	}

	return SolarPanelDataEvents, nil
}
