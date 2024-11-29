package models

type LocalLoginRequest struct {
	ClientId string `json:"clientId" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LocalForgotPasswordRequest struct {
	ClientId string `json:"clientId" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type LocalResetPasswordRequest struct {
	Token           string `json:"token" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
}

type LocalLoginResponse struct {
	ID           string `json:"id,omitempty"`
	AccountID    string `json:"accountId,omitempty"`
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}
