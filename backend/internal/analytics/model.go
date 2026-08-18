package analytics

import "time"

type Click struct {
	Id        int64     `json:"id"`
	UrlId     string    `json:"url_id"`
	ClickedAt time.Time `json:"clicked_at"`
	Referrer  *string   `json:"referrer"`
	UserAgent *string   `json:"user_agent"`
}

type Stats struct {
	ShortCode   string          `json:"short_code"`
	TotalClicks int64           `json:"total_clicks"`
	Timeline    []TimelinePoint `json:"timeline"`
}

type TimelinePoint struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}
