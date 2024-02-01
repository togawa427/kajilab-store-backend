package main

import (
	"database/sql"
	"fmt"
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

	sqlStmt:= `
	create table products (id integer not null primary key, name text);
	delete from foo;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		fmt.Println("強制終了")
		return
	}

	// データ挿入
	_, err = db.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
	if err != nil {
		log.Fatal(err)
	}
}