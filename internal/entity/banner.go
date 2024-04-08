package entity

import "time"

type Banner struct {
	TagID     []uint64 `json:"tag_ids" example:"1"`
	FeatureID uint64   `json:"feature_id" example:"1"`
	Content   Content  `json:"content" example:"some_content"`
	IsActive  bool     `json:"is_active" example:"true"`
}

type Content struct {
	Title string `json:"title" example:"some_title"`
	Text  string `json:"text" example:"some_text"`
	URL   string `json:"url" example:"some_url"`
}

type UserBannerResponse struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	URL   string `json:"url"`
}

type BannersList struct {
	BannerID  int       `json:"banner_id"`
	TagIDs    *string   `json:"tag_ids"`
	FeatureID int       `json:"feature_id"`
	Content   Content   `json:"content"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
