package service

import (
	"context"
	"log"

	"cloud.google.com/go/logging"
)

func createCloudLog(logMessage string)(error) {
	ctx := context.Background()
	client, err := logging.NewClient(ctx, "kajilab-store")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return err
	}
	defer client.Close()

	//ログインスタンスの作成
	logger := client.Logger("backend").StandardLogger(logging.Info)

	// ログ出力の例
	logger.Println(logMessage)
	
	// 返す
	return nil
}