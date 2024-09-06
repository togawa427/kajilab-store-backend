package model

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name   		string
	Barcode   	int64
	Price 		int64
	Stock		int64
	TagId		int64
	IsSold	bool `gorm:"default:true"`
	WarningStock int64 `gorm:"default:0"`
	SafetyStock int64 `gorm:"default:1"`
	ImagePath	string
}

type Asset struct {
	gorm.Model
	Money		int64
	Debt		int64
}

type Payment struct {
	gorm.Model
	Price		int64
	PayAt 		time.Time
	Method		string
	UserId 		int64
}

type PaymentProduct struct {
	gorm.Model
	PaymentId 	int64
	ProductId 	int64
	Quantity 	int64
	UnitPrice 	int64
}

type Arrival struct {
	gorm.Model
	Money 		int64
	ArriveAt 	time.Time
}

type ArrivalProduct struct {
	gorm.Model
	ArrivalId 	int64
	ProductId 	int64
	Quantity 	int64
}

type ProductLog struct {
	gorm.Model
	SourceId 	int64
	ProductId int64
	Quantity 	int64
	UnitPrice int64
	Stock 		int64
}

type User struct {
	gorm.Model
	Name 		string
	Debt 		int64
	Barcode string
}