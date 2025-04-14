package main

import (
	"fmt"
	"log"
	"os"
	"xfyunclient/pkg/interface/translation"

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

	log.Println("appid:", appid, "apiSecret:", apiSecret, ",apiKey:", apiKey)

	serverUrl := "https://itrans.xf-yun.com/v1/its"

	result, err := translation.Translate(appid, apiKey, apiSecret, serverUrl, "cn", "en", "你好")
	if err != nil {
		fmt.Println("Error occurred:", err)
		return
	}
	fmt.Println(result)
}
