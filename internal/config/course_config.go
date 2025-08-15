package config

import (
	"LectioBot/internal/models"
	"encoding/json"
	"os"
)

func LoadCourseConfig(filename string) (*models.CourseConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg models.CourseConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func SaveCourseConfig(filename string, cfg *models.CourseConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
