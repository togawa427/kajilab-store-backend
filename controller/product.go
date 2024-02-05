package controller

import (

	// "kajilab-store-backend/service"

	"fmt"
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

// 購入ログの取得
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
			productInfo, err := ProductService.GetProductById(int64(buyProduct.ProductId))
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get buyproduct from DB")
				return
			}

			buyProductsJson = append(buyProductsJson, model.BuyProductResponse{
				Id: int64(productInfo.ID),
				Name: productInfo.Name,
				Barcode: productInfo.Barcode,
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

// 入荷ログ取得API
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
			productInfo, err := ProductService.GetProductById(int64(arriveProduct.ProductId))
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get buyproduct from DB")
				return
			}

			arriveProductsJson = append(arriveProductsJson, model.ArriveProductJson{
				Id: int64(productInfo.ID,),
				Name: productInfo.Name,
				Barcode: int64(productInfo.Barcode),
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

// 商品登録API
func CreateProduct(c *gin.Context){
	ProductService := service.ProductService{}
	ProductCreateRequest := model.ProductCreateRequest{}
	err := c.Bind(&ProductCreateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "request is not correct")
		return
	}

	// リクエストの商品情報をデータベースの型へ変換
	product := model.Product{
		Name: ProductCreateRequest.Name,
		Barcode: ProductCreateRequest.Barcode,
		Price: ProductCreateRequest.Price,
		Stock: 0,
		TagId: ProductCreateRequest.TagId,
		ImagePath: strconv.Itoa(int(ProductCreateRequest.Barcode))+ ".jpg",
	}

	err = ProductService.CreateProduct(&product)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal create product")
		return
	}

	c.JSON(http.StatusOK, "success")

}

// 商品購入時API
func BuyProducts(c *gin.Context) {
	ProductService := service.ProductService{}
	AssetService := service.AssetService{}
	ProductsBuyRequest := model.ProductsBuyRequest{}
	err := c.Bind(&ProductsBuyRequest)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "request is not correct")
		return
	}

	// 合計金額を算出
	totalPrice := int64(0)
	productsJson := ProductsBuyRequest.Products
	for _, productJson := range productsJson {
		totalPrice += productJson.UnitPrice * productJson.Quantity
	}

	// ユーザ番号からユーザIDを取得
	// 未実装

	// 購入情報を登録
	// リクエストの商品情報をデータベースの型へ変換
	payment := model.Payment{
		Price: totalPrice,
		PayAt: ProductsBuyRequest.PayAt,
		Method: ProductsBuyRequest.Method,
		UserId: 0,
	}
	// DBへ保存
	paymentId, err := ProductService.CreatePayment(&payment)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal create payment")
		return
	}

	fmt.Println(paymentId)


	// 購入商品情報を登録
	for _, productJson := range productsJson {
		fmt.Println("start")
		// リクエストの商品情報をデータベースの型へ変換
		product := model.PaymentProduct{
			PaymentId: paymentId,
			ProductId: productJson.Id,
			Quantity: productJson.Quantity,
			UnitPrice: productJson.UnitPrice,
		}
		// DBへ保存
		err = ProductService.CreatePaymentProduct(&product)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, "fetal create payment_product")
			return
		}	
	}

	// お金を減らす
	err = AssetService.IncreaseMoney(totalPrice)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal increase money")
		return
	}

	c.JSON(http.StatusOK, "success")

}

// 商品入荷API
func ArriveProducts(c *gin.Context) {
	ProductService := service.ProductService{}
	AssetService := service.AssetService{}
	ProductsArriveRequest := model.ProductsArriveRequest{}
	err := c.Bind(&ProductsArriveRequest)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "request is not correct")
		return
	}

	// 購入情報を登録
	// リクエストの商品情報をデータベースの型へ変換
	arrival := model.Arrival{
		Money: ProductsArriveRequest.Money,
		ArriveAt: ProductsArriveRequest.ArriveAt,
	}
	// DBへ保存
	arrivalId, err := ProductService.CreateArrival(&arrival)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal create arrival")
		return
	}

	// 入荷商品情報を登録
	productsJson := ProductsArriveRequest.Products
	for _, productJson := range productsJson {
		// リクエストの商品情報をデータベースの型へ変換
		product := model.ArrivalProduct{
			ArrivalId: arrivalId,
			ProductId: productJson.Id,
			Quantity: productJson.Quantity,
		}
		// DBへ保存
		err = ProductService.CreateArriveProduct(&product)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, "fetal create arrival_product")
			return
		}	
	}

	// お金を減らす
	err = AssetService.IncreaseMoney(0-ProductsArriveRequest.Money)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal decrease money")
		return
	}

	c.JSON(http.StatusOK, "success")
}