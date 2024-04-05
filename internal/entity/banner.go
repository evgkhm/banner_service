package entity

type Banner struct {
	TagID     uint64 `json:"tag_id" example:"1"`
	FeatureID uint64 `json:"feature_id" example:"1"`
	Token     string `json:"token" example:"some_token"`
	Limit     int    `json:"limit" example:"10"`
	Offset    int    `json:"offset" example:"0"`
}

type UserBannerRequest struct {
	TagID          uint64 `json:"tag_id" example:"1"`
	FeatureID      uint64 `json:"feature_id" example:"1"`
	UseLastVersion bool   `json:"use_last_version" example:"true"`
	Token          string `json:"token" example:"some_token"`
}

type UserBannerResponse struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	URL   string `json:"url"`
}
