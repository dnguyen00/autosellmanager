package Structs //this is the struct where i parse collections against, collections.go is the struct i use to store

type TradeCollections struct {
	PreviousPageCursor string `json:"previousPageCursor"`
	NextPageCursor     string `json:"nextPageCursor"`
	Data               []struct {
		UserAssetID                int64       `json:"userAssetId"`
		SerialNumber               int         `json:"serialNumber"`
		AssetID                    int64       `json:"assetId"`
		Name                       string      `json:"name"`
		RecentAveragePrice         int         `json:"recentAveragePrice"`
		OriginalPrice              interface{} `json:"originalPrice"`
		AssetStock                 int         `json:"assetStock"`
		BuildersClubMembershipType int         `json:"buildersClubMembershipType"`
	} `json:"data"`
}
