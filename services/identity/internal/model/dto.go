package model

type LoginResult struct {
	User         *User
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

type TokenResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}
