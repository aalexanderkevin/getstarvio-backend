package auth

type GoogleLoginRequest struct {
	IDToken   string `json:"idToken"`
	GoogleSub string `json:"googleSub"`
	Email     string `json:"email"`
	Name      string `json:"name"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	UserID       string `json:"userId"`
}

type googleIdentity struct {
	Sub   string
	Email string
	Name  string
}
