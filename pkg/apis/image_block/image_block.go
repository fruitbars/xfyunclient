// pkg/image_block/image_block.go

package image_block

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fruitbars/xfyunclient/pkg/ase"
	"github.com/fruitbars/xfyunclient/pkg/models/ocr_multi_lang"
	"github.com/fruitbars/xfyunclient/pkg/models/ocr_universal_2024/ocr_universal_2024_engine_text"
	"github.com/fruitbars/xfyunclient/pkg/models/ocr_universal_2024/ocr_universal_2024_request"
	"github.com/fruitbars/xfyunclient/pkg/models/ocr_universal_2024/ocr_universal_2024_response"
	"github.com/fruitbars/xfyunclient/pkg/utils"
)

// BlockRecognitionClient 定义了块识别客户端的接口
type BlockRecognitionClient interface {
	// RecognizeBlocksFromFile 从文件中识别文本块
	RecognizeBlocksFromFile(filePath string, language string) (map[string]string, error)

	// RecognizeBlocksFromImage 从图像数据中识别文本块
	RecognizeBlocksFromImage(format string, imageData []byte, language string) (map[string]string, error)

	// RecognizeBlocksFromReader 从Reader中识别文本块
	RecognizeBlocksFromReader(reader io.Reader, language string) (map[string]string, error)

	// SaveResultToFile 将识别结果保存到文件
	SaveResultToFile(data []byte, language, protoc, version string) (string, error)
}

// Client 实现了BlockRecognitionClient接口
type Client struct {
	appID        string
	apiKey       string
	apiSecret    string
	serverURL    string
	outputDir    string
	debug        bool
	useUniversal bool // 是否使用2024版通用OCR
}

// NewClient 创建一个新的块识别客户端
func NewClient(appID, apiKey, apiSecret, serverURL string) *Client {
	return &Client{
		appID:        appID,
		apiKey:       apiKey,
		apiSecret:    apiSecret,
		serverURL:    serverURL,
		outputDir:    "./output",
		debug:        false,
		useUniversal: false,
	}
}

// SetDebug 设置是否启用调试模式
func (c *Client) SetDebug(debug bool) {
	c.debug = debug
}

// SetOutputDirectory 设置结果文件的输出目录
func (c *Client) SetOutputDirectory(dir string) error {
	// 检查目录是否存在，不存在则创建
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	c.outputDir = dir
	return nil
}

// UseUniversalOCR 设置是否使用2024版通用OCR
func (c *Client) UseUniversalOCR(enable bool) {
	c.useUniversal = enable
}

// logDebug 打印调试信息
func (c *Client) logDebug(format string, args ...interface{}) {
	if c.debug {
		fmt.Printf("[ImageBlock Debug] "+format+"\n", args...)
	}
}

// RecognizeBlocksFromFile 从文件中识别文本块
func (c *Client) RecognizeBlocksFromFile(filePath string, language string) (map[int]string, error) {
	// 验证文件路径
	if filePath == "" {
		return nil, errors.New("file path cannot be empty")
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", filePath)
	}

	// 验证语言参数
	if language == "" {
		return nil, errors.New("language cannot be empty")
	}

	// 读取图像文件
	format, imageData, err := utils.ReadImageFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read image file: %w", err)
	}

	c.logDebug("Image format: %s, size: %d bytes", format, len(imageData))

	// 调用RecognizeBlocksFromImage方法
	return c.RecognizeBlocksFromImage(format, imageData, language)
}

