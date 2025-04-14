package main

import (
	"github.com/fruitbars/xfyunclient/pkg/apis/ost"
	"github.com/joho/godotenv"
	"log"
	"os"
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

	client := ost.NewClient(appid, apiKey, apiSecret)

	// 识别音频文件
	result, err := client.RecognizeFile("402_1744111790.mp4.wav", nil)
	if err != nil {
		log.Fatalf("识别失败: %v", err)
	}

	ufResult, err := ost.FormatAsUserFriendlyResult(result)
	// 输出识别结果
	log.Println("识别结果: ", ufResult.Data.Result.FullText)

	// 输出识别结果
	log.Println("识别结果: ", ufResult.Data.Result.FullText2)
}
