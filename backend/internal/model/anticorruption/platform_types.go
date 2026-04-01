package anticorruption

type SourceProduct struct {
	ExternalID  string   `json:"externalId"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
	PriceMin    float64  `json:"priceMin"`
	PriceMax    float64  `json:"priceMax"`
	Supplier    Supplier `json:"supplier"`
	Category    string   `json:"category"`
	SalesVolume int      `json:"salesVolume"`
	MinOrder    int      `json:"minOrder"`
}

type Supplier struct {
	Name   string  `json:"name"`
	Rating float64 `json:"rating"`
	Region string  `json:"region"`
}

type PublishResult struct {
	PlatformProductID string `json:"platformProductId"`
	Success           bool   `json:"success"`
	ErrorMessage      string `json:"errorMessage"`
}

type PublishConfig struct {
	CategoryID      string `json:"categoryId"`
	FreightTemplate string `json:"freightTemplate"`
}
