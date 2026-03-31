package dto

type CreatePublishTaskReq struct {
	ProductID       int64  `json:"productId"`
	TargetPlatform  string `json:"targetPlatform"`
	CategoryID      string `json:"categoryId"`
	FreightTemplate string `json:"freightTemplate"`
}

type PublishTaskResp struct {
	ID                int64  `json:"id"`
	ProductID         int64  `json:"productId"`
	TargetPlatform    string `json:"targetPlatform"`
	PlatformProductID string `json:"platformProductId"`
	Status            string `json:"status"`
	ErrorMessage      string `json:"errorMessage"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
}

type PublishTaskFilter struct {
	Status         *string `json:"status,omitempty"`
	TargetPlatform *string `json:"targetPlatform,omitempty"`
	Page           int     `json:"page"`
	PageSize       int     `json:"pageSize"`
}
