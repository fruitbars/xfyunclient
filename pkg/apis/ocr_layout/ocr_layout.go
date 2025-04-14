// pkg/ocr_layout/ocr_layout.go

package ocr_layout

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fruitbars/xfyunclient/pkg/ase"
	"github.com/fruitbars/xfyunclient/pkg/models/ocr_layout"
	"github.com/fruitbars/xfyunclient/pkg/utils"
)

// OCRClient 定义了OCR客户端的接口
type OCRClient interface {
	// ScanFile 分析指定文件中的图像并返回识别的文本
	ScanFile(filePath string) (string, error)

	// ScanImage 分析指定格式和图像数据并返回识别的文本
	ScanImage(format string, imageData []byte) (string, error)

	// ScanReader 从读取器中读取图像数据并返回识别的文本
	ScanReader(reader io.Reader) (string, error)
}

// OCRLayoutClient 实现了OCRClient接口
type OCRLayoutClient struct {
	appID     string
	apiKey    string
	apiSecret string
	serverURL string
	debug     bool
}

// NewOCRLayoutClient 创建一个新的OCR布局识别客户端
func NewOCRLayoutClient(appID, apiKey, apiSecret, serverURL string) *OCRLayoutClient {
	return &OCRLayoutClient{
		appID:     appID,
		apiKey:    apiKey,
		apiSecret: apiSecret,
		serverURL: serverURL,
		debug:     false,
	}
}

// SetDebug 设置是否启用调试模式
func (c *OCRLayoutClient) SetDebug(debug bool) {
	c.debug = debug
}

// logDebug 打印调试信息
func (c *OCRLayoutClient) logDebug(format string, args ...interface{}) {
	if c.debug {
		fmt.Printf("[OCRLayout Debug] "+format+"\n", args...)
	}
}

// ScanFile 分析指定文件中的图像并返回识别的文本
func (c *OCRLayoutClient) ScanFile(filePath string) (string, error) {
	// 验证文件路径
	if filePath == "" {
		return "", errors.New("file path cannot be empty")
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("file does not exist: %s", filePath)
	}

	// 读取图像文件
	format, imageData, err := utils.ReadImageFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read image file: %w", err)
	}

	c.logDebug("Image format: %s, size: %d bytes", format, len(imageData))

	// 调用ScanImage方法
	return c.ScanImage(format, imageData)
}

// ScanImage 分析指定格式和图像数据并返回识别的文本
func (c *OCRLayoutClient) ScanImage(format string, imageData []byte) (string, error) {
	// 验证参数
	if len(imageData) == 0 {
		return "", errors.New("image data cannot be empty")
	}

	// 检查图像格式
	format = strings.ToLower(format)
	if format != "jpg" && format != "jpeg" && format != "png" && format != "bmp" {
		return "", fmt.Errorf("unsupported image format: %s", format)
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

	// 创建OCR布局请求
	req := ocr_layout.NewASEOCRLayoutRequest(c.appID, format, imageData)

	// 调用ASE API
	respData, err := client.CallASEAPIJson(req)
	if err != nil {
		return "", fmt.Errorf("API call failed: %w", err)
	}

	// 验证响应数据
	if respData == nil {
		return "", errors.New("empty response from API")
	}

	c.logDebug("Received %d bytes of response data", len(respData))

	// 解析响应
	var ocrResponse ocr_layout.ASEOCRLayoutResponse
	if err := json.Unmarshal(respData, &ocrResponse); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// 检查响应状态
	if ocrResponse.Header.Code != 0 {
		return "", fmt.Errorf("OCR error: code=%d, message=%s",
			ocrResponse.Header.Code,
			ocrResponse.Header.Message,
		)
	}

	// 解码识别文本
	if ocrResponse.Payload.OcrOutputText.Text == "" {
		return "", nil // 无识别结果但不是错误
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(ocrResponse.Payload.OcrOutputText.Text)
	if err != nil {
		return "", fmt.Errorf("failed to decode text: %w", err)
	}

	recognizedText := string(decodedBytes)
	c.logDebug("OCR result length: %d characters", len(recognizedText))

	return recognizedText, nil
}

// ScanReader 从读取器中读取图像数据并返回识别的文本
func (c *OCRLayoutClient) ScanReader(reader io.Reader) (string, error) {
	// 读取全部数据
	imageData, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read image data: %w", err)
	}

	// 检测图像格式
	format, err := utils.DetectImageFormat(imageData)
	if err != nil {
		return "", fmt.Errorf("failed to detect image format: %w", err)
	}

	// 调用ScanImage方法
	return c.ScanImage(format, imageData)
}

// OCRLayoutFile 是为了向后兼容而保留的函数
// 推荐使用OCRLayoutClient代替
func OCRLayoutFile(aseAppid, aseAPIKey, aseAPISecret string, serverUrl string, fname string) (string, error) {
	client := NewOCRLayoutClient(aseAppid, aseAPIKey, aseAPISecret, serverUrl)
	text, err := client.ScanFile(fname)
	return text, err
}

// OCRLayout 是为了向后兼容而保留的函数
// 推荐使用OCRLayoutClient代替
func OCRLayout(aseAppid, aseAPIKey, aseAPISecret string, serverUrl string, format string, imageData []byte) (string, error) {
	client := NewOCRLayoutClient(aseAppid, aseAPIKey, aseAPISecret, serverUrl)
	text, err := client.ScanImage(format, imageData)
	return text, err
}
