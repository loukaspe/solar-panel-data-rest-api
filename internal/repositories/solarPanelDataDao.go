package repositories

type SolarPanelDataDB map[string]*SolarPanelData

type SolarPanelData struct {
	Solar map[string][][]string
	Wind  interface{}
}
