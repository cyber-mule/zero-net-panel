### 1. "Request password reset code"

1. route definition

- Url: /api/v1/auth/forgot
- Method: POST
- Request: `AuthForgotPasswordRequest`
- Response: `AuthForgotPasswordResponse`

2. request definition



```golang
type AuthForgotPasswordRequest struct {
	Email string 
}
```


3. response definition



```golang
type AuthForgotPasswordResponse struct {
	Message string 
}
```

### 2. "Authenticate user and issue token"

1. route definition

- Url: /api/v1/auth/login
- Method: POST
- Request: `AuthLoginRequest`
- Response: `AuthLoginResponse`

2. request definition



```golang
type AuthLoginRequest struct {
	Email string 
	Password string 
}
```


3. response definition



```golang
type AuthLoginResponse struct {
	Access_token string 
	Refresh_token string 
	Token_type string 
	Expires_in int64 
	Refresh_expires_in int64 
	User AuthenticatedUser 
}

type AuthenticatedUser struct {
	Id uint64 
	Email string 
	Display_name string 
	Roles []string 
	Created_at int64 
	Updated_at int64 
}
```

### 3. "Refresh access token"

1. route definition

- Url: /api/v1/auth/refresh
- Method: POST
- Request: `AuthRefreshRequest`
- Response: `AuthRefreshResponse`

2. request definition



```golang
type AuthRefreshRequest struct {
	Refresh_token string 
}
```


3. response definition



```golang
type AuthRefreshResponse struct {
	Access_token string 
	Refresh_token string 
	Token_type string 
	Expires_in int64 
	Refresh_expires_in int64 
	User AuthenticatedUser 
}

type AuthenticatedUser struct {
	Id uint64 
	Email string 
	Display_name string 
	Roles []string 
	Created_at int64 
	Updated_at int64 
}
```

### 4. "Register new user account"

1. route definition

- Url: /api/v1/auth/register
- Method: POST
- Request: `AuthRegisterRequest`
- Response: `AuthRegisterResponse`

2. request definition



```golang
type AuthRegisterRequest struct {
	Email string 
	Password string 
	Display_name string `form:"display_name,optional" json:"display_name,optional"`
	Invite_code string `form:"invite_code,optional" json:"invite_code,optional"`
}
```


3. response definition



```golang
type AuthRegisterResponse struct {
	Requires_verification bool 
	Access_token string `form:"access_token,optional" json:"access_token,optional"`
	Refresh_token string `form:"refresh_token,optional" json:"refresh_token,optional"`
	Token_type string `form:"token_type,optional" json:"token_type,optional"`
	Expires_in int64 `form:"expires_in,optional" json:"expires_in,optional"`
	Refresh_expires_in int64 `form:"refresh_expires_in,optional" json:"refresh_expires_in,optional"`
	User AuthenticatedUser 
}

type AuthenticatedUser struct {
	Id uint64 
	Email string 
	Display_name string 
	Roles []string 
	Created_at int64 
	Updated_at int64 
}
```

### 5. "Reset password using verification code"

1. route definition

- Url: /api/v1/auth/reset
- Method: POST
- Request: `AuthResetPasswordRequest`
- Response: `AuthResetPasswordResponse`

2. request definition



```golang
type AuthResetPasswordRequest struct {
	Email string 
	Code string 
	Password string 
}
```


3. response definition



```golang
type AuthResetPasswordResponse struct {
	Message string 
}
```

### 6. "Verify user email with code"

1. route definition

- Url: /api/v1/auth/verify
- Method: POST
- Request: `AuthVerifyRequest`
- Response: `AuthVerifyResponse`

2. request definition



```golang
type AuthVerifyRequest struct {
	Email string 
	Code string 
}
```


3. response definition



```golang
type AuthVerifyResponse struct {
	Access_token string 
	Refresh_token string 
	Token_type string 
	Expires_in int64 
	Refresh_expires_in int64 
	User AuthenticatedUser 
}

type AuthenticatedUser struct {
	Id uint64 
	Email string 
	Display_name string 
	Roles []string 
	Created_at int64 
	Updated_at int64 
}
```

