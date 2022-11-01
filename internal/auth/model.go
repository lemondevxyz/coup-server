package auth

type Provider interface {
	ID() string
	URL() string
	// Optional
	Profiler() (Profiler, error)
}

type Profiler interface {
	GetProfileByID(id string) (Profile, error)
}

type Profile struct {
	ID         string `json:"id"`
	PictureURL string `json:"picture_url"`
	Username   string `json:"username"`
	Platform   string `json:"platform"`
}
