package models

type Config struct {
	BotToken     string
	AdminIDs     []int64
	LecturerID   int64
	CourseConfig *CourseConfig

	SpreadsheetID   string
	CredentialsFile string
}

func GetConfig(botToken string, adminIDs []int64, courseConfig *CourseConfig, lecturerID int64, spreadsheetID, credentialsFile string) *Config {
	return &Config{
		BotToken:     botToken,
		AdminIDs:     adminIDs,
		LecturerID:   lecturerID,
		CourseConfig: courseConfig,

		SpreadsheetID:   spreadsheetID,
		CredentialsFile: credentialsFile,
	}
}

func (c *Config) IsAdmin(userID int64) bool {
	for _, id := range c.AdminIDs {
		if id == userID {
			return true
		}
	}
	return false
}
