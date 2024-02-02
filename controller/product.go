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

	// DBから購入物取得
	buyProducts, err := ProductService.GetBuyProducts()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get buyproducts from DB")
		return
	}

	resBuyProducts := []model.BuyProductResponse{}
	for _, buyProduct := range buyProducts {
		resBuyProducts = append(resBuyProducts, model.BuyProductResponse{
			Name: "じゃがりこ",
			Quantity: buyProduct.Quantity,
			UnitPrice: buyProduct.UnitPrice,
		})
	}

	resLogs := []model.BuyLogsGetResponse{}
	for _, log := range logs {
		resLogs = append(resLogs, model.BuyLogsGetResponse{
			Id: int64(log.ID),
			Price: log.Price,
			PayAt: log.PayAt,
			Method: log.Method,
			Products: resBuyProducts,
		})
	}

	c.JSON(http.StatusOK, resLogs)
}