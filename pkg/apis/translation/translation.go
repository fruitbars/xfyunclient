// pkg/translation/translation.go

package translation

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fruitbars/xfyunclient/pkg/ase"
	"github.com/fruitbars/xfyunclient/pkg/models/translation_v1_its"
)

// TranslationClient 定义了翻译客户端的接口
type TranslationClient interface {
	// Translate 将文本从源语言翻译到目标语言
	Translate(text, fromLang, toLang string) (string, error)

	// TranslateAuto 自动检测源语言并翻译到目标语言
	TranslateAuto(text, toLang string) (string, error)

	// BatchTranslate 批量翻译多个文本
	BatchTranslate(texts []string, fromLang, toLang string) ([]string, error)
}

// Client 实现了TranslationClient接口
type Client struct {
	appID     string
	apiKey    string
	apiSecret string
	serverURL string
	debug     bool
}

// New 创建一个新的翻译客户端
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
		fmt.Printf("[Translation Debug] "+format+"\n", args...)
	}
}

// Translate 将文本从源语言翻译到目标语言
func (c *Client) Translate(text, fromLang, toLang string) (string, error) {
	// 验证参数
	if text == "" {
		return "", errors.New("text cannot be empty")
	}

	if fromLang == "" {
		return "", errors.New("source language cannot be empty")
	}

	if toLang == "" {
		return "", errors.New("target language cannot be empty")
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
	c.logDebug("Translating from %s to %s", fromLang, toLang)

	// 创建翻译请求
	req := translation_v1_its.NewASETranslationRequest(c.appID, fromLang, toLang, text)

	// 调用ASE API
	response, err := client.CallASEAPIJson(req)
	if err != nil {
		return "", fmt.Errorf("API call failed: %w", err)
	}

	// 验证响应数据
	if response == nil {
		return "", errors.New("empty response from API")
	}

	c.logDebug("Received %d bytes of response data", len(response))

	// 解析响应
	var transResponse translation_v1_its.ASETranslationResult
	if err := json.Unmarshal(response, &transResponse); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// 检查响应状态
	if transResponse.Header.Code != 0 {
		return "", fmt.Errorf("translation error: code=%d, message=%s",
			transResponse.Header.Code,
			transResponse.Header.Message,
		)
	}

	// 解码翻译结果
	if transResponse.Payload.Result.Text == "" {
		return "", errors.New("empty translation result")
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(transResponse.Payload.Result.Text)
	if err != nil {
		return "", fmt.Errorf("failed to decode result: %w", err)
	}

	// 解析翻译引擎结果
	var engResult translation_v1_its.ASETranslationEngineResult
	if err := json.Unmarshal(decodedBytes, &engResult); err != nil {
		return "", fmt.Errorf("failed to parse engine result: %w", err)
	}

	c.logDebug("Translation successful")

	return engResult.TransResult.Dst, nil
}

// TranslateAuto 自动检测源语言并翻译到目标语言
func (c *Client) TranslateAuto(text, toLang string) (string, error) {
	return c.Translate(text, "auto", toLang)
}

// BatchTranslate 批量翻译多个文本
func (c *Client) BatchTranslate(texts []string, fromLang, toLang string) ([]string, error) {
	if len(texts) == 0 {
		return nil, errors.New("texts array cannot be empty")
	}

	results := make([]string, 0, len(texts))

	for _, text := range texts {
		result, err := c.Translate(text, fromLang, toLang)
		if err != nil {
			return results, fmt.Errorf("failed to translate text '%s': %w", truncateText(text, 30), err)
		}
		results = append(results, result)
	}

	return results, nil
}

// truncateText 截断文本以便日志输出
func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}

// 为了兼容旧代码而保留的函数
func Translate(aseAppid, aseAPIKey, aseAPISecret string, serverUrl string, fromLang, toLang, text string) (string, error) {
	client := New(aseAppid, aseAPIKey, aseAPISecret, serverUrl)
	return client.Translate(text, fromLang, toLang)
}
