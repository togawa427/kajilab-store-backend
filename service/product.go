package service

import (
	"fmt"
	"kajilab-store-backend/model"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ProductService struct{}

// 全ての商品を取得する
func (ProductService) GetAllProducts() ([]model.Product, error){
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})

	if err != nil {
		fmt.Println(err)
	}

	// cmd := "SELECT * FROM products"
	// out, err := db.Exec(cmd)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(out)
	products := make([]model.Product, 0)
	result := db.Find(&products)
	if result.Error != nil {
		fmt.Printf("商品取得失敗 %v", result.Error)
		return nil, result.Error
	}
	return products, nil
}

