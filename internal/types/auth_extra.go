package types

// AuthRegisterRequest registers a new user account.
type AuthRegisterRequest struct {
	Email       string  `json:"email"`
	Password    string  `json:"password"`
	DisplayName *string `json:"display_name,omitempty,optional"`
	InviteCode  *string `json:"invite_code,omitempty,optional"`
}

// AuthRegisterResponse returns account info and optional tokens.
type AuthRegisterResponse struct {
	RequiresVerification bool              `json:"requires_verification"`
	AccessToken          string            `json:"access_token,omitempty"`
	RefreshToken         string            `json:"refresh_token,omitempty"`
	TokenType            string            `json:"token_type,omitempty"`
	ExpiresIn            int64             `json:"expires_in,omitempty"`
	RefreshExpiresIn     int64             `json:"refresh_expires_in,omitempty"`
	User                 AuthenticatedUser `json:"user"`
}

// AuthVerifyRequest verifies the email with a code.
type AuthVerifyRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

// AuthVerifyResponse returns issued tokens after verification.
type AuthVerifyResponse struct {
	AccessToken      string            `json:"access_token"`
	RefreshToken     string            `json:"refresh_token"`
	TokenType        string            `json:"token_type"`
	ExpiresIn        int64             `json:"expires_in"`
	RefreshExpiresIn int64             `json:"refresh_expires_in"`
	User             AuthenticatedUser `json:"user"`
}

// AuthForgotPasswordRequest requests a reset code.
type AuthForgotPasswordRequest struct {
	Email string `json:"email"`
}

// AuthForgotPasswordResponse acknowledges reset code dispatch.
type AuthForgotPasswordResponse struct {
	Message string `json:"message"`
}

// AuthResetPasswordRequest resets password using a verification code.
type AuthResetPasswordRequest struct {
	Email    string `json:"email"`
	Code     string `json:"code"`
	Password string `json:"password"`
}

// AuthResetPasswordResponse acknowledges password reset.
type AuthResetPasswordResponse struct {
	Message string `json:"message"`
}
