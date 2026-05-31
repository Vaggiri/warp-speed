package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Record struct {
	Date     time.Time `json:"date"`
	Server   string    `json:"server"`
	Latency  float64   `json:"latency"` // ms
	Download float64   `json:"download"` // Mbps
	Upload   float64   `json:"upload"` // Mbps
}

func getHistoryPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "history.json"
	}
	dir := filepath.Join(home, ".warp-speed")
	os.MkdirAll(dir, 0755)
	return filepath.Join(dir, "history.json")
}

func SaveRecord(r Record) error {
	records, err := LoadRecords()
	if err != nil {
		records = []Record{}
	}

	records = append(records, r)

	// Keep last 100
	if len(records) > 100 {
		records = records[len(records)-100:]
	}

	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(getHistoryPath(), data, 0644)
}

func LoadRecords() ([]Record, error) {
	data, err := os.ReadFile(getHistoryPath())
	if err != nil {
		return nil, err
	}

	var records []Record
	if len(data) > 0 {
		err = json.Unmarshal(data, &records)
		if err != nil {
			return nil, err
		}
	}
	return records, nil
}
