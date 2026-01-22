### 1. "Get user balance"

1. route definition

- Url: /api/v1/user/account/balance
- Method: GET
- Request: `UserBalanceRequest`
- Response: `UserBalanceResponse`

2. request definition



```golang
type UserBalanceRequest struct {
	Page int `form:"page,optional" json:"page,optional"`
	Per_page int `form:"per_page,optional" json:"per_page,optional"`
	Entry_type string `form:"entry_type,optional" json:"entry_type,optional"`
}
```


3. response definition



```golang
type UserBalanceResponse struct {
	User_id uint64 
	Balance_cents int64 
	Currency string 
	Updated_at int64 
	Transactions []BalanceTransactionSummary 
	Pagination PaginationMeta 
}

type PaginationMeta struct {
	Page int 
	Per_page int 
	Total_count int64 
	Has_next bool 
	Has_prev bool 
}
```

### 2. "Rotate credential"

1. route definition

- Url: /api/v1/user/account/credentials/rotate
- Method: POST
- Request: `UserRotateCredentialRequest`
- Response: `UserRotateCredentialResponse`

2. request definition



```golang
type UserRotateCredentialRequest struct {
}
```


3. response definition



```golang
type UserRotateCredentialResponse struct {
	Credential CredentialSummary 
}

type CredentialSummary struct {
	Version int 
	Status int 
	Issued_at int64 
	Deprecated_at *int64 
	Revoked_at *int64 
	Last_seen_at *int64 
}
```

### 3. "Change email"

1. route definition

- Url: /api/v1/user/account/email
- Method: POST
- Request: `UserChangeEmailRequest`
- Response: `UserChangeEmailResponse`

2. request definition



```golang
type UserChangeEmailRequest struct {
	Email string 
	Code string 
	Password string 
}
```


3. response definition



```golang
type UserChangeEmailResponse struct {
	Profile UserProfile 
}

type UserProfile struct {
	Id uint64 
	Email string 
	Display_name string 
	Status int 
	Email_verified_at *int64 
	Created_at int64 
	Updated_at int64 
}
```

### 4. "Send email change code"

1. route definition

- Url: /api/v1/user/account/email/code
- Method: POST
- Request: `UserEmailChangeCodeRequest`
- Response: `UserEmailChangeCodeResponse`

2. request definition



```golang
type UserEmailChangeCodeRequest struct {
	Email string 
}
```


3. response definition



```golang
type UserEmailChangeCodeResponse struct {
	Message string 
}
```

### 5. "Change password"

1. route definition

- Url: /api/v1/user/account/password
- Method: POST
- Request: `UserChangePasswordRequest`
- Response: `UserChangePasswordResponse`

2. request definition



```golang
type UserChangePasswordRequest struct {
	Current_password string 
	New_password string 
}
```


3. response definition



```golang
type UserChangePasswordResponse struct {
	Message string 
}
```

### 6. "Get user profile"

1. route definition

- Url: /api/v1/user/account/profile
- Method: GET
- Request: `UserProfileRequest`
- Response: `UserProfileResponse`

2. request definition



```golang
type UserProfileRequest struct {
}
```


3. response definition



```golang
type UserProfileResponse struct {
	Profile UserProfile 
}

type UserProfile struct {
	Id uint64 
	Email string 
	Display_name string 
	Status int 
	Email_verified_at *int64 
	Created_at int64 
	Updated_at int64 
}
```

### 7. "Update user profile"

1. route definition

- Url: /api/v1/user/account/profile
- Method: PATCH
- Request: `UserUpdateProfileRequest`
- Response: `UserProfileResponse`

2. request definition



```golang
type UserUpdateProfileRequest struct {
	Display_name string `form:"display_name,optional" json:"display_name,optional"`
}
```


3. response definition



```golang
type UserProfileResponse struct {
	Profile UserProfile 
}

type UserProfile struct {
	Id uint64 
	Email string 
	Display_name string 
	Status int 
	Email_verified_at *int64 
	Created_at int64 
	Updated_at int64 
}
```

