### 1. "List users with filters"

1. route definition

- Url: /api/v1/admin/users
- Method: GET
- Request: `AdminListUsersRequest`
- Response: `AdminUserListResponse`

2. request definition



```golang
type AdminListUsersRequest struct {
	Page int `form:"page,optional" json:"page,optional"`
	Per_page int `form:"per_page,optional" json:"per_page,optional"`
	Q string `form:"q,optional" json:"q,optional"`
	Status string `form:"status,optional" json:"status,optional"`
	Role string `form:"role,optional" json:"role,optional"`
}
```


3. response definition



```golang
type AdminUserListResponse struct {
	Users []AdminUserSummary 
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

### 2. "Create a new user"

1. route definition

- Url: /api/v1/admin/users
- Method: POST
- Request: `AdminCreateUserRequest`
- Response: `AdminUserResponse`

2. request definition



```golang
type AdminCreateUserRequest struct {
	Email string 
	Password string 
	Display_name string `form:"display_name,optional" json:"display_name,optional"`
	Roles []string `form:"roles,optional" json:"roles,optional"`
	Status string `form:"status,optional" json:"status,optional"`
	Email_verified bool `form:"email_verified,optional" json:"email_verified,optional"`
}
```


3. response definition



```golang
type AdminUserResponse struct {
	User AdminUserSummary 
}

type AdminUserSummary struct {
	Id uint64 
	Email string 
	Display_name string 
	Roles []string 
	Status string 
	Email_verified_at int64 `form:"email_verified_at,optional" json:"email_verified_at,optional"`
	Failed_login_attempts int 
	Locked_until int64 `form:"locked_until,optional" json:"locked_until,optional"`
	Last_login_at int64 `form:"last_login_at,optional" json:"last_login_at,optional"`
	Created_at int64 
	Updated_at int64 
}
```

### 3. "Rotate user credential"

1. route definition

- Url: /api/v1/admin/users/:id/credentials/rotate
- Method: POST
- Request: `AdminRotateUserCredentialRequest`
- Response: `AdminRotateUserCredentialResponse`

2. request definition



```golang
type AdminRotateUserCredentialRequest struct {
	Id uint64 `path:"id"`
}
```


3. response definition



```golang
type AdminRotateUserCredentialResponse struct {
	User_id uint64 
	Credential CredentialSummary 
}

type CredentialSummary struct {
	Version int 
	Status string 
	Issued_at int64 
	Deprecated_at *int64 
	Revoked_at *int64 
	Last_seen_at *int64 
}
```

### 4. "Force user logout"

1. route definition

- Url: /api/v1/admin/users/:id/force-logout
- Method: POST
- Request: `AdminForceLogoutRequest`
- Response: `AdminForceLogoutResponse`

2. request definition



```golang
type AdminForceLogoutRequest struct {
	Id uint64 `path:"id"`
}
```


3. response definition



```golang
type AdminForceLogoutResponse struct {
	Message string 
}
```

### 5. "Reset user password"

1. route definition

- Url: /api/v1/admin/users/:id/reset-password
- Method: POST
- Request: `AdminResetUserPasswordRequest`
- Response: `AdminResetUserPasswordResponse`

2. request definition



```golang
type AdminResetUserPasswordRequest struct {
	Id uint64 `path:"id"`
	Password string 
}
```


3. response definition



```golang
type AdminResetUserPasswordResponse struct {
	Message string 
}
```

### 6. "Update user roles"

1. route definition

- Url: /api/v1/admin/users/:id/roles
- Method: PATCH
- Request: `AdminUpdateUserRolesRequest`
- Response: `AdminUserResponse`

2. request definition



```golang
type AdminUpdateUserRolesRequest struct {
	Id uint64 `path:"id"`
	Roles []string 
}
```


3. response definition



```golang
type AdminUserResponse struct {
	User AdminUserSummary 
}

type AdminUserSummary struct {
	Id uint64 
	Email string 
	Display_name string 
	Roles []string 
	Status string 
	Email_verified_at int64 `form:"email_verified_at,optional" json:"email_verified_at,optional"`
	Failed_login_attempts int 
	Locked_until int64 `form:"locked_until,optional" json:"locked_until,optional"`
	Last_login_at int64 `form:"last_login_at,optional" json:"last_login_at,optional"`
	Created_at int64 
	Updated_at int64 
}
```

### 7. "Update user status (active/disabled)"

1. route definition

- Url: /api/v1/admin/users/:id/status
- Method: PATCH
- Request: `AdminUpdateUserStatusRequest`
- Response: `AdminUserResponse`

2. request definition



```golang
type AdminUpdateUserStatusRequest struct {
	Id uint64 `path:"id"`
	Status string 
}
```


3. response definition



```golang
type AdminUserResponse struct {
	User AdminUserSummary 
}

type AdminUserSummary struct {
	Id uint64 
	Email string 
	Display_name string 
	Roles []string 
	Status string 
	Email_verified_at int64 `form:"email_verified_at,optional" json:"email_verified_at,optional"`
	Failed_login_attempts int 
	Locked_until int64 `form:"locked_until,optional" json:"locked_until,optional"`
	Last_login_at int64 `form:"last_login_at,optional" json:"last_login_at,optional"`
	Created_at int64 
	Updated_at int64 
}
```