// RecognizeBlocksFromImage 从图像数据中识别文本块
func (c *Client) RecognizeBlocksFromImage(format string, imageData []byte, language string) (map[int]string, error) {
	// 验证参数
	if len(imageData) == 0 {
		return nil, errors.New("image data cannot be empty")
	}

	if language == "" {
		return nil, errors.New("language cannot be empty")
	}

	// 检查图像格式
	format = strings.ToLower(format)
	if format != "jpg" && format != "jpeg" && format != "png" && format != "bmp" {
		return nil, fmt.Errorf("unsupported image format: %s", format)
	}

	// 创建ASE HTTP客户端
	client := ase.NewASEHttpClient(
		c.serverURL,
		c.appID,
		c.apiKey,
		c.apiSecret,
		"", // 使用默认HTTP协议
		"", // 使用默认算法
	)

	c.logDebug("Created ASE client for server URL: %s", c.serverURL)

	// 根据设置选择OCR版本
	if c.useUniversal {
		return c.recognizeUniversal2024(client, format, imageData, language)
	} else {
		return c.recognizeMultiLang(client, format, imageData, language)
	}
}

// recognizeMultiLang 使用多语言OCR进行识别
func (c *Client) recognizeMultiLang(client *ase.ASEHttpClient, format string, imageData []byte, language string) (map[int]string, error) {
	// 创建OCR请求
	req := ocr_multi_lang.NewOcrRequest(c.appID, language, format, imageData)

	// 调用ASE API
	response, err := client.CallASEAPIJson(req)
	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}

	// 验证响应数据
	if response == nil {
		return nil, errors.New("empty response from API")
	}

	c.logDebug("Received %d bytes of response data", len(response))

	// 解析响应
	var ocrResponse ocr_multi_lang.Response
	if err := json.Unmarshal(response, &ocrResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// 检查响应状态
	if ocrResponse.Header.Code != 0 {
		return nil, fmt.Errorf("OCR error: code=%d, message=%s",
			ocrResponse.Header.Code,
			ocrResponse.Header.Message,
		)
	}

	// 解码识别文本
	if ocrResponse.Payload.OCROutputText.Text == "" {
		return nil, errors.New("empty OCR result")
	}

	textData, err := base64.StdEncoding.DecodeString(ocrResponse.Payload.OCROutputText.Text)
	if err != nil {
		return nil, fmt.Errorf("failed to decode result text: %w", err)
	}

	// 解析OCR引擎结果
	var ocrEngineText ocr_multi_lang.OCRResponseText
	if err := json.Unmarshal(textData, &ocrEngineText); err != nil {
		return nil, fmt.Errorf("failed to parse engine result: %w", err)
	}

	// 保存结果到文件
	outputFile, err := c.SaveResultToFile(textData, language, ocrEngineText.Protoc, ocrEngineText.Version)
	if err != nil {
		c.logDebug("Warning: failed to save result file: %v", err)
	} else {
		c.logDebug("Result saved to file: %s", outputFile)
	}

	// 获取块内容
	blockContents := ocr_multi_lang.GetBlockContents(&ocrEngineText)
	c.logDebug("Recognized %d text blocks", len(blockContents))

	return blockContents, nil
}

// recognizeUniversal2024 使用2024版通用OCR进行识别
func (c *Client) recognizeUniversal2024(client *ase.ASEHttpClient, format string, imageData []byte, language string) (map[int]string, error) {
	// 创建OCR请求
	req := ocr_universal_2024_request.NewOcrUniversal2024Request(c.appID, language, format, imageData)

	// 调用ASE API
	response, err := client.CallASEAPIJson(req)
	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}

	// 验证响应数据
	if response == nil {
		return nil, errors.New("empty response from API")
	}

	c.logDebug("Received %d bytes of response data", len(response))

	// 解析响应
	var ocrResponse ocr_universal_2024_response.OCRResponse
	if err := json.Unmarshal(response, &ocrResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// 检查响应状态
	if ocrResponse.Header.Code != 0 {
		return nil, fmt.Errorf("OCR error: code=%d, message=%s",
			ocrResponse.Header.Code,
			ocrResponse.Header.Message,
		)
	}

	// 解码识别文本
	if ocrResponse.Payload.OCROutputText.Text == "" {
		return nil, errors.New("empty OCR result")
	}

	textData, err := base64.StdEncoding.DecodeString(ocrResponse.Payload.OCROutputText.Text)
	if err != nil {
		return nil, fmt.Errorf("failed to decode result text: %w", err)
	}

	// 解析OCR引擎结果
	var ocrEngineText ocr_universal_2024_engine_text.OCREngineText
	if err := json.Unmarshal(textData, &ocrEngineText); err != nil {
		return nil, fmt.Errorf("failed to parse engine result: %w", err)
	}

	// 保存结果到文件
	outputFile, err := c.SaveResultToFile(textData, language, ocrEngineText.Protoc, ocrEngineText.Version)
	if err != nil {
		c.logDebug("Warning: failed to save result file: %v", err)
	} else {
		c.logDebug("Result saved to file: %s", outputFile)
	}

	// 获取块内容
	blockContents := ocr_multi_lang.GetUniversal2024BlockContents(&ocrEngineText)
	c.logDebug("Recognized %d text blocks", len(blockContents))

	return blockContents, nil
}

