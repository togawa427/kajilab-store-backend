package service

import (
	"errors"
	"fmt"
	"kajilab-store-backend/model"
	"os"
	"time"

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
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer sqlDB.Close()

	products := make([]model.Product, 0)
	//result := db.Order("name").Find(&products)
	result := db.Order("stock DESC").Offset(int(offset)).Limit(int(limit)).Find(&products)
	if result.Error != nil {
		fmt.Printf("商品取得失敗 %v", result.Error)
		return nil, result.Error
	}
	return products, nil
}

// バーコードから商品情報取得
func (ProductService) GetProductByBarcode (barcode int64) (model.Product, error) {
	product := model.Product{}
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return product, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer sqlDB.Close()

	// 取得
	result := db.Where(&model.Product{Barcode: barcode}).First(&product)
	if result.Error != nil {
		fmt.Printf("購入商品取得失敗 %v", result.Error)
		return product, result.Error
	}
	return product, nil
}

// 購入ログを取得
func (ProductService) GetBuyLogs(limit int64) ([]model.Payment, error){
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer sqlDB.Close()

	logs := make([]model.Payment, 0)
	//result := db.Order("name").Find(&products)
	result := db.Order("ID desc").Limit(int(limit)).Find(&logs)
	if result.Error != nil {
		fmt.Printf("購入履歴取得失敗 %v", result.Error)
		return nil, result.Error
	}
	return logs, nil
}

// ユーザIDから購入ログを取得
func (ProductService) GetBuyLogsByUserId(offset int64, limit int64, userId int64) ([]model.Payment, error){
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer sqlDB.Close()

	logs := make([]model.Payment, 0)
	//result := db.Order("name").Find(&products)
	result := db.Where("user_id = ?", int(userId)).Order("ID desc").Offset(int(offset)).Limit(int(limit)).Find(&logs)
	if result.Error != nil {
		fmt.Printf("購入履歴取得失敗 %v", result.Error)
		return nil, result.Error
	}
	return logs, nil
}

// IDから商品情報を取得
func (ProductService) GetProductById(id int64) (model.Product, error) {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return model.Product{}, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer sqlDB.Close()

	var product model.Product
	result := db.First(&product, id)
	if result.Error != nil {
		fmt.Printf("購入物取得失敗 %v", result.Error)
		return product, result.Error
	}
	return product, nil
}

func (ProductService) GetProductLogsByDay(day int64, productId int64) ([]model.ProductLog, error){
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer sqlDB.Close()

	dayAgo := time.Now().Add((0-time.Duration(day)) * 24 * time.Hour)
	logs := make([]model.ProductLog, 0)
	result := db.Where("created_at >= ?", dayAgo).Find(&logs)
	if result.Error != nil {
		fmt.Printf("商品ログ取得失敗 %v", result.Error)
		return nil, result.Error
	}
	return logs, nil
}

// 入荷ログを取得
func (ProductService) GetArriveLogs(limit int64) ([]model.Arrival, error){
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer sqlDB.Close()

	logs := make([]model.Arrival, 0)
	result := db.Order("ID desc").Limit(int(limit)).Find(&logs)
	if result.Error != nil {
		fmt.Printf("入荷履歴取得失敗 %v", result.Error)
		return nil, result.Error
	}
	return logs, nil
}

// 購入情報を取得
func (ProductService) GetPaymentById(id int64) (model.Payment, error) {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return model.Payment{}, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer sqlDB.Close()

	payment := model.Payment{}
	result := db.First(&payment, id)
	if result.Error != nil {
		fmt.Printf("購入情報取得失敗 %v", result.Error)
		return payment, result.Error
	}
	return payment, nil
}

// 入荷情報を取得
func (ProductService) GetArrivalById(id int64) (model.Arrival, error) {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return model.Arrival{}, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer sqlDB.Close()

	arrival := model.Arrival{}
	result := db.First(&arrival, id)
	if result.Error != nil {
		fmt.Printf("入荷商品取得失敗 %v", result.Error)
		return arrival, result.Error
	}
	return arrival, nil
}

// 商品情報を登録
func (ProductService) CreateProduct(product *model.Product) error {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer sqlDB.Close()

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
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer sqlDB.Close()

	// 購入情報をDBへ登録
	result := db.Create(payment)
	if result.Error != nil {
		fmt.Printf("購入情報登録失敗 %v", result.Error)
		return 0,result.Error
	}
	return int64(payment.ID),nil
}


