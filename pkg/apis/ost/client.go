package ost

// client.go - 主客户端接口
// ====================================

import (
	"fmt"
	"log"
	"time"
)

// Client 是讯飞语音识别客户端
type Client struct {
	config Config
}

// NewClient 创建一个新的语音识别客户端
func NewClient(appID, apiKey, apiSecret string) *Client {
	return &Client{
		config: Config{
			AppID:     appID,
			APIKey:    apiKey,
			APISecret: apiSecret,
			Timeout:   30,
		},
	}
}

// WithTimeout 设置超时时间
func (c *Client) WithTimeout(timeout int) *Client {
	c.config.Timeout = timeout
	return c
}

// WithHTTPClient 设置自定义HTTP客户端
func (c *Client) WithHTTPClient(client interface{}) *Client {
	c.config.HTTPClient = client
	return c
}

// RecognizeFile 识别文件
func (c *Client) RecognizeFile(filePath string, options *RecognitionOptions) (*RawQueryResponse, error) {
	source := &FileAudioSource{FilePath: filePath}
	return c.RecognizeAudio(source, options)
}

// RecognizeBytes 识别内存中的音频数据
func (c *Client) RecognizeBytes(data []byte, fileName string, options *RecognitionOptions) (*RawQueryResponse, error) {
	source := &MemoryAudioSource{Data: data, FileName: fileName}
	return c.RecognizeAudio(source, options)
}

// RecognizeAudio 识别音频
func (c *Client) RecognizeAudio(source AudioSource, options *RecognitionOptions) (*RawQueryResponse, error) {
	// 合并默认选项
	if options == nil {
		options = DefaultOptions()
	}

	// 创建上传器
	uploader := newUploader(c.config, source)

	// 上传文件
	fileURL, err := uploader.Upload()
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("文件上传失败: %v", err)
	}

	// 创建识别任务
	recognizer := newRecognizer(c.config)
	return recognizer.RecognizeAudio(fileURL, options)
}

// RecognizeFileAsync 异步识别文件
func (c *Client) RecognizeFileAsync(filePath string, options *RecognitionOptions,
	resultCallback func(*RawQueryResponse, error)) {

	// 启动goroutine异步处理
	go func() {
		result, err := c.RecognizeFile(filePath, options)
		resultCallback(result, err)
	}()
}

// RecognizeBytesAsync 异步识别内存中的音频数据
func (c *Client) RecognizeBytesAsync(data []byte, fileName string, options *RecognitionOptions,
	resultCallback func(*RawQueryResponse, error)) {

	// 启动goroutine异步处理
	go func() {
		result, err := c.RecognizeBytes(data, fileName, options)
		resultCallback(result, err)
	}()
}

// GetResult 获取识别结果
func (c *Client) GetResult(taskID string) (*RawQueryResponse, error) {
	recognizer := newRecognizer(c.config)
	return recognizer.GetResult(taskID)
}

// GetResultWithCallback 获取结果并使用回调
func (c *Client) GetResultWithCallback(taskID string, statusCallback StatusCallback, resultCallback ResultCallback, options *RecognitionOptions) error {
	if options == nil {
		options = DefaultOptions()
	}

	opts := *options
	opts.StatusCallback = statusCallback
	opts.ResultCallback = resultCallback

	recognizer := newRecognizer(c.config)
	return recognizer.WaitForResult(taskID, &opts)
}

// GetResultWithTimeout 带超时获取结果
func (c *Client) GetResultWithTimeout(taskID string, timeout time.Duration) (*RawQueryResponse, error) {
	recognizer := newRecognizer(c.config)

	// 创建带超时的options
	options := DefaultOptions()
	options.Timeout = int(timeout.Seconds())

	// 创建结果通道
	resultChan := make(chan *RawQueryResponse, 1)
	errorChan := make(chan error, 1)

	// 创建完成通道
	doneChan := make(chan struct{})

	// 启动goroutine来等待结果
	go func() {
		err := recognizer.WaitForResult(taskID, options)
		if err != nil {
			errorChan <- err
			close(doneChan)
			return
		}

		result, err := recognizer.GetResult(taskID)
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- result
		}
		close(doneChan)
	}()

	// 设置超时
	select {
	case <-doneChan:
		// 检查是否有错误
		select {
		case err := <-errorChan:
			return nil, err
		default:
			// 没有错误，获取结果
			select {
			case result := <-resultChan:
				return result, nil
			default:
				return nil, fmt.Errorf("未知错误，没有返回结果")
			}
		}
	case <-time.After(timeout):
		return nil, fmt.Errorf("获取结果超时")
	}
}

// CreateTask 创建语音识别任务但不等待结果
func (c *Client) CreateTask(fileURL string, options *RecognitionOptions) (string, error) {
	if options == nil {
		options = DefaultOptions()
	}

	recognizer := newRecognizer(c.config)
	return recognizer.createTask(fileURL, options)
}
