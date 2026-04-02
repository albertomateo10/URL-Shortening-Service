package model

import "time"

type ClickEvent struct {
	ID        int64     `json:"id"`
	URLID     int64     `json:"url_id"`
	ClickedAt time.Time `json:"clicked_at"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Referer   string    `json:"referer"`
	Country   string    `json:"country"`
}

type DailyClickCount struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type ClicksOverTimeResponse struct {
	URLID       int64             `json:"url_id"`
	Period      string            `json:"period"`
	TotalClicks int               `json:"total_clicks"`
	ClicksPerDay []DailyClickCount `json:"clicks_per_day"`
}

type BrowserCount struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type CountryCount struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type SourcesResponse struct {
	URLID     int64          `json:"url_id"`
	Period    string         `json:"period"`
	Browsers  []BrowserCount `json:"browsers"`
	Countries []CountryCount `json:"countries"`
}