func (ProductService) CreateProductLog(productLog *model.ProductLog) (error) {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer sqlDB.Close()

	// 購入入荷商品をDBへ保存
	result := db.Create(productLog)
	if result.Error != nil {
		fmt.Printf("購入ログ登録失敗 %v", result.Error)
		return result.Error
	}
	return nil
}

// 入荷情報を登録
func (ProductService) CreateArrival(arrival *model.Arrival) (int64, error) {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer sqlDB.Close()

	// 入荷情報をDBへ保存
	result := db.Create(arrival)
	if result.Error != nil {
		fmt.Printf("入荷情報登録失敗 %v", result.Error)
		return 0,result.Error
	}
	return int64(arrival.ID), nil
}

// 商品情報の更新
func (ProductService) UpdateProduct(id int64,product *model.Product) (error) {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer sqlDB.Close()

	// 商品情報をDBへ登録
	result := db.Model(&model.Product{}).Where("id = ?", id).Updates(&product)
	if result.Error != nil {
		fmt.Printf("商品情報更新失敗 %v", result.Error)
		return result.Error
	}
	// 0の場合更新対象から外れるためポインタを使用
	stock := product.Stock
	result = db.Model(&model.Product{}).Where("id = ?", id).Update("stock", &stock)
	if result.Error != nil {
		fmt.Printf("商品個数更新失敗 %v", result.Error)
		return result.Error
	}
	return nil
}

// 商品画像パスの更新
func (ProductService) UpdateProductImagePath(id int64, imagePath string) (error) {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer sqlDB.Close()

	// 商品情報をDBへ登録
	result := db.Model(&model.Product{}).Where("id = ?", id).Update("image_path", imagePath)
	if result.Error != nil {
		fmt.Printf("商品画像パス更新失敗 %v", result.Error)
		return result.Error
	}
	return nil
}

// 購入情報の削除
func (ProductService) DeletePayment(id int64) (error) {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer sqlDB.Close()

	// 購入情報をDBから削除
	result := db.Delete(&model.Payment{}, id)
	if result.Error != nil {
		fmt.Printf("購入情報削除失敗 %v", result.Error)
		return result.Error
	}

	// 購入商品情報ログをDBから削除
	result = db.Where("source_id LIKE ?", id).Delete(&model.ProductLog{})
	if result.Error != nil {
		fmt.Printf("購入商品情報削除失敗 %v", result.Error)
		return result.Error
	}

	return nil
}

// 入荷情報の削除
func (ProductService) DeleteArrival(id int64) (error) {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer sqlDB.Close()

	// 入荷情報をDBから削除
	result := db.Delete(&model.Arrival{}, id)
	if result.Error != nil {
		fmt.Printf("入荷情報削除失敗 %v", result.Error)
		return result.Error
	}

	// 入荷商品情報をDBから削除
	result = db.Where("source_id LIKE ?", id).Delete(&model.ProductLog{})
	if result.Error != nil {
		fmt.Printf("入荷商品情報削除失敗 %v", result.Error)
		return result.Error
	}
	return nil
}

// 在庫を増やす
func (ProductService) IncreaseStock(productId int64, quantity int64)(error) {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer sqlDB.Close()

	// 商品情報を取得
	product := model.Product{}
	result := db.First(&product, productId)
	if result.Error != nil {
		fmt.Printf("商品情報取得失敗 %v", result.Error)
		return result.Error
	}
	// 在庫情報を更新
	result = db.Model(&model.Product{}).Where("id = ?", productId).Update("stock", product.Stock + quantity)
	if result.Error != nil {
		fmt.Printf("在庫更新失敗 %v", result.Error)
		return result.Error
	}
	return nil
}

// 商品ログを取得
func (ProductService) GetProductLogsBySourceId(sourceId int64)([]model.ProductLog, error) {
	productLogs := []model.ProductLog{}
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return productLogs, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer sqlDB.Close()

	// 取得
	result := db.Where(&model.ProductLog{SourceId: sourceId}).Find(&productLogs)
	if result.Error != nil {
		fmt.Printf("購入商品取得失敗 %v", result.Error)
		return productLogs, result.Error
	}
	return productLogs, nil
}