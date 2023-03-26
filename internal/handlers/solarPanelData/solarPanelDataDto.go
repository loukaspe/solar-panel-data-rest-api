package solarPanelData

type Dto struct {
	Solar map[string][][]string `json:"solar"`
	Wind  interface{}           `json:"wind"`
}

type CreateSolarPanelDataResponse struct {
	InsertedId    string `json:"id,omitempty"`
	DataSubmitted *Dto   `json:"dataSubmitted,omitempty"`
	ErrorMessage  string `json:"errorMessage,omitempty"`
}

type GetSolarPanelDataResponse struct {
	SolarPanelDataEvents [][]string
}

type DeleteSolarPanelDataResponse struct {
	ErrorMessage string `json:"errorMessage,omitempty"`
}
