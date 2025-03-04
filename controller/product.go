package controller

import (

	// "kajilab-store-backend/service"

	"fmt"
	"kajilab-store-backend/model"
	"kajilab-store-backend/service"
	"net/http"
	"strconv"
	"time"

	// // "strconv"

	"github.com/gin-gonic/gin"
)

func GetAllProducts(c *gin.Context) {
	ProductService := service.ProductService{}

	// limitとoffsetの取得
	limit, err := strconv.ParseInt(c.Query("limit"), 10, 64)
	if err != nil {
		limit = 100
	}
	offset, err := strconv.ParseInt(c.Query("offset"), 10, 64)
	if err != nil {
		offset = 0
	}
	updateDays, err := strconv.ParseInt(c.Query("updated_days"), 10, 64)
	if err != nil {
		updateDays = 0
	}

	// DBから商品情報取得
	products, err := ProductService.GetAllProducts(limit, offset, updateDays)
	if err != nil {
		fmt.Println(err)
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

// バーコードから商品取得API
func GetProductByBarcode(c *gin.Context) {
	ProductService := service.ProductService{}
	barcode, err := strconv.ParseInt(c.Param("barcode"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "barcode is not number")
		return
	}
	// 商品情報を取得
	product, err := ProductService.GetProductByBarcode(barcode)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get product by barcode")
		return
	}
	// レスポンスの型に変換
	productJson := model.ProductGetResponse{
		Id: int64(product.ID),
		Name: product.Name,
		Barcode: product.Barcode,
		Price: product.Price,
		Stock: product.Stock,
		TagId: product.TagId,
		ImagePath: product.ImagePath,
	}

	c.JSON(http.StatusOK, productJson)
}

// 現在からxx日分の在庫
func GetProductStockLogsById(c *gin.Context) {
	ProductService := service.ProductService{}
	// 商品ID取得
	productId, err := strconv.ParseInt(c.Param("productId"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "productId is not number")
		return
	}

	// 何日分なのかを取得
	day, err := strconv.ParseInt(c.Query("day"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "day is not number")
		return
	}

	// day日分の商品ログを取得
	productLogs, err := ProductService.GetProductLogsByDay(day, productId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal to is get product logs by days")
		return
	}
	fmt.Println("GetProductStockLogsByIdの結果")
	fmt.Println(productLogs)
	var resLogs [20]model.ProductStockLogJson
	for i := range resLogs {
		resLogs[i].Date = time.Now()
		resLogs[i].Stock = 3
	}
	fmt.Println(resLogs)
}


// 購入ログの取得
func GetBuyLogs(c *gin.Context) {
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
		year = 0
	}
	// month
	month, err := strconv.ParseInt(c.Query("month"), 10, 64)
	if err != nil {
		month = 0
	}

	// DBから購入ログ取得
	logs, err := ProductService.GetBuyLogs(limit, year, month)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get logs from DB")
		return
	}

	resLogs := []model.BuyLogsGetResponse{}
	for _, log := range logs {

		// 購入物の商品情報を取得
		buyProducts, err := ProductService.GetProductLogsBySourceId(int64(log.ID))
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
			PayAtStr: log.PayAt.Format("2006/01/02 15:04:05.000"),
			Method: log.Method,
			UserName: "",
			Products: buyProductsJson,
		})
	}

	c.JSON(http.StatusOK, resLogs)
}

