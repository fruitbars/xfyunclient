package v2tts

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/fruitbars/xfyunclient/pkg/ase"
	"io"
	"os"
	"time"
)

// pkg/v2tts/client.go

// TTSClient 定义了语音合成客户端的接口
type TTSClient interface {
	// TextToSpeech 将文本转换为语音并返回音频数据
	TextToSpeech(text string, options ...TTSOption) ([]byte, error)

	// TextToSpeechToFile 将文本转换为语音并保存到文件
	TextToSpeechToFile(text, fileName string, options ...TTSOption) error

	// TextToSpeechToWriter 将文本转换为语音并写入到指定的writer
	TextToSpeechToWriter(text string, writer io.Writer, options ...TTSOption) error
}

// V2TTSClient 实现了TTSClient接口
type V2TTSClient struct {
	appID      string
	serverURL  string
	aseAppID   string
	aseAPIKey  string
	aseSecret  string
	defaultAue string
	defaultVcn string
	defaultTte string
}

// NewV2TTSClient 创建新的V2TTS客户端
func NewV2TTSClient(appID, aseAPIKey, aseSecret, serverURL string) *V2TTSClient {
	return &V2TTSClient{
		appID:      appID,
		serverURL:  serverURL,
		aseAPIKey:  aseAPIKey,
		aseSecret:  aseSecret,
		defaultAue: "raw",     // 默认音频格式
		defaultVcn: "xiaoyan", // 默认发音人
		defaultTte: "utf8",    // 默认文本编码
	}
}

// SetDefaultAue 设置默认音频格式
func (c *V2TTSClient) SetDefaultAue(aue string) {
	c.defaultAue = aue
}

// SetDefaultVoice 设置默认发音人
func (c *V2TTSClient) SetDefaultVoice(vcn string) {
	c.defaultVcn = vcn
}

// SetDefaultTextEncoding 设置默认文本编码
func (c *V2TTSClient) SetDefaultTextEncoding(tte string) {
	c.defaultTte = tte
}

// TextToSpeechWithCallback 使用回调函数处理每一块音频数据
func (c *V2TTSClient) TextToSpeechWithCallback(
	text string,
	audioCallback func(audioChunk []byte) error,
	options ...TTSOption,
) error {
	// 创建基础请求
	req := CreateDefaultV2TTSRequest(c.appID, text, c.defaultAue, c.defaultVcn, c.defaultTte)

	// 应用自定义选项
	ApplyOptions(&req, options...)

	// 创建ASE客户端
	client := ase.NewASEWebsocketClient(
		c.serverURL,
		c.aseAppID,
		c.aseAPIKey,
		c.aseSecret,
		ase.DefaultASEHttpProto,
		ase.DefaultASEAlgorithm,
	)

	// 用于通知处理完成或出错
	doneChan := make(chan struct{})
	errChan := make(chan error, 1)

	// 回调函数处理响应
	callback := func(p []byte) {
		var resp V2TTSResponse
		if err := json.Unmarshal(p, &resp); err != nil {
			errChan <- fmt.Errorf("failed to unmarshal TTS response: %w", err)
			return
		}

		// 检查是否是最后一条消息
		if resp.Code == 0 && resp.Data.Status == 2 {
			// 最后一条消息，处理完成
			close(doneChan)
			return
		}

		if resp.Data.Audio == "" {
			return
		}

		audioBytes, err := base64.StdEncoding.DecodeString(resp.Data.Audio)
		if err != nil {
			errChan <- fmt.Errorf("failed to decode base64 audio: %w", err)
			return
		}

		// 调用用户提供的回调函数处理这块音频数据
		if err := audioCallback(audioBytes); err != nil {
			errChan <- fmt.Errorf("audio callback failed: %w", err)
			return
		}
	}

	// 调用API
	_, err := client.CallASEAPICallBack(req, callback)
	if err != nil {
		return fmt.Errorf("TTS request failed: %w", err)
	}

	// 等待所有数据接收完成或出错
	select {
	case <-doneChan:
		// 所有数据已接收
		return nil
	case err := <-errChan:
		return err
	case <-time.After(30 * time.Second): // 设置超时
		return fmt.Errorf("TTS request timed out")
	}
}

// TextToSpeech 将文本转换为语音，返回完整的音频数据
func (c *V2TTSClient) TextToSpeech(text string, options ...TTSOption) ([]byte, error) {
	var audioBuffers [][]byte

	// 定义回调函数，收集所有音频数据
	collectCallback := func(audioChunk []byte) error {
		audioBuffers = append(audioBuffers, audioChunk)
		return nil
	}

	// 调用带回调的方法
	err := c.TextToSpeechWithCallback(text, collectCallback, options...)
	if err != nil {
		return nil, err
	}

	// 合并所有收集到的音频数据
	totalSize := 0
	for _, buf := range audioBuffers {
		totalSize += len(buf)
	}

	result := make([]byte, 0, totalSize)
	for _, buf := range audioBuffers {
		result = append(result, buf...)
	}

	return result, nil
}

// TextToSpeechToFile 将文本转换为语音并保存到文件
func (c *V2TTSClient) TextToSpeechToFile(text, fileName string, options ...TTSOption) error {
	// 创建或打开文件
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open or create file: %w", err)
	}
	defer file.Close()

	// 调用TextToSpeechToWriter
	return c.TextToSpeechToWriter(text, file, options...)
}

// TextToSpeechToWriter 将文本转换为语音并写入到指定的writer
func (c *V2TTSClient) TextToSpeechToWriter(text string, writer io.Writer, options ...TTSOption) error {
	// 获取音频数据
	audioData, err := c.TextToSpeech(text, options...)
	if err != nil {
		return err
	}

	// 写入到writer
	_, err = writer.Write(audioData)
	if err != nil {
		return fmt.Errorf("failed to write audio data: %w", err)
	}

	return nil
}
