package service

import (
	"errors"
	"fmt"
	"kajilab-store-backend/model"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ProductService struct{}

// 全ての商品を取得する
func (ProductService) GetAllProducts(limit int64, offset int64) ([]model.Product, error){
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	products := make([]model.Product, 0)
	//result := db.Order("name").Find(&products)
	result := db.Order("stock DESC").Offset(int(offset)).Limit(int(limit)).Find(&products)
	if result.Error != nil {
		fmt.Printf("商品取得失敗 %v", result.Error)
		return nil, result.Error
	}
	return products, nil
}

// 購入ログを取得
func (ProductService) GetBuyLogs(limit int64) ([]model.Payment, error){
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	logs := make([]model.Payment, 0)
	//result := db.Order("name").Find(&products)
	result := db.Limit(int(limit)).Find(&logs)
	if result.Error != nil {
		fmt.Printf("購入履歴取得失敗 %v", result.Error)
		return nil, result.Error
	}
	return logs, nil
}

func (ProductService) GetBuyProductsByPaymentId(paymentId int64) ([]model.PaymentProduct, error){
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	products := make([]model.PaymentProduct, 0)
	result := db.Where("payment_id = ?", paymentId).Find(&products)
	if result.Error != nil {
		fmt.Printf("購入物取得失敗 %v", result.Error)
		return nil, result.Error
	}
	return products, nil
}

// 購入物を取得
func (ProductService) GetBuyProducts() ([]model.PaymentProduct, error){
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	products := make([]model.PaymentProduct, 0)
	result := db.Find(&products)
	if result.Error != nil {
		fmt.Printf("購入物取得失敗 %v", result.Error)
		return nil, result.Error
	}
	return products, nil
}

// IDから商品情報を取得
func (ProductService) GetProductById(id int64) (model.Product, error) {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return model.Product{}, err
	}

	var product model.Product
	result := db.First(&product, id)
	if result.Error != nil {
		fmt.Printf("購入物取得失敗 %v", result.Error)
		return product, result.Error
	}
	return product, nil
}

// 入荷ログを取得
func (ProductService) GetArriveLogs(limit int64) ([]model.Arrival, error){
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	logs := make([]model.Arrival, 0)
	result := db.Limit(int(limit)).Find(&logs)
	if result.Error != nil {
		fmt.Printf("入荷履歴取得失敗 %v", result.Error)
		return nil, result.Error
	}
	return logs, nil
}

// 入荷IDから入荷商品の取得
func (ProductService) GetArriveProductsByArriveId(arriveId int64) ([]model.ArrivalProduct, error){
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	products := make([]model.ArrivalProduct, 0)
	result := db.Where("arrival_id = ?", arriveId).Find(&products)
	if result.Error != nil {
		fmt.Printf("入荷商品取得失敗 %v", result.Error)
		return nil, result.Error
	}
	return products, nil
}

// 商品情報を登録
func (ProductService) CreateProduct(product *model.Product) error {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	// バーコードのかぶりチェック
	existProduct := model.Product{}
	result := db.Where("barcode = ?", product.Barcode).First(&existProduct)
	if result.Error == nil {
		err := errors.New("the barcode is existing")
		fmt.Printf("%v",err)
		return err
	}

	// 商品をDBへ登録
	result = db.Create(product)
	if result.Error != nil {
		fmt.Printf("入荷登録失敗 %v", result.Error)
		return result.Error
	}
	return nil
}

// 購入情報を登録
func (ProductService) CreatePayment(payment *model.Payment) (int64, error) {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return 0,err
	}

	// 購入情報をDBへ登録
	result := db.Create(payment)
	if result.Error != nil {
		fmt.Printf("購入情報登録失敗 %v", result.Error)
		return 0,result.Error
	}
	return int64(payment.ID),nil
}

// 購入した商品の情報を登録
func (ProductService) CreatePaymentProduct(payment *model.PaymentProduct) (error) {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 購入商品をDBへ保存
	result := db.Create(payment)
	if result.Error != nil {
		fmt.Printf("購入商品登録失敗 %v", result.Error)
		return result.Error
	}
	return nil
}