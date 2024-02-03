package controller

import (

	// "kajilab-store-backend/service"

	"kajilab-store-backend/model"
	"kajilab-store-backend/service"
	"net/http"
	"strconv"

	// // "strconv"

	"github.com/gin-gonic/gin"
)

func GetAllProducts(c *gin.Context) {
	ProductService := service.ProductService{}

	// limitとoffsetの取得
	limit, err := strconv.ParseInt(c.Query("limit"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "limit is not number")
		return
	}
	offset, err := strconv.ParseInt(c.Query("offset"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "offset is not number")
		return
	}

	// DBから商品情報取得
	products, err := ProductService.GetAllProducts(limit, offset)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get products from DB")
		return
	}

	resProducts := []model.AllProductsGetResponse{}

	for _, product := range products {
		resProducts = append(resProducts, model.AllProductsGetResponse{
			Id: int64(product.ID),
			Name: product.Name,
			Barcode: product.Barcode,
			Price: product.Price,
			Stock: product.Stock,
			TagId: product.TagId,
			ImagePath: product.ImagePath,
		})
	}

	c.JSON(http.StatusOK, resProducts)
}

func GetBuyLogs(c *gin.Context) {
	ProductService := service.ProductService{}

	// limitの取得
	limit, err := strconv.ParseInt(c.Query("limit"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "limit is not number")
		return
	}

	// DBから購入ログ取得
	logs, err := ProductService.GetBuyLogs(limit)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get logs from DB")
		return
	}

	resLogs := []model.BuyLogsGetResponse{}
	for _, log := range logs {

		// 購入物の商品情報を取得
		buyProducts, err := ProductService.GetBuyProductsByPaymentId(int64(log.ID))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get buyproducts from DB")
			return
		}
		buyProductsJson := []model.BuyProductResponse{}
		for _, buyProduct := range buyProducts {

			// 商品名を取得
			productInfo, err := ProductService.GetProductById(int64(buyProduct.ID))
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get buyproduct from DB")
				return
			}

			buyProductsJson = append(buyProductsJson, model.BuyProductResponse{
				Name: productInfo.Name,
				Quantity: buyProduct.Quantity,
				UnitPrice: buyProduct.UnitPrice,
			})
		}

		resLogs = append(resLogs, model.BuyLogsGetResponse{
			Id: int64(log.ID),
			Price: log.Price,
			PayAt: log.PayAt,
			Method: log.Method,
			UserName: "",
			Products: buyProductsJson,
		})
	}

	c.JSON(http.StatusOK, resLogs)
}

func GetArriveLogs(c *gin.Context){
	ProductService := service.ProductService{}

	// limitの取得
	limit, err := strconv.ParseInt(c.Query("limit"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "limit is not number")
		return
	}

	// 最終的に返す値
	logsJson := []model.ArriveLogGetResponse{}

	// DBから入荷ログ取得
	logs, err := ProductService.GetArriveLogs(limit)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get logs from DB")
		return
	}
	for _, log := range logs {
		totalValue := int64(0)

		// 入荷物の商品情報を取得
		arriveProducts, err := ProductService.GetArriveProductsByArriveId(int64(log.ID))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get arrivalproducts from DB")
			return
		}
		arriveProductsJson := []model.ArriveProductJson{}
		for _, arriveProduct := range arriveProducts {

			// 商品名を取得
			productInfo, err := ProductService.GetProductById(int64(arriveProduct.ID))
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get buyproduct from DB")
				return
			}

			arriveProductsJson = append(arriveProductsJson, model.ArriveProductJson{
				Name: productInfo.Name,
				Quantity: arriveProduct.Quantity,
				Value: productInfo.Price,
			})
			totalValue = totalValue + (productInfo.Price*arriveProduct.Quantity)
		}

		logsJson = append(logsJson, model.ArriveLogGetResponse{
			Id: int64(log.ID),
			Money: log.Money,
			Value: totalValue,
			ArriveAt: log.ArriveAt,
			Products: arriveProductsJson,
		})
	}

	c.JSON(http.StatusOK, logsJson)
}