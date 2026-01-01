### 1. "List active announcements"

1. route definition

- Url: /api/v1/user/announcements
- Method: GET
- Request: `UserAnnouncementListRequest`
- Response: `UserAnnouncementListResponse`

2. request definition



```golang
type UserAnnouncementListRequest struct {
	Audience string `form:"audience,optional" json:"audience,optional"`
	Limit int `form:"limit,optional" json:"limit,optional"`
}
```


3. response definition



```golang
type UserAnnouncementListResponse struct {
	Announcements []UserAnnouncementSummary 
}
```

