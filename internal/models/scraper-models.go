package models

type GameDurations struct {
	MainStory     string `json:"main_story"`
	MainsSides    string `json:"main_sides"`
	Completionist string `json:"completionist"`
}

type ScraperSuccessPayload struct {
	Query         string        `json:"query"`
	GameTitle     string        `json:"game_title"`
	GameId        string        `json:"game_id"`
	GameDurations GameDurations `json:"game_durations"`
}

type ScraperSuccessResponse struct {
	Message string                `json:"message"`
	Data    ScraperSuccessPayload `json:"data"`
}
