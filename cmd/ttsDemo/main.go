package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"xfyunclient/pkg/apis/v2tts/backup"
	"xfyunclient/pkg/utils"
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

	serverUrl := "wss://tts-api.xfyun.cn/v2/tts"

	// 调用 OCR 函数
	sid, err := backup.V2TextToSpeechToFileDefault(appid, apiKey, apiSecret, serverUrl, "这真的有问题吗", "raw", "xiaoyan", "UTF8", "fanyi.pcm")
	if err != nil {
		fmt.Println("Error occurred:", err)
		return
	}

	utils.PcmToWav("fanyi.pcm", "fanyi.wav", 16000, 1, 16)

	log.Println("sid:", sid)
}
