package main

import (
	"fmt"
	"github.com/fruitbars/xfyunclient/pkg/interface/ocr_layout"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// 加载.env 文件
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading.env file: %v", err)
	}

	// 从环境变量中读取相应的值
	appid := os.Getenv("APPID")
	apiKey := os.Getenv("APIKEY")
	apiSecret := os.Getenv("APISECRET")

	serverUrl := "https://cn-huadong-1.xf-yun.com/v1/private/s8ea6b8fa"
	fname := "fanyi1.jpg"
	// 调用 OCR 函数
	err, result := ocr_layout.OCRLayoutFile(appid, apiKey, apiSecret, serverUrl, fname)
	if err != nil {
		fmt.Println("Error occurred:", err)
		return
	}
	fmt.Println(result)
}
