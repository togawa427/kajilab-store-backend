package main

import (
	"log"
	"os"
	"time"

	controller "kajilab-store-backend/controller"
	"kajilab-store-backend/model"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {

	// 環境変数の設定
	os.Setenv("DB_FILE_NAME", "kajilabstore.db")
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "kajilab-store-256c6c01f9cb.json")

	// DBの設定
	log.Println("Start Database")
	_, err := os.Stat(os.Getenv("DB_FILE_NAME"))
	if err != nil {
		// ファイルが存在しないとき初期化を行う
		SetUpDatabase()
	}

	// サーバ起動
	log.Println("Start Server")
	SetUpServer().Run(":8080")
	// v1.GET("/list/simultaneous/:user_id", controller.SimultaneousStayUserList
}

func SetUpServer() *gin.Engine {
	engine := gin.Default()
	// ミドルウェア
	// engine.Use(middleware.RecordUaAndTime)
	// CRUD 書籍
	engine.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*", "http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
		},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))

	versionEngine := engine.Group("api/v1")
	{
		// products
		versionEngine.GET("/products", controller.GetAllProducts)
		versionEngine.GET("/products/:barcode", controller.GetProductByBarcode)
		versionEngine.GET("/products/buy/logs", controller.GetBuyLogs)
		versionEngine.GET("/products/buy/logs/user/:userId", controller.GetBuyLogsByUser)
		versionEngine.GET("/products/arrive/logs", controller.GetArriveLogs)
		versionEngine.GET("/products/stock/:productId", controller.GetProductStockLogsById)
		versionEngine.POST("/products", controller.CreateProduct)
		versionEngine.POST("/products/buy", controller.BuyProducts)
		versionEngine.POST("/products/arrive", controller.ArriveProducts)
		versionEngine.PUT("/products", controller.UpdateProduct)
		versionEngine.PUT("/products/image", controller.UpdateProductImagePath)
		versionEngine.DELETE("/products/buy/:paymentId", controller.DeletePayment)
		versionEngine.DELETE("/products/arrival/:arrivalId", controller.DeleteArrival)

		// assets
		versionEngine.GET("/assets", controller.GetAsset)
		versionEngine.GET("/assets/history", controller.GetAssetHistory)
		versionEngine.PUT("/assets", controller.UpdateAsset)

		// users
		versionEngine.GET("/users/:barcode", controller.GetUserByBarcode)
		versionEngine.POST("/users", controller.CreateUser)
		versionEngine.PUT("/users/debt", controller.UpdateUserDebt)
		versionEngine.PUT("/users/barcode", controller.UpdateUserBarcode)

		//versionEngine.GET("/products/:product_id", controller.GetProductByProductId)

	}

	return engine
}

func SetUpDatabase() {
	log.Println("init database")
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// テーブルの作成
	db.AutoMigrate(&model.Product{})
	db.AutoMigrate(&model.Asset{})
	db.AutoMigrate(&model.Payment{})
	db.AutoMigrate(&model.PaymentProduct{})
	db.AutoMigrate(&model.Arrival{})
	db.AutoMigrate(&model.ArrivalProduct{})
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.ProductLog{})

	// サンプルデータの挿入
	db.Create(&model.Product{Name: "じゃがりこサラダ味", Barcode: 134912341234, Price: 120, Stock: 9, TagId: 1, ImagePath: "public/images/jagariko.jpg"})
	db.Create(&model.Product{Name: "じゃがりこチーズ味", Barcode: 134912341233, Price: 120, Stock: 4, TagId: 1, ImagePath: "public/images/jagariko.jpg"})
	db.Create(&model.Product{Name: "じゃがりこたらこ味", Barcode: 134912341232, Price: 120, Stock: 5, TagId: 1, ImagePath: "public/images/jagariko.jpg"})
	db.Create(&model.Product{Name: "じゃがりこじゃがバター味", Barcode: 134912341231, Price: 120, Stock: 9, TagId: 1, ImagePath: "public/images/jagariko.jpg"})
	db.Create(&model.Product{Name: "ごつもりソース味", Barcode: 134912341221, Price: 140, Stock: 11, TagId: 2, ImagePath: "public/images/jagariko.jpg"})
	db.Create(&model.Product{Name: "ごつもり塩味", Barcode: 134912341222, Price: 140, Stock: 10, TagId: 2, ImagePath: "public/images/jagariko.jpg"})

	db.Create(&model.Payment{Price: 240, PayAt: time.Now(), Method: "cash"})
	db.Create(&model.Payment{Price: 120, PayAt: time.Now(), Method: "cash"})
	db.Create(&model.Payment{Price: 400, PayAt: time.Now(), Method: "cash"})
	db.Create(&model.Payment{Price: 240, PayAt: time.Now(), Method: "cash"})
	db.Create(&model.Payment{Price: 400, PayAt: time.Now(), Method: "cash"})
	db.Create(&model.Payment{Price: 140, PayAt: time.Now(), Method: "cash"})

	db.Create(&model.Arrival{Money: 2400, ArriveAt: time.Now()})
	db.Create(&model.Arrival{Money: 1231, ArriveAt: time.Now()})
	db.Create(&model.Arrival{Money: 413, ArriveAt: time.Now()})
	db.Create(&model.Arrival{Money: 2234, ArriveAt: time.Now()})
	db.Create(&model.Arrival{Money: 1231, ArriveAt: time.Now()})
	db.Create(&model.Arrival{Money: 1941, ArriveAt: time.Now()})

	db.Create(&model.PaymentProduct{PaymentId: 1, ProductId: 1, Quantity: 2, UnitPrice: 120})
	db.Create(&model.PaymentProduct{PaymentId: 3, ProductId: 2, Quantity: 1, UnitPrice: 120})
	db.Create(&model.PaymentProduct{PaymentId: 3, ProductId: 5, Quantity: 1, UnitPrice: 140})
	db.Create(&model.PaymentProduct{PaymentId: 3, ProductId: 6, Quantity: 1, UnitPrice: 140})

	db.Create(&model.ArrivalProduct{ArrivalId: 1, ProductId: 1, Quantity: 10})
	db.Create(&model.ArrivalProduct{ArrivalId: 3, ProductId: 2, Quantity: 8})
	db.Create(&model.ArrivalProduct{ArrivalId: 1, ProductId: 3, Quantity: 12})
	db.Create(&model.ArrivalProduct{ArrivalId: 1, ProductId: 5, Quantity: 4})
	db.Create(&model.ArrivalProduct{ArrivalId: 4, ProductId: 2, Quantity: 10})

	db.Create(&model.Asset{Money: 10000, Debt: 0})
}
