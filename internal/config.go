package internal

import "os"

type Config struct {
	AllowedAPIKeys map[string]string
}

func LoadConfig() *Config {
	return &Config{
		AllowedAPIKeys: map[string]string{
			"website": os.Getenv("WEBSITE_API_KEY"),
			"bot":     os.Getenv("BOT_API_KEY"),
		},
	}
}
