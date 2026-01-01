### 1. "Get third-party security settings"

1. route definition

- Url: /api/v1/admin/security-settings
- Method: GET
- Request: `-`
- Response: `AdminSecuritySettingResponse`

2. request definition



3. response definition



```golang
type AdminSecuritySettingResponse struct {
	Setting SecuritySetting 
}

type SecuritySetting struct {
	Id uint64 
	Third_party_api_enabled bool 
	Api_key string 
	Api_secret string 
	Encryption_algorithm string 
	Nonce_ttl_seconds int 
	Created_at int64 
	Updated_at int64 
}
```

### 2. "Update third-party security settings"

1. route definition

- Url: /api/v1/admin/security-settings
- Method: PATCH
- Request: `AdminUpdateSecuritySettingRequest`
- Response: `AdminSecuritySettingResponse`

2. request definition



```golang
type AdminUpdateSecuritySettingRequest struct {
	Third_party_api_enabled bool `form:"third_party_api_enabled,optional" json:"third_party_api_enabled,optional"`
	Api_key string `form:"api_key,optional" json:"api_key,optional"`
	Api_secret string `form:"api_secret,optional" json:"api_secret,optional"`
	Encryption_algorithm string `form:"encryption_algorithm,optional" json:"encryption_algorithm,optional"`
	Nonce_ttl_seconds int `form:"nonce_ttl_seconds,optional" json:"nonce_ttl_seconds,optional"`
}
```


3. response definition



```golang
type AdminSecuritySettingResponse struct {
	Setting SecuritySetting 
}

type SecuritySetting struct {
	Id uint64 
	Third_party_api_enabled bool 
	Api_key string 
	Api_secret string 
	Encryption_algorithm string 
	Nonce_ttl_seconds int 
	Created_at int64 
	Updated_at int64 
}
```

