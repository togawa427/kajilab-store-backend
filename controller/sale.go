package controller

import (
	"net/http"
	"strconv"
	"time"

	"kajilab-store-backend/model"
	"kajilab-store-backend/service"

	"github.com/gin-gonic/gin"
)

func findPaymentByDay(payments []model.GetSalesResponsePaymentsDay, day int64) *int {
	for i := range payments {
		if payments[i].Day == day {
			return &i
		}
	}
	return nil
}

func findPaymentProduct(paymentProducts []model.GetSalesResponsePayment, name string) *int {
	for i := range paymentProducts {
		if paymentProducts[i].Name == name {
			return &i
		}
	}
	return nil
}

func GetMonthSales(c *gin.Context) {
	ProductService := service.ProductService{}

	// クエリパラメータの取得
	// limit
	limit, err := strconv.ParseInt(c.Query("limit"), 10, 64)
	if err != nil {
		limit = 0
	}
	// year
	year, err := strconv.ParseInt(c.Query("year"), 10, 64)
	if err != nil {
		year = int64(time.Now().Year())
	}
	// month
	month, err := strconv.ParseInt(c.Query("month"), 10, 64)
	if err != nil {
		month = int64(time.Now().Month())
	}

	// DBから購入ログ取得
	logs, err := ProductService.GetBuyLogs(limit, year, month)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get logs from DB")
		return
	}

	resPaymentsPerDay := []model.GetSalesResponsePaymentsDay{}
	salesMonth := int64(0)

	for _, log := range logs {

		// 購入物の商品情報を取得
		reqBuyProducts, err := ProductService.GetProductLogsBySourceId(int64(log.ID))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get buyproducts from DB")
			return
		}

		// 購入物の支払い日時を取得
		indexPaymentByDay := findPaymentByDay(resPaymentsPerDay, int64(log.PayAt.Day()))
		if indexPaymentByDay != nil {
			// 既にその日が存在する場合
			for _, reqBuyProduct := range reqBuyProducts {
				// 商品名を取得
				reqBuyProductInfo, err := ProductService.GetProductById(int64(reqBuyProduct.ProductId))
				if err != nil {
					c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get buyproduct from DB")
					return
				}
				indexPaymentProduct := findPaymentProduct(resPaymentsPerDay[*indexPaymentByDay].Payments, reqBuyProductInfo.Name)
				if indexPaymentProduct != nil {
					// 既にその商品情報がレスポンスに追加済みである場合、数を更新
					resPaymentsPerDay[*indexPaymentByDay].Payments[*indexPaymentProduct].Quantity += reqBuyProduct.Quantity
				} else {
					// その商品がレスポンスに未登録である場合、商品を追加
					resPaymentsPerDay[*indexPaymentByDay].Payments = append(resPaymentsPerDay[*indexPaymentByDay].Payments, model.GetSalesResponsePayment{
						Name:      reqBuyProductInfo.Name,
						Quantity:  reqBuyProduct.Quantity,
						UnitPrice: reqBuyProduct.UnitPrice,
					})
				}
				// その日、その月の売り上げ金額を更新
				resPaymentsPerDay[*indexPaymentByDay].TotalSale += reqBuyProduct.UnitPrice * reqBuyProduct.Quantity
				salesMonth += reqBuyProduct.UnitPrice * reqBuyProduct.Quantity
			}
		} else {
			// 新しい日付の場合
			tmpProducts := []model.GetSalesResponsePayment{}
			salesDay := int64(0)
			for _, reqBuyProduct := range reqBuyProducts {
				// 商品名を取得
				reqBuyProductInfo, err := ProductService.GetProductById(int64(reqBuyProduct.ProductId))
				if err != nil {
					c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get buyproduct from DB")
					return
				}
				// 追加する
				tmpProducts = append(tmpProducts, model.GetSalesResponsePayment{
					Name:      reqBuyProductInfo.Name,
					Quantity:  reqBuyProduct.Quantity,
					UnitPrice: reqBuyProduct.UnitPrice,
				})
				// その日、その月の売り上げ金額を更新
				salesDay += reqBuyProduct.UnitPrice * reqBuyProduct.Quantity
				salesMonth += reqBuyProduct.UnitPrice * reqBuyProduct.Quantity
			}
			resPaymentsPerDay = append(resPaymentsPerDay, model.GetSalesResponsePaymentsDay{
				Day:       int64(log.PayAt.Day()),
				TotalSale: salesDay,
				Payments:  tmpProducts,
			})
		}
	}

	paymentPerMonth := model.GetSalesResponse{
		Year:           year,
		Month:          month,
		TotalMonthSale: salesMonth,
		ResponseDate:   time.Now().Format(time.DateTime),
		Sales:          resPaymentsPerDay,
	}

	c.JSON(http.StatusOK, paymentPerMonth)
}
