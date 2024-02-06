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

// 購入情報を取得
func (ProductService) GetPaymentById(id int64) (model.Payment, error) {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return model.Payment{}, err
	}

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

// 入荷情報を登録
func (ProductService) CreateArrival(arrival *model.Arrival) (int64, error) {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	// 入荷情報をDBへ保存
	result := db.Create(arrival)
	if result.Error != nil {
		fmt.Printf("入荷情報登録失敗 %v", result.Error)
		return 0,result.Error
	}
	return int64(arrival.ID), nil
}

// 入荷した商品情報を登録
func (ProductService) CreateArriveProduct(arriveProduct *model.ArrivalProduct) (error) {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 入荷商品をDBへ登録
	result := db.Create(arriveProduct)
	if result.Error != nil {
		fmt.Printf("入荷商品情報登録失敗 %v", result.Error)
		return result.Error
	}
	return nil
}

// 商品情報の更新
func (ProductService) UpdateProduct(id int64,product *model.Product) (error) {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 商品情報をDBへ登録
	result := db.Model(&model.Product{}).Where("id = ?", id).Updates(&product)
	if result.Error != nil {
		fmt.Printf("商品情報更新失敗 %v", result.Error)
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

	// 購入情報をDBから削除
	result := db.Delete(&model.Payment{}, id)
	if result.Error != nil {
		fmt.Printf("購入情報削除失敗 %v", result.Error)
		return result.Error
	}

	// 購入商品情報をDBから削除
	result = db.Where("payment_id LIKE ?", id).Delete(&model.PaymentProduct{})
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

	// 入荷情報をDBから削除
	result := db.Delete(&model.Arrival{}, id)
	if result.Error != nil {
		fmt.Printf("入荷情報削除失敗 %v", result.Error)
		return result.Error
	}

	// 入荷商品情報をDBから削除
	result = db.Where("arrival_id LIKE ?", id).Delete(&model.ArrivalProduct{})
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

// 購入商品情報を取得
func (ProductService) GetPaymentProductsByPaymentId(paymentId int64)([]model.PaymentProduct, error) {
	paymentProducts := []model.PaymentProduct{}
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return paymentProducts, err
	}
	// 取得
	result := db.Where(&model.PaymentProduct{PaymentId: paymentId}).Find(&paymentProducts)
	if result.Error != nil {
		fmt.Printf("購入商品取得失敗 %v", result.Error)
		return paymentProducts, result.Error
	}
	return paymentProducts, nil
}

// 入荷商品情報を取得
func (ProductService) GetArrivalProductsByArrivalId(arrivalId int64) ([]model.ArrivalProduct, error) {
	arrivalProducts := []model.ArrivalProduct{}
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return arrivalProducts, err
	}
	// 取得
	result := db.Where(&model.ArrivalProduct{ArrivalId: arrivalId}).Find(&arrivalProducts)
	if result.Error != nil {
		fmt.Printf("入荷商品取得失敗 %v", result.Error)
		return arrivalProducts, result.Error
	}
	return arrivalProducts, nil
}