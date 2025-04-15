package main

import (
	"fmt"
	"github.com/fruitbars/xfyunclient/pkg/apis/v2tts"
	"github.com/fruitbars/xfyunclient/pkg/utils"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func ttsclientTest(appID, apiKey, apiSecret, serverUrl string) {
	ttsConfig := v2tts.DefaultConfig()
	ttsConfig.AppID = appID
	ttsConfig.ASEAPIKey = apiKey
	ttsConfig.ASESecret = apiSecret
	ttsConfig.ServerURL = serverUrl
	ttsConfig.DefaultAue = "raw" // PCM原始音频格式

	client := v2tts.NewClientWithConfig(ttsConfig)

	client.TextToSpeechToFile("这真的有问题吗", "ttssyn.pcm")

	utils.PcmToWav("ttssyn.pcm", "ttssyn.pcm.wav", 16000, 1, 16)
}

func ttsAPItest(appID, apiKey, apiSecret, serverUrl string) {
	sid, err := v2tts.V2TextToSpeechToFileDefault(appID, apiKey, apiSecret, serverUrl, "这真的有问题吗", "raw", "xiaoyan", "UTF8", "fanyi.pcm")
	if err != nil {
		fmt.Println("Error occurred:", err)
		return
	}

	utils.PcmToWav("fanyi.pcm", "fanyi.wav", 16000, 1, 16)

	log.Println("sid:", sid)
}

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

	ttsclientTest(appid, apiKey, apiSecret, serverUrl)

}
