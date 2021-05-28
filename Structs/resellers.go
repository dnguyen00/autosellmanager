package Structs

type Resellers struct {
	PreviousPageCursor string `json:"previousPageCursor"`
	NextPageCursor     string `json:"nextPageCursor"`
	Data               []struct {
		UserAssetID int64 `json:"userAssetId"`
		Seller      struct {
			ID   int64  `json:"id"`
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"seller"`
		Price        int `json:"price"`
		SerialNumber int `json:"serialNumber"`
	} `json:"data"`
}
