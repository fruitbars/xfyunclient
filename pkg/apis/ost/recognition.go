package ost

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// recognizer 语音识别器
type recognizer struct {
	config     Config
	httpClient *http.Client
}

// newRecognizer 创建识别器
func newRecognizer(config Config) *recognizer {
	return &recognizer{
		config:     config,
		httpClient: &http.Client{Timeout: time.Duration(config.Timeout) * time.Second},
	}
}

// RecognizeAudio 创建识别任务并获取结果
func (r *recognizer) RecognizeAudio(fileURL string, options *RecognitionOptions) (*RawQueryResponse, error) {
	// 合并选项
	opts := mergeOptions(options)

	// 创建任务
	taskID, err := r.createTask(fileURL, opts)
	if err != nil {
		return nil, err
	}

	// 等待结果
	err = r.WaitForResult(taskID, opts)
	if err != nil {
		return nil, err
	}

	// 获取最终结果
	return r.GetResult(taskID)
}

// 生成签名
func (r *recognizer) generateSignature(digest, uri string) string {
	date := time.Now().UTC().Format(time.RFC1123)
	date = strings.Replace(date, "UTC", "GMT", 1)

	signatureStr := fmt.Sprintf("host: %s\ndate: %s\nPOST %s HTTP/1.1\ndigest: %s",
		OstApiHost, date, uri, digest)

	mac := hmac.New(sha256.New, []byte(r.config.APISecret))
	mac.Write([]byte(signatureStr))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return signature
}

// 初始化请求头
func (r *recognizer) initHeader(jsonData []byte, uri string) http.Header {
	date := time.Now().UTC().Format(time.RFC1123)
	date = strings.Replace(date, "UTC", "GMT", 1)

	// 生成摘要
	hash := sha256.Sum256(jsonData)
	digest := "SHA-256=" + base64.StdEncoding.EncodeToString(hash[:])

	// 生成签名
	sign := r.generateSignature(digest, uri)

	// 创建授权头
	authHeader := fmt.Sprintf(`api_key="%s", algorithm="hmac-sha256", headers="host date request-line digest", signature="%s"`,
		r.config.APIKey, sign)

	headers := http.Header{
		"Content-Type":  {"application/json"},
		"Accept":        {"application/json"},
		"Method":        {"POST"},
		"Host":          {OstApiHost},
		"Date":          {date},
		"Digest":        {digest},
		"Authorization": {authHeader},
	}

	return headers
}

// 调用OST API
func (r *recognizer) callOSTAPI(apiURL string, jsonData []byte, uri string) ([]byte, error) {
	headers := r.initHeader(jsonData, uri)

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, err
	}

	for key, values := range headers {
		for _, v := range values {
			req.Header.Add(key, v)
		}
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API请求失败，状态码：%d", resp.StatusCode)
	}

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println(string(respBody))

	return respBody, nil
}

// createTask 创建识别任务
func (r *recognizer) createTask(fileURL string, options *RecognitionOptions) (string, error) {
	// 准备请求体
	req := CreateTaskRequest{
		Common: struct {
			AppID string `json:"app_id"`
		}{
			AppID: r.config.AppID,
		},
		Business: struct {
			Language    string `json:"language"`
			Domain      string `json:"domain"`
			Accent      string `json:"accent"`
			CallbackURL string `json:"callback_url,omitempty"`
			// 其他可选参数
			VsppOn         string `json:"vspp_on,omitempty"`
			SpeakerNum     string `json:"speaker_num,omitempty"`
			OutputType     string `json:"output_type,omitempty"`
			PostprocOn     string `json:"postproc_on,omitempty"`
			Pd             string `json:"pd,omitempty"`
			Duration       string `json:"duration,omitempty"`
			EnableSubtitle string `json:"enable_subtitle,omitempty"`
			Smoothproc     string `json:"smoothproc,omitempty"`
			Colloqproc     string `json:"colloqproc,omitempty"`
			LanguageType   string `json:"language_type,omitempty"`
			Vto            string `json:"vto,omitempty"`
			Dhw            string `json:"dhw,omitempty"`
		}{
			Language:    options.Language,
			Accent:      options.Accent,
			Domain:      options.Domain,
			CallbackURL: options.CallbackURL,
		},
		Data: struct {
			AudioURL  string `json:"audio_url"`
			AudioSrc  string `json:"audio_src"`
			AudioSize string `json:"audio_size,omitempty"`
			Format    string `json:"format,omitempty"`
			Encoding  string `json:"encoding"`
		}{
			AudioSrc: "http",
			AudioURL: fileURL,
			Encoding: options.Encoding,
			Format:   options.Format,
		},
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %v", err)
	}

	// 创建URL和URI
	uri := RequestUriCreate
	apiURL := "https://" + OstApiHost + uri

	// 发出API调用
	respBody, err := r.callOSTAPI(apiURL, jsonData, uri)
	if err != nil {
		return "", err
	}

	log.Println(string(respBody))

	var result CreateTaskResponse
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if result.Code != 0 {
		return "", fmt.Errorf("创建任务失败: %s", result.Message)
	}

	return result.Data.TaskID, nil
}

// GetResult 获取识别结果
func (r *recognizer) GetResult(taskID string) (*RawQueryResponse, error) {
	// 准备请求体
	req := QueryTaskRequest{
		Common: struct {
			AppID string `json:"app_id"`
		}{
			AppID: r.config.AppID,
		},
		Business: struct {
			TaskID string `json:"task_id"`
		}{
			TaskID: taskID,
		},
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	// 创建URL和URI
	uri := RequestUriQuery
	apiURL := "https://" + OstApiHost + uri

	// 发出API调用
	respBody, err := r.callOSTAPI(apiURL, jsonData, uri)
	if err != nil {
		return nil, err
	}

	// 解析原始响应
	return parseQueryResponse(respBody)

}

// queryTask 查询任务状态
func (r *recognizer) queryTask(taskID string) (*RawQueryResponse, error) {
	// 准备请求体
	req := QueryTaskRequest{
		Common: struct {
			AppID string `json:"app_id"`
		}{
			AppID: r.config.AppID,
		},
		Business: struct {
			TaskID string `json:"task_id"`
		}{
			TaskID: taskID,
		},
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	// 创建URL和URI
	uri := RequestUriQuery
	apiURL := "https://" + OstApiHost + uri

	// 发出API调用
	respBody, err := r.callOSTAPI(apiURL, jsonData, uri)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println(string(respBody))

	return parseQueryResponse(respBody)
}

// WaitForResult 等待识别结果
func (r *recognizer) WaitForResult(taskID string, options *RecognitionOptions) error {
	opts := mergeOptions(options)
	maxRetries := opts.MaxRetries
	retryInterval := time.Duration(opts.RetryInterval) * time.Second

	for i := 0; i < maxRetries; i++ {
		result, err := r.queryTask(taskID)
		if err != nil {
			log.Println(err)
			return fmt.Errorf("查询任务状态失败: %v", err)
		}

		status := result.Data.TaskStatus

		// 调用状态回调
		if opts.StatusCallback != nil {
			//parsed, _ := parseResult(result)
			opts.StatusCallback(status, result)
		}

		// 提取文本并调用结果回调
		text := extractText(result)
		if text != "" && opts.ResultCallback != nil {
			//parsed, _ := parseResult(result)
			opts.ResultCallback(text, result)
		}

		// 检查任务是否完成
		if isTaskCompleted(status) {
			return nil
		}

		// 等待下一次轮询
		time.Sleep(retryInterval)
	}

	return fmt.Errorf("达到最大重试次数(%d)，任务仍未完成", maxRetries)
}
