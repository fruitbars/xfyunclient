package v2tts

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/fruitbars/xfyunclient/pkg/ase"
	"io"
	"os"
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

// TextToSpeech 将文本转换为语音并返回音频数据
func (c *V2TTSClient) TextToSpeech(text string, options ...TTSOption) ([]byte, error) {
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

	// 用于存储音频数据的通道
	audioChan := make(chan []byte, 1)
	errChan := make(chan error, 1)

	// 回调函数处理响应
	callback := func(p []byte) {
		var resp V2TTSResponse
		if err := json.Unmarshal(p, &resp); err != nil {
			errChan <- fmt.Errorf("failed to unmarshal TTS response: %w", err)
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

		audioChan <- audioBytes
	}

	// 调用API
	_, err := client.CallASEAPICallBack(req, callback)
	if err != nil {
		return nil, fmt.Errorf("TTS request failed: %w", err)
	}

	// 等待结果
	select {
	case audioData := <-audioChan:
		return audioData, nil
	case err := <-errChan:
		return nil, err
	}
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
