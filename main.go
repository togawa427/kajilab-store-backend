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
		versionEngine.GET("/products", controller.GetAllProducts)
		//versionEngine.GET("/products/:product_id", controller.GetProductByProductId)

		// versionEngine.GET("/stayers", controller.Stayer)
		// versionEngine.POST("/stayers", controller.Beacon)
		// versionEngine.GET("/logs", controller.Log)
		// versionEngine.GET("/logs/gantt", controller.LogGantt)
		// versionEngine.GET("/users", controller.PastUserList) // 編集機能のフロントのブランチがマージされたら消す
		// versionEngine.GET("/users/:communityId", controller.UserList)
		// versionEngine.GET("/users/extended", controller.ExtendedUserList)
		// versionEngine.POST("/users", controller.CreateUser)
		// versionEngine.PUT("/users", controller.UpdateUser)
		// versionEngine.DELETE("/users/:userId", controller.DeleteUser)
		// versionEngine.GET("/admin/users/:communityId", controller.AdminUserList)
		// versionEngine.GET("/check", controller.Check)
		// versionEngine.POST("/attendance", controller.Attendance)
		// versionEngine.GET("/rooms/:communityID", controller.GetRoomsByCommunityID)
		// versionEngine.PUT("/rooms", controller.UpdateRoom)
		// versionEngine.GET("/tags/:communityId/names", controller.GetTagNamesByCommunityId)
		// versionEngine.GET("/tags/:communityId", controller.GetTagsByCommunityIdHandler)
		// versionEngine.GET("/beacons", controller.GetBeacon)
		// versionEngine.GET("/communities/:userId", controller.GetCommunityByUserIdHandler)
		// versionEngine.GET("/buildings/editor", controller.GetBuildingsEditor)
		// versionEngine.GET("/signup", controller.SignUp)
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

	// サンプルデータの挿入
	db.Create(&model.Product{Name: "じゃがりこサラダ味", Barcode: 134912341234, Price: 120, Stock: 9, ImagePath: "public/images/jagariko.jpg"})
	db.Create(&model.Product{Name: "じゃがりこチーズ味", Barcode: 134912341233, Price: 120, Stock: 4, ImagePath: "public/images/jagariko.jpg"})
	db.Create(&model.Product{Name: "じゃがりこたらこ味", Barcode: 134912341232, Price: 120, Stock: 5, ImagePath: "public/images/jagariko.jpg"})
	db.Create(&model.Product{Name: "じゃがりこじゃがバター味", Barcode: 134912341231, Price: 120, Stock: 9, ImagePath: "public/images/jagariko.jpg"})
	db.Create(&model.Product{Name: "ごつもりソース味", Barcode: 134912341221, Price: 140, Stock: 11, ImagePath: "public/images/jagariko.jpg"})
	db.Create(&model.Product{Name: "ごつもり塩味", Barcode: 134912341222, Price: 140, Stock: 10, ImagePath: "public/images/jagariko.jpg"})

	
}