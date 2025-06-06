package model

import "time"

type ProductsGetResponseProduct struct {
	Id        int64               `json:"id"`
	Name      string              `json:"name"`
	Barcode   int64               `json:"barcode"`
	Price     int64               `json:"price"`
	Stock     int64               `json:"stock"`
	TagId     int64               `json:"tag_id"`
	ImagePath string              `json:"image_path"`
	Tags      []TagGetResponseTag `json:"tags"`
}

type ProductsGetResponse struct {
	Products   []ProductsGetResponseProduct `json:"products"`
	TotalCount int64                        `json:"total_count"`
}

type ProductGetResponse struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Barcode   int64  `json:"barcode"`
	Price     int64  `json:"price"`
	Stock     int64  `json:"stock"`
	TagId     int64  `json:"tag_id"`
	ImagePath string `json:"image_path"`
}

type ProductStockLogsGetResponse struct {
	Logs []ProductStockLogJson `json:"logs"`
}

type ProductStockLogJson struct {
	Date  time.Time `json:"date"`
	Stock int64     `json:"stock"`
}

type BuyLogsGetResponse struct {
	Id          int64                `json:"id"`
	Price       int64                `json:"price"`
	PayAt       time.Time            `json:"pay_at"`
	PayAtStr    string               `json:"pay_at_str"`
	Method      string               `json:"method"`
	UserName    string               `json:"user_name"`
	Content     string               `json:"content"`
	PrevDebt    int64                `json:"prev_debt"`
	CurrentDebt int64                `json:"current_debt"`
	Products    []BuyProductResponse `json:"products"`
}

type BuyProductResponse struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Barcode   int64  `json:"barcode"`
	Quantity  int64  `json:"quantity"`
	UnitPrice int64  `json:"unit_price"`
}

type ArriveLogGetResponse struct {
	Id       int64               `json:"id"`
	Money    int64               `json:"money"`
	Value    int64               `json:"value"`
	ArriveAt time.Time           `json:"arrive_at"`
	Products []ArriveProductJson `json:"products"`
}

type ArriveProductJson struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Barcode  int64  `json:"barcode"`
	Quantity int64  `json:"quantity"`
	Value    int64  `json:"value"`
}

type ProductCreateRequestTag struct {
	Name string `json:"name"`
}

type ProductCreateRequest struct {
	Name    string                    `json:"name"`
	Barcode int64                     `json:"barcode"`
	Price   int64                     `json:"price"`
	TagId   int64                     `json:"tag_id"`
	Tags    []ProductCreateRequestTag `json:"tags"`
}

type ProductsBuyRequest struct {
	PayAt      time.Time        `json:"pay_at"`
	Method     string           `json:"method"`
	UserNumber string           `json:"user_number"`
	Products   []ProductBuyJson `json:"products"`
}

type ProductBuyJson struct {
	Id        int64 `json:"id"`
	Quantity  int64 `json:"quantity"`
	UnitPrice int64 `json:"unit_price"`
}

type ProductsArriveRequest struct {
	ArriveAt    time.Time           `json:"arrive_at"`
	Money       int64               `json:"money"`
	UserBarcode *string             `json:"user_barcode"`
	Products    []ProductArriveJson `json:"products"`
}

type ProductArriveJson struct {
	Id       int64 `json:"id"`
	Quantity int64 `json:"quantity"`
}

type ProductUpdateRequest struct {
	Id           int64                `json:"id"`
	Name         *string              `json:"name"`
	Barcode      *int64               `json:"barcode"`
	Price        *int64               `json:"price"`
	Stock        *int64               `json:"stock"`
	TagId        *int64               `json:"tag_id"`
	IsSold       *bool                `json:"is_sold"`
	WarningStock *int64               `json:"warning_stock"`
	SafetyStock  *int64               `json:"safety_stock"`
	Tags         *[]TagPostRequestTag `json:"tags"`
}

type ProductImageUpdateRequest struct {
	Id        int64  `json:"id"`
	ImagePath string `json:"image_path"`
}

type UserGetResponse struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Debt    int64  `json:"debt"`
	Barcode string `json:"barcode"`
}

type UserCreateRequest struct {
	Name    string `json:"name"`
	Barcode string `json:"barcode"`
}

type UserUpdateDebtRequest struct {
	Id   int64 `json:"id"`
	Debt int64 `json:"debt"`
}

type UserUpdateBarcodeRequest struct {
	Id      int64  `json:"id"`
	Barcode string `json:"barcode"`
}
