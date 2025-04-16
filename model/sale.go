package model

type GetSalesResponse struct {
	Year           int64                         `json:"year"`
	Month          int64                         `json:"month"`
	TotalMonthSale int64                         `json:"total_month_sale"`
	Sales          []GetSalesResponsePaymentsDay `json:"name"`
}

type GetSalesResponsePaymentsDay struct {
	Day       int64 `json:"day"`
	TotalSale int64 `json:"total_sale"`
	Payments  []GetSalesResponsePayment
}

type GetSalesResponsePayment struct {
	Name      string `json:"name"`
	Quantity  int64  `json:"quantity"`
	UnitPrice int64  `json:"unit_price"`
}
