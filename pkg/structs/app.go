package structs

// App struct represents the request payload for creating an app
type App struct {
	App           string  `json:"app" validate:"required"`
	Category      string  `json:"category" validate:"required"`
	Rating        float64 `json:"rating" validate:"required"`
	Reviews       int     `json:"reviews" validate:"required"`
	Size          string  `json:"size" validate:"required"`
	Installs      string  `json:"installs" validate:"required"`
	Type          string  `json:"type" validate:"required"`
	Price         string  `json:"price" validate:"required"`
	ContentRating string  `json:"content_rating" validate:"required"`
	Genres        string  `json:"genres" validate:"required"`
	LastUpdated   string  `json:"last_updated" validate:"required"`
	CurrentVer    string  `json:"current_ver" validate:"required"`
	AndroidVer    string  `json:"android_ver" validate:"required"`
}
