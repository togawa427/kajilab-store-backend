package model

import "time"

type AllProductsGetResponse struct {
	Id     		int64	`json:"id"`
	Name   		string 	`json:"name"`
	Barcode   	int64 	`json:"barcode"`
	Price 		int64  	`json:"price"`
	Stock		int64	`json:"stock"`
	TagId 		int64 	`json:"tag_id"`
	ImagePath	string	`json:"image_path"`
}

type BuyLogsGetResponse struct {
	Id     		int64					`json:"id"`
	Price 		int64  					`json:"price"`
	PayAt 		time.Time				`json:"pay_at"`
	Method 		string					`json:"method"`
	Products 	[]BuyProductResponse	`json:"products"`
}

type BuyProductResponse struct {
	Name		string	`json:"name"`
	Quantity 	int64	`json:"quantity"`
	UnitPrice	int64	`json:"unit_price"`
}
