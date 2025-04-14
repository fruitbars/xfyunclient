package ost

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// uploader 文件上传器
type uploader struct {
	config     Config
	source     AudioSource
	requestID  string
	cloudID    string
	httpClient *http.Client
}

// newUploader 创建上传器
func newUploader(config Config, source AudioSource) *uploader {
	return &uploader{
		config:     config,
		source:     source,
		requestID:  time.Now().Format("200601021504"),
		cloudID:    "0",
		httpClient: &http.Client{Timeout: time.Duration(config.Timeout) * time.Second},
	}
}

// Upload 上传文件，返回URL
func (u *uploader) Upload() (string, error) {
	// 检查文件大小，决定使用单次上传还是分片上传
	size, err := u.source.Size()
	if err != nil {
		return "", fmt.Errorf("获取文件大小失败: %v", err)
	}

	if size < 30*1024*1024 { // 30MB
		return u.uploadSingle()
	} else {
		return u.uploadChunked()
	}
}

// 生成SHA-256哈希
func (u *uploader) generateSHA256(data string) string {
	hash := sha256.Sum256([]byte(data))
	digest := "SHA-256=" + base64.StdEncoding.EncodeToString(hash[:])
	return digest
}

// 创建授权头
func (u *uploader) assembleAuthHeader(requestURL, contentType, method string, body []byte) http.Header {
	parsedURL, _ := url.Parse(requestURL)
	host := parsedURL.Host
	path := parsedURL.Path

	now := time.Now().UTC()
	date := now.Format(time.RFC1123)
	date = strings.Replace(date, "UTC", "GMT", 1)

	digest := "SHA256=" + u.generateSHA256("")

	signatureOrigin := fmt.Sprintf("host: %s\ndate: %s\n%s %s HTTP/1.1\ndigest: %s",
		host, date, method, path, digest)

	mac := hmac.New(sha256.New, []byte(u.config.APISecret))
	mac.Write([]byte(signatureOrigin))
	signatureSha := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	authorization := fmt.Sprintf(`api_key="%s", algorithm="hmac-sha256", headers="host date request-line digest", signature="%s"`,
		u.config.APIKey, signatureSha)

	headers := http.Header{
		"Host":          {host},
		"Date":          {date},
		"Authorization": {authorization},
		"Digest":        {digest},
		"Content-Type":  {contentType},
	}

	return headers
}

// POST请求
func (u *uploader) postRequest(url string, data []byte, contentType string) ([]byte, error) {
	headers := u.assembleAuthHeader(url, contentType, "POST", data)

	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	for key, values := range headers {
		for _, v := range values {
			req.Header.Add(key, v)
		}
	}

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// 单次上传文件
func (u *uploader) uploadSingle() (string, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// 读取文件内容
	fileContent, err := u.source.Read()
	if err != nil {
		return "", fmt.Errorf("读取文件内容失败: %v", err)
	}

	// 添加文件
	fileField, err := w.CreateFormFile("data", u.source.Name())
	if err != nil {
		return "", fmt.Errorf("创建表单字段失败: %v", err)
	}

	_, err = fileField.Write(fileContent)
	if err != nil {
		return "", fmt.Errorf("写入文件内容失败: %v", err)
	}

	// 添加其他表单字段
	_ = w.WriteField("app_id", u.config.AppID)
	_ = w.WriteField("request_id", u.requestID)

	w.Close()

	url := UploadHost + ApiUpload
	contentType := w.FormDataContentType()

	body, err := u.postRequest(url, b.Bytes(), contentType)
	if err != nil {
		return "", fmt.Errorf("上传请求失败: %v", err)
	}

	var result UploadResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if result.Code != 0 {
		return "", fmt.Errorf("上传失败: %s", result.Message)
	}

	return result.Data.URL, nil
}

// 初始化分块上传
func (u *uploader) initUploadChunked() (string, error) {
	req := InitUploadRequest{
		AppID:     u.config.AppID,
		RequestID: u.requestID,
		CloudID:   u.cloudID,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %v", err)
	}

	url := UploadHost + ApiInit
	contentType := "application/json"

	body, err := u.postRequest(url, jsonData, contentType)
	if err != nil {
		return "", fmt.Errorf("初始化上传请求失败: %v", err)
	}

	var resp InitUploadResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if resp.Code != 0 {
		return "", fmt.Errorf("初始化上传失败: %s", resp.Message)
	}

	return resp.Data.UploadID, nil
}

// 上传单个分块
func (u *uploader) uploadChunk(fileContent []byte, uploadID string, sliceID int) error {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// 添加文件块
	fileField, err := w.CreateFormFile("data", u.source.Name())
	if err != nil {
		return fmt.Errorf("创建表单字段失败: %v", err)
	}

	_, err = fileField.Write(fileContent)
	if err != nil {
		return fmt.Errorf("写入文件内容失败: %v", err)
	}

	// 添加其他表单字段
	_ = w.WriteField("app_id", u.config.AppID)
	_ = w.WriteField("request_id", u.requestID)
	_ = w.WriteField("upload_id", uploadID)
	_ = w.WriteField("slice_id", fmt.Sprintf("%d", sliceID))

	w.Close()

	url := UploadHost + ApiCut
	contentType := w.FormDataContentType()

	// 失败时最多重试3次
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		body, err := u.postRequest(url, b.Bytes(), contentType)
		if err == nil {
			var result map[string]interface{}
			err = json.Unmarshal(body, &result)
			if err != nil {
				continue
			}

			code, ok := result["code"].(float64)
			if ok && code == 0 {
				return nil
			}
		}

		if i < maxRetries-1 {
			time.Sleep(time.Second)
		}
	}

	return fmt.Errorf("分块上传失败，已重试%d次", maxRetries)
}

