package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"path/filepath"
	"strings"
	"xfyunclient/pkg/interface/image_block"
)

type Config struct {
	ASE struct {
		Appid     string `json:"appid"`
		APIKey    string `json:"apikey"`
		APISecret string `json:"apisecret"`
	} `json:"ase"`
	LLM struct {
		Appid     string `json:"appid"`
		APIKey    string `json:"apikey"`
		APISecret string `json:"apisecret"`
	} `json:"llm"`
}

var config Config

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 读取配置文件
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("failed to open config file: %v", err)
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	if err := decoder.Decode(&config); err != nil {
		log.Fatalf("failed to decode config file: %v", err)
	}

	directory := "./OCRdata_xsq"
	serverUrl := "https://cn-huabei-1.xf-yun.com/v1/private/s2a094db9"
	lang := "ch_en"
	outputExcel := "./ocr_results.xlsx"

	if err := processDirectory(directory, serverUrl, lang, outputExcel); err != nil {
		log.Fatalf("failed to process directory: %v", err)
	}

	log.Println("OCR results have been written to", outputExcel)
}

func processDirectory(directory, serverUrl, lang, outputExcel string) error {
	f := excelize.NewFile()
	index, _ := f.NewSheet("Sheet1")
	f.SetCellValue("Sheet1", "A1", "文件名")
	f.SetCellValue("Sheet1", "B1", "OCR结果")
	f.SetCellValue("Sheet1", "C1", "LLM结果general")
	f.SetCellValue("Sheet1", "D1", "LLM结果generalv3.5")

	row := 2
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			result, err := image_block.OcrBlockTest(config.ASE.Appid, config.ASE.APIKey, config.ASE.APISecret, serverUrl, lang, path)
			if err != nil {
				log.Printf("failed to process file %s: %v", info.Name(), err)
				return nil
			}

			kvResultGeneral, _ := llmKeyValue("general", result)
			kvResultGeneralv35, _ := llmKeyValue("generalv3.5", result)

			// 写入文件名和OCR结果
			f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), path)
			f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), result)
			f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), kvResultGeneral)
			f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), kvResultGeneralv35)

			row++
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk directory: %v", err)
	}

	f.SetActiveSheet(index)
	if err := f.SaveAs(outputExcel); err != nil {
		return fmt.Errorf("failed to save Excel file: %v", err)
	}

	return nil
}

func llmKeyValue(model string, ocrResult string) (string, error) {
	prompt := fmt.Sprintf("从以下文本中找出显示器或者电视机的品牌，只输出品牌名称即可：\n\n%s", ocrResult)

	prompt = strings.ReplaceAll(prompt, "\n", "\r\n")
	token := config.LLM.APIKey + ":" + config.LLM.APISecret
	conf := openai.DefaultConfig(token)
	conf.BaseURL = "https://spark-api-open.xf-yun.com/v1"
	client := openai.NewClientWithConfig(conf)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens: 4096,
		},
	)

	if err != nil {
		log.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
