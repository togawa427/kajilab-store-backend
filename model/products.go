package model

type AllProductsGetResponse struct {
	Id     		int64	`json:"id"`
	Name   		string 	`json:"name"`
	Barcode   	int64 	`json:"barcode"`
	Price 		int64  	`json:"price"`
	Stock		int64	`json:"stock"`
	ImagePath	string	`json:"image_path"`
}
