package dto

type ImportSourceItemReq struct {
	Platform  string `json:"platform"`
	SourceURL string `json:"sourceUrl"`
}

type SourceItemResp struct {
	ID          int64    `json:"id"`
	Platform    string   `json:"platform"`
	SourceURL   string   `json:"sourceUrl"`
	ExternalID  string   `json:"externalId"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
	PriceMin    float64  `json:"priceMin"`
	PriceMax    float64  `json:"priceMax"`
	Supplier    struct {
		Name   string  `json:"name"`
		Rating float64 `json:"rating"`
		Region string  `json:"region"`
	} `json:"supplier"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	SalesVolume int      `json:"salesVolume"`
	MinOrder    int      `json:"minOrder"`
	Status      string   `json:"status"`
	FetchedAt   string   `json:"fetchedAt"`
	CreatedAt   string   `json:"createdAt"`
}

type SourceItemFilter struct {
	Platform    *string  `json:"platform,omitempty"`
	Category    *string  `json:"category,omitempty"`
	PriceMin    *float64 `json:"priceMin,omitempty"`
	PriceMax    *float64 `json:"priceMax,omitempty"`
	SupplierMin *float64 `json:"supplierMin,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Status      *string  `json:"status,omitempty"`
	Keyword     *string  `json:"keyword,omitempty"`
	Page        int      `json:"page"`
	PageSize    int      `json:"pageSize"`
}

type UpdateSourceItemStatusReq struct {
	Status string `json:"status"`
}

type AddTagReq struct {
	Tag string `json:"tag"`
}
