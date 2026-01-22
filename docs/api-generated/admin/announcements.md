### 1. "List announcements"

1. route definition

- Url: /api/v1/admin/announcements
- Method: GET
- Request: `AdminListAnnouncementsRequest`
- Response: `AdminAnnouncementListResponse`

2. request definition



```golang
type AdminListAnnouncementsRequest struct {
	Page int `form:"page,optional" json:"page,optional"`
	Per_page int `form:"per_page,optional" json:"per_page,optional"`
	Status int `form:"status,optional" json:"status,optional"`
	Category string `form:"category,optional" json:"category,optional"`
	Audience string `form:"audience,optional" json:"audience,optional"`
	Q string `form:"q,optional" json:"q,optional"`
	Sort string `form:"sort,optional" json:"sort,optional"`
	Direction string `form:"direction,optional" json:"direction,optional"`
}
```


3. response definition



```golang
type AdminAnnouncementListResponse struct {
	Announcements []AnnouncementSummary 
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

### 2. "Create announcement"

1. route definition

- Url: /api/v1/admin/announcements
- Method: POST
- Request: `AdminCreateAnnouncementRequest`
- Response: `AnnouncementSummary`

2. request definition



```golang
type AdminCreateAnnouncementRequest struct {
	Title string 
	Content string 
	Category string `form:"category,optional" json:"category,optional"`
	Audience string `form:"audience,optional" json:"audience,optional"`
	Is_pinned bool `form:"is_pinned,optional" json:"is_pinned,optional"`
	Priority int `form:"priority,optional" json:"priority,optional"`
	Created_by string `form:"created_by,optional" json:"created_by,optional"`
}
```


3. response definition



```golang
type AnnouncementSummary struct {
	Id uint64 
	Title string 
	Content string 
	Category string 
	Status int 
	Audience string 
	Is_pinned bool 
	Priority int 
	Visible_from int64 
	Visible_to int64 `form:"visible_to,optional" json:"visible_to,optional"`
	Published_at int64 `form:"published_at,optional" json:"published_at,optional"`
	Published_by string 
	Created_by string 
	Updated_by string 
	Created_at int64 
	Updated_at int64 
}
```

### 3. "Publish announcement"

1. route definition

- Url: /api/v1/admin/announcements/:id/publish
- Method: POST
- Request: `AdminPublishAnnouncementRequest`
- Response: `AnnouncementSummary`

2. request definition



```golang
type AdminPublishAnnouncementRequest struct {
	Id uint64 
	Visible_to int64 `form:"visible_to,optional" json:"visible_to,optional"`
	Operator string `form:"operator,optional" json:"operator,optional"`
}
```


3. response definition



```golang
type AnnouncementSummary struct {
	Id uint64 
	Title string 
	Content string 
	Category string 
	Status int 
	Audience string 
	Is_pinned bool 
	Priority int 
	Visible_from int64 
	Visible_to int64 `form:"visible_to,optional" json:"visible_to,optional"`
	Published_at int64 `form:"published_at,optional" json:"published_at,optional"`
	Published_by string 
	Created_by string 
	Updated_by string 
	Created_at int64 
	Updated_at int64 
}
```

