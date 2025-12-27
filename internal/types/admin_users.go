package types

// AdminListUsersRequest describes user listing filters.
type AdminListUsersRequest struct {
	Page    int    `form:"page,optional" json:"page,optional"`
	PerPage int    `form:"per_page,optional" json:"per_page,optional"`
	Query   string `form:"q,optional" json:"q,optional"`
	Status  string `form:"status,optional" json:"status,optional"`
	Role    string `form:"role,optional" json:"role,optional"`
}

// AdminUserSummary summarizes user account status for admin views.
type AdminUserSummary struct {
	ID                  uint64   `json:"id"`
	Email               string   `json:"email"`
	DisplayName         string   `json:"display_name"`
	Roles               []string `json:"roles"`
	Status              string   `json:"status"`
	EmailVerifiedAt     *int64   `json:"email_verified_at,omitempty"`
	FailedLoginAttempts int      `json:"failed_login_attempts"`
	LockedUntil         *int64   `json:"locked_until,omitempty"`
	LastLoginAt         *int64   `json:"last_login_at,omitempty"`
	CreatedAt           int64    `json:"created_at"`
	UpdatedAt           int64    `json:"updated_at"`
}

// AdminUserListResponse returns paginated users.
type AdminUserListResponse struct {
	Users      []AdminUserSummary `json:"users"`
	Pagination PaginationMeta     `json:"pagination"`
}

// AdminUserResponse returns a single user.
type AdminUserResponse struct {
	User AdminUserSummary `json:"user"`
}

// AdminCreateUserRequest provisions a user account.
type AdminCreateUserRequest struct {
	Email         string   `json:"email"`
	Password      string   `json:"password"`
	DisplayName   *string  `json:"display_name,omitempty,optional"`
	Roles         []string `json:"roles,omitempty,optional"`
	Status        *string  `json:"status,omitempty,optional"`
	EmailVerified *bool    `json:"email_verified,omitempty,optional"`
}

// AdminUpdateUserStatusRequest updates user status.
type AdminUpdateUserStatusRequest struct {
	UserID uint64 `path:"id"`
	Status string `json:"status"`
}

// AdminUpdateUserRolesRequest updates user roles.
type AdminUpdateUserRolesRequest struct {
	UserID uint64   `path:"id"`
	Roles  []string `json:"roles"`
}

// AdminResetUserPasswordRequest resets a user password.
type AdminResetUserPasswordRequest struct {
	UserID   uint64 `path:"id"`
	Password string `json:"password"`
}

// AdminForceLogoutRequest forces logout for a user.
type AdminForceLogoutRequest struct {
	UserID uint64 `path:"id"`
}

// AdminResetUserPasswordResponse acknowledges reset.
type AdminResetUserPasswordResponse struct {
	Message string `json:"message"`
}

// AdminForceLogoutResponse acknowledges forced logout.
type AdminForceLogoutResponse struct {
	Message string `json:"message"`
}

// AdminRotateUserCredentialRequest rotates a user's credential.
type AdminRotateUserCredentialRequest struct {
	UserID uint64 `path:"id"`
}

// AdminRotateUserCredentialResponse returns new credential metadata.
type AdminRotateUserCredentialResponse struct {
	UserID     uint64            `json:"user_id"`
	Credential CredentialSummary `json:"credential"`
}