// 完成分块上传
func (u *uploader) completeChunkedUpload(uploadID string) (string, error) {
	req := UploadCompleteRequest{
		AppID:     u.config.AppID,
		RequestID: u.requestID,
		UploadID:  uploadID,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %v", err)
	}

	url := UploadHost + ApiCutComplete
	contentType := "application/json"

	body, err := u.postRequest(url, jsonData, contentType)
	if err != nil {
		return "", fmt.Errorf("完成上传请求失败: %v", err)
	}

	var result UploadResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if result.Code != 0 {
		return "", fmt.Errorf("完成上传失败: %s", result.Message)
	}

	return result.Data.URL, nil
}

// 分块上传文件
func (u *uploader) uploadChunked() (string, error) {
	// 初始化上传
	uploadID, err := u.initUploadChunked()
	if err != nil {
		return "", err
	}

	// 读取文件内容
	fileContent, err := u.source.Read()
	if err != nil {
		return "", fmt.Errorf("读取文件内容失败: %v", err)
	}

	fileSize := len(fileContent)
	chunks := int(math.Ceil(float64(fileSize) / float64(FilePieceSize)))

	// 上传每个分块
	for i := 1; i <= chunks; i++ {
		start := (i - 1) * FilePieceSize
		end := start + FilePieceSize

		if end > fileSize {
			end = fileSize
		}

		err = u.uploadChunk(fileContent[start:end], uploadID, i)
		if err != nil {
			return "", fmt.Errorf("上传分块%d失败: %v", i, err)
		}
	}

	// 完成上传
	return u.completeChunkedUpload(uploadID)
}

// Read 实现 AudioSource 接口
func (f *FileAudioSource) Read() ([]byte, error) {
	return os.ReadFile(f.FilePath)
}

// Size 实现 AudioSource 接口
func (f *FileAudioSource) Size() (int64, error) {
	info, err := os.Stat(f.FilePath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// Name 实现 AudioSource 接口
func (f *FileAudioSource) Name() string {
	return filepath.Base(f.FilePath)
}

// Read 实现 AudioSource 接口
func (m *MemoryAudioSource) Read() ([]byte, error) {
	return m.Data, nil
}

// Size 实现 AudioSource 接口
func (m *MemoryAudioSource) Size() (int64, error) {
	return int64(len(m.Data)), nil
}

// Name 实现 AudioSource 接口
func (m *MemoryAudioSource) Name() string {
	if m.FileName == "" {
		return "audio_data"
	}
	return m.FileName
}
