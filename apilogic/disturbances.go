package apilogic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type BaseDisturbancesAPIResponse struct {
	Version string `json:"version"`
	Timestamp string `json:"timestamp"`
	Disturbance []Disturbance `json:"disturbance"`
}

type Disturbance struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Link string `json:"link"`
	Time string `json:"timestamp"`
}

// FormattedTime converts the Unix timestamp to a human-readable time in CET/CEST.
func (d Disturbance) FormattedTime() string {
	timestamp, err := strconv.ParseInt(d.Time, 10, 64)
	if err != nil {
		return "Invalid time"
	}
	t := time.Unix(timestamp, 0)
	loc, err := time.LoadLocation("Europe/Brussels")
	if err != nil {
		return t.Format("15:04") // Fallback to UTC
	}
	return t.In(loc).Format("15:04")
}

// IsCurrentDate checks if the disturbance is from the current date in CET/CEST.
func (d Disturbance) IsCurrentDate() bool {
	timestamp, err := strconv.ParseInt(d.Time, 10, 64)
	if err != nil {
		return false
	}
	t := time.Unix(timestamp, 0)
	loc, err := time.LoadLocation("Europe/Brussels")
	if err != nil {
		return false
	}
	currentDate := time.Now().In(loc).Truncate(24 * time.Hour)
	disturbanceDate := t.In(loc).Truncate(24 * time.Hour)
	return currentDate.Equal(disturbanceDate)
}

func GetDisturbanceInformation() ([]Disturbance, error) {
	baseURL := "https://api.irail.be/v1/disturbances/?format=json&lineBreakCharacter=%27%27&lang=en"

	resp, err := http.Get(baseURL)
	if err != nil {
		return nil, fmt.Errorf("could not send get request to IRail Disturbance api: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IRail Disturbance API returned non-200 status: %s", resp.Status)
	}

	var apiResponse BaseDisturbancesAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response from IRail Disturbance API: %w", err)
	}

	// Filter disturbances for the current date
	var currentDisturbances []Disturbance
	for _, disturbance := range apiResponse.Disturbance {
		if disturbance.IsCurrentDate() {
			currentDisturbances = append(currentDisturbances, disturbance)
		}
	}

	return currentDisturbances, nil
}