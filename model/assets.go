package model

type AssetGetResponse struct {
	Money 	int64 	`json:"money"`
	Debt 	int64 	`json:"debt"`
}

type AssetUpdateRequest struct {
	Money 	int64 	`json:"money"`
	Debt	int64 	`json:"debt"`
}