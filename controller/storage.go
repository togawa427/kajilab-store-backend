package controller

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)


func GetStorage(c *gin.Context) {

	bucketName  := "kajilab-store.appspot.com"	// GCSバケット名
	credentialsFile := "kajilab-store-cf40dbeb6615.json"	// サービスアカウント鍵ファイルのパス

	now := time.Now()
	formatedNow := now.Format("2006-01-02_15-04-05")
	localFileName := "backup/db_backup_" + formatedNow + ".bk"
	gcsFileName := localFileName	// GCSバケットのアップロード先のパス

	// ==DBバックアップ実行==
	// 実行したいシェルスクリプトのパス
	scriptPath := "./backup.sh"

	// コマンドを準備
	cmd := exec.Command("bash", scriptPath, localFileName)

	// コマンドの標準出力と標準エラー出力を取得
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// コマンドを実行
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Command execution failed: %v\n%s", err, stderr.String())
	}

	// 成功時の出力を表示
	fmt.Printf("Command output:\n%s", stdout.String())

	// == Cloud Storageへ保存
	// Google Cloud Storageクライアントの作成
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
			log.Fatal(err)
	}

	// アップロードするファイルを開く
	file, err := os.Open(localFileName)
	if err != nil {
			log.Fatal(err)
	}
	defer file.Close()

	// バケットオブジェクトの作成
	bucket := client.Bucket(bucketName)

	// バケット内のアップロード先のオブジェクトを作成
	obj := bucket.Object(gcsFileName)

	// ファイルをアップロード
	wc := obj.NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
			log.Fatal(err)
	}
	if err := wc.Close(); err != nil {
			log.Fatal(err)
	}

	fmt.Println("file upload success")
		

	c.JSON(http.StatusOK, "sucess")
}