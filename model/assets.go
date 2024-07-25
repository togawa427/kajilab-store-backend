package model

import "time"

type AssetGetResponse struct {
	Money 	int64 	`json:"money"`
	Debt 	int64 	`json:"debt"`
}

type AssetHistoryGetResponse struct {
	Date 		time.Time		`json:"date"`
	Money 	int64 			`json:"money"`
	Debt 		int64 			`json:"debt"`
	Product int64 			`json:"product"`
}

type AssetUpdateRequest struct {
	Money 	int64 	`json:"money"`
	Debt	int64 	`json:"debt"`
}