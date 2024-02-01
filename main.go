package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	controller "kajilab-store-backend/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
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
		versionEngine.GET("/products", controller.GetProducts)
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
	db, err := sql.Open("sqlite3", os.Getenv("DB_FILE_NAME"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// productsテーブルの作成
	cmd := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER NOT NULL PRIMARY KEY,
		name TEXT NOT NULL,
		barcode INTEGER UNIQUE,
		price INTEGER,
		stock INTEGER
	);
	`
	_, err = db.Exec(cmd)
	if err != nil {
		log.Printf("failed create products table %q: %s\n", err, cmd)
		return
	}

	// assetsテーブルの作成
	cmd = `
	CREATE TABLE IF NOT EXISTS assets (
		id INTEGER NOT NULL PRIMARY KEY,
		money INTEGER,
		debt INTEGER
	);
	`
	_, err = db.Exec(cmd)
	if err != nil {
		log.Printf("failed create assets table %q: %s\n", err, cmd)
		return
	}

	// paymentsテーブルの作成
	cmd = `
	CREATE TABLE IF NOT EXISTS payments (
		id INTEGER NOT NULL PRIMARY KEY,
		price INTEGER,
		date DATETIME,
		method TEXT
	);
	`
	_, err = db.Exec(cmd)
	if err != nil {
		log.Printf("failed create payments table %q: %s\n", err, cmd)
		return
	}

	// payment_productsテーブルの作成
	cmd = `
	CREATE TABLE IF NOT EXISTS payment_products (
		id INTEGER NOT NULL PRIMARY KEY,
		payment_id INTEGER,
		product_id INTEGER,
		quantity INTEGER,
		unit_price INTEGER
	);
	`
	_, err = db.Exec(cmd)
	if err != nil {
		log.Printf("failed create payment_products table %q: %s\n", err, cmd)
		return
	}

	// arrivalsテーブルの作成
	cmd = `
	CREATE TABLE IF NOT EXISTS arrivals (
		id INTEGER NOT NULL PRIMARY KEY,
		money INTEGER,
		date DATETIME
	);
	`
	_, err = db.Exec(cmd)
	if err != nil {
		log.Printf("failed create arrivals table %q: %s\n", err, cmd)
		return
	}

	// arrival_productsテーブルの作成
	cmd = `
	CREATE TABLE IF NOT EXISTS arrival_products (
		id INTEGER NOT NULL PRIMARY KEY,
		arrival_id INTEGER,
		product_id INTEGER,
		quantity INTEGER
	);
	`
	_, err = db.Exec(cmd)
	if err != nil {
		log.Printf("failed create arrival_products table %q: %s\n", err, cmd)
		return
	}


	// データ挿入
	cmd = "INSERT INTO products (name, barcode, price, stock) VALUES (?, ?, ?, ?);"
	_, err = db.Exec(cmd, "じゃがりこサラダ味", 134912341234, 120, 9)
	if err != nil {
		log.Fatal(err)
	}

	
}