// ユーザIDから購入ログの取得
func GetBuyLogsByUser(c *gin.Context) {
	ProductService := service.ProductService{}

	// ユーザID取得
	userId, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "productId is not number")
		return
	}

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

	// DBから購入ログ取得
	logs, err := ProductService.GetBuyLogsByUserId(offset, limit, userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get logs from DB")
		return
	}

	resLogs := []model.BuyLogsGetResponse{}
	for _, log := range logs {

		// 購入物の商品情報を取得
		buyProducts, err := ProductService.GetProductLogsBySourceId(int64(log.ID))
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
	// 

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
		arriveProducts, err := ProductService.GetProductLogsBySourceId(int64(log.ID))
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
	UserService := service.UserService{}
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

	// リクエストの商品情報をデータベースの型へ変換
	payment := model.Payment{
		Price: totalPrice,
		PayAt: ProductsBuyRequest.PayAt,
		Method: ProductsBuyRequest.Method,
		UserId: 0,
	}

	user := model.User{}
	// ユーザバーコードがリクエストに含まれる場合ユーザIDを取得
	if(payment.Method == "card" && ProductsBuyRequest.UserNumber != ""){
		user, err = UserService.GetUserByBarcode(ProductsBuyRequest.UserNumber)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get user by barcode")
			return
		}
		payment.UserId = int64(user.ID)

		// 残高が足りているかチェック
		err = UserService.IsEnoughUserDebt(int64(user.ID), totalPrice)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusPaymentRequired, "Not enough user debt")
			return
		}
	}

	// 購入情報を登録
	// DBへ保存
	paymentId, err := ProductService.CreatePayment(&payment)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal create payment")
		return
	}

	// 商品ログを登録，在庫を更新
	for _, productJson := range productsJson {
		fmt.Println("start")

		// 現在の商品情報（在庫数取得のため）を取得
		productInfo, err := ProductService.GetProductById(productJson.Id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get product_stock")
			return
		}

		productLog := model.ProductLog{
			SourceId: paymentId,
			ProductId: productJson.Id,
			Quantity: productJson.Quantity,
			UnitPrice: productJson.UnitPrice,
			Stock: productInfo.Stock - productJson.Quantity,
		}

		// 商品ログをDBへ保存
		err = ProductService.CreateProductLog(&productLog)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, "fetal create payment_product")
			return
		}
		// 在庫情報をDBへ保存
		err = ProductService.IncreaseStock(productLog.ProductId, 0-productLog.Quantity)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, "fetal decrease stock")
			return
		}
	}

	// お金を減らす
	if(payment.Method == "card" && ProductsBuyRequest.UserNumber != ""){
		// カード支払いの時
		// ユーザのDebtを減らす
		err = UserService.IncreaseUserDebt(payment.UserId, 0-totalPrice)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, "fetal decrease user debt")
			return
		}
		// 商店のDebtを減らす
		err = AssetService.IncreaseDebt(0-totalPrice)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, "fetal decrease debt")
			return
		}
	} else {
		// 現金支払いの時は商店残高を増やす
		err = AssetService.IncreaseMoney(totalPrice)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, "fetal increase money")
			return
		}
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

	// 入荷商品情報を登録, 在庫情報を更新
	productsJson := ProductsArriveRequest.Products
	for _, productJson := range productsJson {

		// 現在の商品情報（在庫数取得のため）を取得
		productInfo, err := ProductService.GetProductById(productJson.Id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get product_stock")
			return
		}
		// リクエストの商品情報をデータベースの型へ変換
		product := model.ProductLog{
			SourceId: arrivalId,
			ProductId: productJson.Id,
			Quantity: productJson.Quantity,
			UnitPrice: productInfo.Price,
			Stock: productInfo.Stock + productJson.Quantity,
		}
		// 入荷商品情報をDBへ保存DBへ保存
		err = ProductService.CreateProductLog(&product)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, "fetal create arrival_product")
			return
		}
		// 在庫情報をDBへ保存
		err = ProductService.IncreaseStock(product.ProductId, product.Quantity)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, "fetal increase stock")
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

