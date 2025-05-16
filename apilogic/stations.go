package apilogic

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type IRailStationInformation struct {
	Version string `json:"version"`
	Timestamp string `json:"timestamp"`
	Station []StationInformation `json:"station"`
}

type StationInformation struct {
	LinkID string `json:"@id"`
	ID string `json:"id"`
	Name string `json:"name"`
	LocationX string `json:"locationX"`
	LocationY string `json:"locationY"`
	StandardName string `json:"standardname"`
}

func GetStationInformation() ([]StationInformation, error) {
	baseURL := "https://api.irail.be/v1/stations/?format=json&lang=en"

	resp, err := http.Get(baseURL)
	if err != nil {
		return nil, fmt.Errorf("could not send get request to IRail Station api: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IRail Station API returned non-200 status: %s", resp.Status)
	}

	var apiResponse IRailStationInformation
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response from IRail Station API: %w", err)
	}

	return apiResponse.Station, nil
}