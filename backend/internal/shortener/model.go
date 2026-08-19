package shortener

import "time"

type URL struct {
	Id          int64     `json:"id"`
	UserId      int64     `json:"user_id"`
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
}
