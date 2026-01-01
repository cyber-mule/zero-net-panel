### 1. "List subscription templates"

1. route definition

- Url: /api/v1/admin/subscription-templates
- Method: GET
- Request: `AdminListSubscriptionTemplatesRequest`
- Response: `AdminSubscriptionTemplateListResponse`

2. request definition



```golang
type AdminListSubscriptionTemplatesRequest struct {
	Page int `form:"page,optional" json:"page,optional"`
	Per_page int `form:"per_page,optional" json:"per_page,optional"`
	Sort string `form:"sort,optional" json:"sort,optional"`
	Direction string `form:"direction,optional" json:"direction,optional"`
	Q string `form:"q,optional" json:"q,optional"`
	Client_type string `form:"client_type,optional" json:"client_type,optional"`
	Format string `form:"format,optional" json:"format,optional"`
	Include_drafts bool `form:"include_drafts,optional" json:"include_drafts,optional"`
}
```


3. response definition



```golang
type AdminSubscriptionTemplateListResponse struct {
	Templates []SubscriptionTemplateSummary 
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

### 2. "Create subscription template"

1. route definition

- Url: /api/v1/admin/subscription-templates
- Method: POST
- Request: `AdminCreateSubscriptionTemplateRequest`
- Response: `SubscriptionTemplateSummary`

2. request definition



```golang
type AdminCreateSubscriptionTemplateRequest struct {
	Name string 
	Description string `form:"description,optional" json:"description,optional"`
	Client_type string 
	Format string `form:"format,optional" json:"format,optional"`
	Content string 
	Variables map[string]TemplateVariable `form:"variables,optional" json:"variables,optional"`
	Is_default bool `form:"is_default,optional" json:"is_default,optional"`
}
```


3. response definition



```golang
type SubscriptionTemplateSummary struct {
	Id uint64 
	Name string 
	Description string 
	Client_type string 
	Format string 
	Content string `form:"content,optional" json:"content,optional"`
	Variables map[string]TemplateVariable 
	Is_default bool 
	Version uint32 
	Updated_at int64 
	Published_at int64 
	Last_published_by string 
}
```

### 3. "Update subscription template"

1. route definition

- Url: /api/v1/admin/subscription-templates/:id
- Method: PATCH
- Request: `AdminUpdateSubscriptionTemplateRequest`
- Response: `SubscriptionTemplateSummary`

2. request definition



```golang
type AdminUpdateSubscriptionTemplateRequest struct {
	Id uint64 
	Name string `form:"name,optional" json:"name,optional"`
	Description string `form:"description,optional" json:"description,optional"`
	Format string `form:"format,optional" json:"format,optional"`
	Content string `form:"content,optional" json:"content,optional"`
	Variables map[string]TemplateVariable `form:"variables,optional" json:"variables,optional"`
	Is_default bool `form:"is_default,optional" json:"is_default,optional"`
}
```


3. response definition



```golang
type SubscriptionTemplateSummary struct {
	Id uint64 
	Name string 
	Description string 
	Client_type string 
	Format string 
	Content string `form:"content,optional" json:"content,optional"`
	Variables map[string]TemplateVariable 
	Is_default bool 
	Version uint32 
	Updated_at int64 
	Published_at int64 
	Last_published_by string 
}
```

### 4. "List template publish history"

1. route definition

- Url: /api/v1/admin/subscription-templates/:id/history
- Method: GET
- Request: `AdminSubscriptionTemplateHistoryRequest`
- Response: `AdminSubscriptionTemplateHistoryResponse`

2. request definition



```golang
type AdminSubscriptionTemplateHistoryRequest struct {
	Id uint64 
}
```


3. response definition



```golang
type AdminSubscriptionTemplateHistoryResponse struct {
	Template_id uint64 
	History []SubscriptionTemplateHistoryEntry 
}
```

### 5. "Publish subscription template"

1. route definition

- Url: /api/v1/admin/subscription-templates/:id/publish
- Method: POST
- Request: `AdminPublishSubscriptionTemplateRequest`
- Response: `AdminPublishSubscriptionTemplateResponse`

2. request definition



```golang
type AdminPublishSubscriptionTemplateRequest struct {
	Id uint64 
	Changelog string `form:"changelog,optional" json:"changelog,optional"`
	Operator string `form:"operator,optional" json:"operator,optional"`
}
```


3. response definition



```golang
type AdminPublishSubscriptionTemplateResponse struct {
	Template SubscriptionTemplateSummary 
	History SubscriptionTemplateHistoryEntry 
}

type SubscriptionTemplateSummary struct {
	Id uint64 
	Name string 
	Description string 
	Client_type string 
	Format string 
	Content string `form:"content,optional" json:"content,optional"`
	Variables map[string]TemplateVariable 
	Is_default bool 
	Version uint32 
	Updated_at int64 
	Published_at int64 
	Last_published_by string 
}

type SubscriptionTemplateHistoryEntry struct {
	Version uint32 
	Changelog string 
	Published_at int64 
	Published_by string 
	Variables map[string]TemplateVariable 
}
```