// 商品情報更新API
func UpdateProduct(c *gin.Context) {
	ProductService := service.ProductService{}
	ProductUpdateRequest := model.ProductUpdateRequest{}
	err := c.Bind(&ProductUpdateRequest)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "request is not correct")
		return
	}

	// 現在の商品情報を取得
	product, err := ProductService.GetProductById(ProductUpdateRequest.Id)
	if (err != nil) {
		c.AbortWithStatusJSON(http.StatusBadRequest, "product not found")
    return
	}

	// 値が存在するフィールドのみ更新
	if ProductUpdateRequest.Name != nil {
		product.Name = *ProductUpdateRequest.Name	// nilでない場合のみ更新
	}
	if ProductUpdateRequest.Barcode != nil {
		product.Barcode = *ProductUpdateRequest.Barcode
	}
	if ProductUpdateRequest.Price != nil {
		product.Price = *ProductUpdateRequest.Price
	}
	if ProductUpdateRequest.Stock != nil {
		product.Stock = *ProductUpdateRequest.Stock
	}
	if ProductUpdateRequest.TagId != nil {
		product.TagId = *ProductUpdateRequest.TagId
	}
	if ProductUpdateRequest.IsSold != nil {
		product.IsSold = *ProductUpdateRequest.IsSold
	}
	if ProductUpdateRequest.WarningStock != nil {
		product.WarningStock = *ProductUpdateRequest.WarningStock
	}
	if ProductUpdateRequest.SafetyStock != nil {
		product.SafetyStock = *ProductUpdateRequest.SafetyStock
	}

	// DBへ保存
	err = ProductService.UpdateProduct(ProductUpdateRequest.Id, &product)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal update product")
		return
	}
	
	productLog := model.ProductLog{
		SourceId: -1,
		ProductId: int64(product.ID),
		Quantity: -1,
		UnitPrice: product.Price,
		Stock: product.Stock,
	}
	// 入荷商品情報をDBへ保存DBへ保存
	err = ProductService.CreateProductLog(&productLog)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal create arrival_product")
		return
	}

	c.JSON(http.StatusOK, "success")
}

// 商品画像パス更新API
func UpdateProductImagePath(c *gin.Context) {
	ProductService := service.ProductService{}
	ProductUpdateRequest := model.ProductImageUpdateRequest{}
	err := c.Bind(&ProductUpdateRequest)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "request is not correct")
		return
	}

	// 商品情報更新
	// DBへ保存
	err = ProductService.UpdateProductImagePath(ProductUpdateRequest.Id, ProductUpdateRequest.ImagePath)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal update product image")
		return
	}

	c.JSON(http.StatusOK, "success")
}

// 購入情報削除API
func DeletePayment(c *gin.Context) {
	ProductService := service.ProductService{}
	AssetService := service.AssetService{}
	UserService := service.UserService{}
	paymentId, err := strconv.ParseInt(c.Param("paymentId"), 10, 64)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "request is not correct")
		return
	}

	// 購入情報取得
	payment, err := ProductService.GetPaymentById(paymentId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get payment")
		return
	}
	// 購入商品情報取得
	productLogs, err := ProductService.GetProductLogsBySourceId(paymentId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get payment_products")
		return
	}


	// 購入情報削除
	// DBへ保存
	err = ProductService.DeletePayment(paymentId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal delete payment")
		return
	}

	if(payment.Method == "card"){
		// ユーザ残高を増やす
		err = UserService.IncreaseUserDebt(payment.UserId,payment.Price)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, "fetal increase user debt")
			return
		}
		// 商店残高(debt)を増やす
		err = AssetService.IncreaseDebt(payment.Price)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, "fetal increase debt")
			return
		}
	}else {
		// お金を減らす
		err = AssetService.IncreaseMoney(0-payment.Price)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, "fetal decrease money")
			return
		}
	}
	
	// 在庫を増やす
	for _, productLog := range productLogs {
		err = ProductService.IncreaseStock(productLog.ProductId, productLog.Quantity)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, "fetal increase stock")
			return
		}
	}

	c.JSON(http.StatusOK, "success")
}

// 入荷情報削除API
func DeleteArrival(c *gin.Context) {
	ProductService := service.ProductService{}
	AssetService := service.AssetService{}
	arrivalId, err := strconv.ParseInt(c.Param("arrivalId"), 10, 64)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "request is not correct")
		return
	}

	// 入荷情報取得
	arrival, err := ProductService.GetArrivalById(arrivalId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get arrival")
		return
	}
	arrivalProducts, err := ProductService.GetProductLogsBySourceId(arrivalId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get arrival_product")
		return
	}

	// 入荷情報削除
	// DBへ保存
	err = ProductService.DeleteArrival(arrivalId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal delete arrival")
		return
	}
	// お金を増やす
	err = AssetService.IncreaseMoney(arrival.Money)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal increase money")
		return
	}
	// 在庫を減らす
	for _, arrivalProduct := range arrivalProducts {
		err = ProductService.IncreaseStock(arrivalProduct.ProductId, 0-arrivalProduct.Quantity)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, "fetal decrease stock")
			return
		}
	}

	c.JSON(http.StatusOK, "success")
}
