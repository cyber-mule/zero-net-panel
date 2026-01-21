### 1. "Get site settings"

1. route definition

- Url: /api/v1/admin/site-settings
- Method: GET
- Request: `-`
- Response: `AdminSiteSettingResponse`

2. request definition



3. response definition



```golang
type AdminSiteSettingResponse struct {
	Setting SiteSetting 
}

type SiteSetting struct {
	Id uint64 
	Name string 
	Logo_url string 
	Access_domain string 
	Created_at int64 
	Updated_at int64 
}
```

### 2. "Update site settings"

1. route definition

- Url: /api/v1/admin/site-settings
- Method: PATCH
- Request: `AdminUpdateSiteSettingRequest`
- Response: `AdminSiteSettingResponse`

2. request definition



```golang
type AdminUpdateSiteSettingRequest struct {
	Name string `form:"name,optional" json:"name,optional"`
	Logo_url string `form:"logo_url,optional" json:"logo_url,optional"`
	Access_domain string `form:"access_domain,optional" json:"access_domain,optional"`
}
```


3. response definition



```golang
type AdminSiteSettingResponse struct {
	Setting SiteSetting 
}

type SiteSetting struct {
	Id uint64 
	Name string 
	Logo_url string 
	Access_domain string 
	Created_at int64 
	Updated_at int64 
}
```

