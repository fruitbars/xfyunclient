// pkg/ocr_universal/ocr_universal.go

package ocr_universal

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"xfyunclient/pkg/ase"
	"xfyunclient/pkg/models/ocr_universal"
	"xfyunclient/pkg/utils"
)

// OCRUniversalClient 定义了通用OCR客户端的接口
type OCRUniversalClient interface {
	// RecognizeFile 识别指定文件中的图像文本
	RecognizeFile(filePath string, language string) (*ocr_universal.Result, error)

	// RecognizeImage 识别指定格式和图像数据中的文本
	RecognizeImage(format string, imageData []byte, language string) (*ocr_universal.Result, error)

	// RecognizeReader 从读取器中读取图像数据并识别文本
	RecognizeReader(reader io.Reader, language string) (*ocr_universal.Result, error)
}

// Client 实现了OCRUniversalClient接口
type Client struct {
	appID     string
	apiKey    string
	apiSecret string
	serverURL string
	debug     bool
}

// New 创建一个新的通用OCR客户端
func New(appID, apiKey, apiSecret, serverURL string) *Client {
	return &Client{
		appID:     appID,
		apiKey:    apiKey,
		apiSecret: apiSecret,
		serverURL: serverURL,
		debug:     false,
	}
}

// SetDebug 设置是否启用调试模式
func (c *Client) SetDebug(debug bool) {
	c.debug = debug
}

// logDebug 打印调试信息
func (c *Client) logDebug(format string, args ...interface{}) {
	if c.debug {
		fmt.Printf("[OCR Universal Debug] "+format+"\n", args...)
	}
}

// RecognizeFile 识别指定文件中的图像文本
func (c *Client) RecognizeFile(filePath string, language string) (*ocr_universal.RespResult, error) {
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

	// 调用RecognizeImage方法
	return c.RecognizeImage(format, imageData, language)
}

// RecognizeImage 识别指定格式和图像数据中的文本
func (c *Client) RecognizeImage(format string, imageData []byte, language string) (*ocr_universal.RespResult, error) {
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

	// 创建OCR通用请求
	req := ocr_universal.NewOcrUniversalRequest(c.appID, language, format, imageData)

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
	var ocrResponse ocr_universal.Response
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

	// 解码识别结果
	if ocrResponse.Payload.Result.Text == "" {
		return &ocrResponse.Payload.Result, nil // 无识别结果文本但不是错误
	}

	textData, err := base64.StdEncoding.DecodeString(ocrResponse.Payload.Result.Text)
	if err != nil {
		return nil, fmt.Errorf("failed to decode result text: %w", err)
	}

	// 将解码后的文本设置回结果
	result := ocrResponse.Payload.Result
	result.Text = string(textData)

	c.logDebug("OCR result length: %d characters", len(result.Text))

	return &result, nil
}

// RecognizeReader 从读取器中读取图像数据并识别文本
func (c *Client) RecognizeReader(reader io.Reader, language string) (*ocr_universal.RespResult, error) {
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

	// 调用RecognizeImage方法
	return c.RecognizeImage(format, imageData, language)
}

// OcrUniversal 是为了向后兼容而保留的函数
// 这里有一个修正，原始函数返回的是空指针和可能的错误
func OcrUniversal(aseAppid, aseAPIKey, aseAPISecret string, serverUrl string, language string, fname string) (*ocr_universal.Response, error) {
	client := New(aseAppid, aseAPIKey, aseAPISecret, serverUrl)

	result, err := client.RecognizeFile(fname, language)
	if err != nil {
		return nil, err
	}

	// 将结果转换为旧格式
	response := &ocr_universal.Response{
		Header: struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			SID     string `json:"sid"`
		}{
			Code:    0,
			Message: "",
			SID:     "",
		},
		Payload: struct {
			Result ocr_universal.RespResult `json:"result"`
		}{
			Result: *result,
		},
	}

	return response, nil
}
