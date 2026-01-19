package model

type LoginResult struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
}

type TokenResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}
