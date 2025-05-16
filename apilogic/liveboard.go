package apilogic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type BaseLiveboardAPIResponse struct {
	Version string `json:"version"`
	Timestamp string `json:"timestamp"`
	Station string `json:"station"`
	StationInfo StationInformation
	Departures Departures `json:"departures"`
}

type Departures struct {
	Number string `json:"number"`
	DepartureInfo []DepartureInfo `json:"departure"`
}

type DepartureInfo struct {
	ID string `json:"id"`
	StationName string `json:"station"`
	StationInfo StationInformation `json:"stationinfo"`
	Time string `json:"time"`
	Delay string `json:"delay"`
	Canceled string `json:"canceled"`
	Left string `json:"left"`
	Vehicle string `json:"vehicle"`
	VehicleInfo VehicleInfo `json:"vehicleinfo"`
	Platform string `json:"platform"`
	PlatformInfo PlatformInfo `json:"platforminfo"`
	Occupancy Occupancy `json:"occupancy"`
	DepartureConnection string `json:"departureConnection"`
}

type VehicleInfo struct {
	Name string `json:"name"`
	Shortname string `json:"shortname"`
	Number string `json:"number"`
	Type string `json:"type"`
	LinkID string `json:"@id"`
}

type PlatformInfo struct {
	Name string `json:"name"`
	Normal string `json:"normal"`
}

type Occupancy struct {
	LinkID string `json:"@id"`
	Name string `json:"name"`
}

// FormattedTime converts the Unix timestamp to a human-readable time in CET/CEST.
func (d DepartureInfo) FormattedTime() string {
	// Parse the timestamp string to an integer
	timestamp, err := strconv.ParseInt(d.Time, 10, 64)
	if err != nil {
		return "Invalid time"
	}

	// Convert to time.Time in UTC
	t := time.Unix(timestamp, 0)

	// Load the Belgian timezone (CET/CEST)
	loc, err := time.LoadLocation("Europe/Brussels")
	if err != nil {
		return t.Format("15:04") // Fallback to UTC if timezone loading fails
	}

	// Convert to Belgian timezone and format as HH:MM (24-hour)
	return t.In(loc).Format("15:04")
}

func GetLiveboardInformation(station string) (*BaseLiveboardAPIResponse, error) {
	baseURL := fmt.Sprintf("https://api.irail.be/liveboard/?station=%s&format=json&lang=en", station)

	resp, err := http.Get(baseURL)
	if err != nil {
		return nil, fmt.Errorf("could not send GET request to iRail Liveboard API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("iRail Liveboard API returned non-200 status: %s", resp.Status)
	}

	var apiResponse BaseLiveboardAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response from iRail Liveboard API: %w", err)
	}

	return &apiResponse, nil
}