// RecognizeBlocksFromReader 从Reader中识别文本块
func (c *Client) RecognizeBlocksFromReader(reader io.Reader, language string) (map[int]string, error) {
	// 读取全部数据
	imageData, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	// 检测图像格式
	format, err := utils.DetectImageFormat(imageData)
	if err != nil {
		return nil, fmt.Errorf("failed to detect image format: %w", err)
	}

	// 调用RecognizeBlocksFromImage方法
	return c.RecognizeBlocksFromImage(format, imageData, language)
}

// SaveResultToFile 将识别结果保存到文件
func (c *Client) SaveResultToFile(data []byte, language, protoc, version string) (string, error) {
	// 确保输出目录存在
	if err := os.MkdirAll(c.outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// 创建文件名
	timestamp := time.Now().Format("20060102_150405")
	sanitizedURL := strings.ReplaceAll(c.serverURL, "/", "_")
	sanitizedURL = strings.ReplaceAll(sanitizedURL, ":", "_")

	fileName := fmt.Sprintf("%s_resp_%s_%s_%s_%s.json",
		sanitizedURL,
		language,
		protoc,
		version,
		timestamp,
	)

	// 完整文件路径
	filePath := filepath.Join(c.outputDir, fileName)

	// 写入文件
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write result file: %w", err)
	}

	return filePath, nil
}

// ProcessDirectory 处理目录中的所有图像文件
func (c *Client) ProcessDirectory(dirPath, language string) (map[string]map[int]string, error) {
	// 检查目录是否存在
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory does not exist: %s", dirPath)
	}

	// 读取目录内容
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	// 结果映射: 文件名 -> 文本块
	results := make(map[string]map[int]string)

	// 处理每个文件
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		ext := strings.ToLower(filepath.Ext(fileName))

		// 只处理图像文件
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".bmp" {
			continue
		}

		filePath := filepath.Join(dirPath, fileName)
		c.logDebug("Processing file: %s", filePath)

		// 识别文本块
		blocks, err := c.RecognizeBlocksFromFile(filePath, language)
		if err != nil {
			c.logDebug("Error processing %s: %v", filePath, err)
			continue
		}

		// 存储结果
		results[fileName] = blocks
	}

	return results, nil
}

// 为了兼容旧代码而保留的函数
func OcrBlockTest(aseAppid, aseAPIKey, aseAPISecret string, serverUrl string, language string, fname string) (string, error) {
	client := NewClient(aseAppid, aseAPIKey, aseAPISecret, serverUrl)

	blocks, err := client.RecognizeBlocksFromFile(fname, language)
	if err != nil {
		return "", err
	}

	result := utils.MapToString(blocks)
	return result, nil
}

// 为了兼容旧代码而保留的函数
func OcrUniversal2024Test(aseAppid, aseAPIKey, aseAPISecret string, serverUrl string, language string, fname string) (string, error) {
	client := NewClient(aseAppid, aseAPIKey, aseAPISecret, serverUrl)
	client.UseUniversalOCR(true)

	blocks, err := client.RecognizeBlocksFromFile(fname, language)
	if err != nil {
		return "", err
	}

	result := utils.MapToString(blocks)
	return result, nil
}
