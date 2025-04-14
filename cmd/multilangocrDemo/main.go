package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"xfyunclient/pkg/interface/ocr_universal"
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

	log.Println("appid:", appid, "apiSecret:", apiSecret, ",apiKey:", apiKey)

	serverUrl := "https://api.xf-yun.com/v1/private/sf8e6aca1"
	fname := "fanyi1.jpg"

	ocr_universal.OcrUniversal(appid, apiKey, apiSecret, serverUrl, "ch_en_public_cloud", fname)

}
