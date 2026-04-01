package dto

type CreateProductFromSourceReq struct {
	SourceItemID int64 `json:"sourceItemId"`
}

type UpdateProductReq struct {
	Name        *string   `json:"name,omitempty"`
	Description *string   `json:"description,omitempty"`
	Images      []string  `json:"images,omitempty"`
	CostPrice   *float64  `json:"costPrice,omitempty"`
	SellPrice   *float64  `json:"sellPrice,omitempty"`
	CategoryID  *string   `json:"categoryId,omitempty"`
	SKUs        []SKUItem `json:"skus,omitempty"`
}

type SKUItem struct {
	ID        int64   `json:"id,omitempty"`
	SpecName  string  `json:"specName"`
	SpecValue string  `json:"specValue"`
	Price     float64 `json:"price"`
	Stock     int     `json:"stock"`
}

type ProductResp struct {
	ID           int64     `json:"id"`
	SourceItemID int64     `json:"sourceItemId"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Images       []string  `json:"images"`
	CostPrice    float64   `json:"costPrice"`
	SellPrice    float64   `json:"sellPrice"`
	CategoryID   string    `json:"categoryId"`
	Status       string    `json:"status"`
	SKUs         []SKUItem `json:"skus"`
	CreatedAt    string    `json:"createdAt"`
	UpdatedAt    string    `json:"updatedAt"`
}

type ProductFilter struct {
	Status   *string `json:"status,omitempty"`
	Keyword  *string `json:"keyword,omitempty"`
	Page     int     `json:"page"`
	PageSize int     `json:"pageSize"`
}
