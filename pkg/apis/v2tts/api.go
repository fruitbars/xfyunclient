package v2tts

// pkg/v2tts/api.go
import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/fruitbars/xfyunclient/pkg/ase"
	"os"
)

// 为了兼容性保留原有的函数名和签名，但内部实现使用新的架构

// V2TextToSpeech 是基础版本，带回调处理
// 这是为了兼容旧版API而保留的函数
func V2TextToSpeech(appid, aseAppid, aseAPIKey, aseAPISecret, serverUrl, text, aue, vcn, tte string, cb ase.WSCB) (string, error) {
	req := CreateDefaultV2TTSRequest(appid, text, aue, vcn, tte)
	client := ase.NewASEWebsocketClient(serverUrl, aseAppid, aseAPIKey, aseAPISecret, ase.DefaultASEHttpProto, ase.DefaultASEAlgorithm)
	return client.CallASEAPICallBack(req, cb)
}

// handleAudioMessageToFile 将TTS响应中的base64音频解码并写入文件
func handleAudioMessageToFile(p []byte, file *os.File) error {
	var resp V2TTSResponse
	if err := json.Unmarshal(p, &resp); err != nil {
		return fmt.Errorf("failed to unmarshal TTS response: %w", err)
	}

	if resp.Data.Audio == "" {
		return nil
	}

	audioBytes, err := base64.StdEncoding.DecodeString(resp.Data.Audio)
	if err != nil {
		return fmt.Errorf("failed to decode base64 audio: %w", err)
	}

	if _, err := file.Write(audioBytes); err != nil {
		return fmt.Errorf("failed to write audio to file: %w", err)
	}

	return nil
}

// V2TextToSpeechToFile 提供完整配置并写入文件
// 这是为了兼容旧版API而保留的函数
func V2TextToSpeechToFile(
	aseAppid, aseAPIKey, aseAPISecret, serverUrl, text, aue string,
	sfl int, auf, vcn string, speed, volume, pitch, bgs int,
	tte, reg, rdn, fileName string,
) (string, error) {
	// 使用新的选项模式创建请求
	request := CreateDefaultV2TTSRequest(aseAppid, text, aue, vcn, tte)

	// 手动设置选项，以保持与原API相同的行为
	options := []TTSOption{
		WithStreamMode(sfl == 1),
		WithAudioFormat(auf),
		WithSpeed(speed),
		WithVolume(volume),
		WithPitch(pitch),
		WithBackgroundSound(bgs == 1),
		WithRegionalAccent(reg),
		WithNumberReading(rdn),
	}

	ApplyOptions(&request, options...)

	// 创建ASE客户端
	client := ase.NewASEWebsocketClient(serverUrl, aseAppid, aseAPIKey, aseAPISecret, ase.DefaultASEHttpProto, ase.DefaultASEAlgorithm)

	// 创建或打开文件
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to open or create file: %w", err)
	}
	defer file.Close()

	// 使用回调处理响应并写入文件
	sid, err := client.CallASEAPICallBack(request, func(p []byte) {
		if err := handleAudioMessageToFile(p, file); err != nil {
			// 只能记录错误，因为回调函数无法返回错误
			fmt.Printf("Error handling audio: %v\n", err)
		}
	})

	if err != nil {
		return sid, fmt.Errorf("TTS request failed: %w", err)
	}

	return sid, nil
}

// V2TextToSpeechToFileDefault 使用默认配置的简化版本
// 这是为了兼容旧版API而保留的函数
func V2TextToSpeechToFileDefault(
	aseAppid, aseAPIKey, aseAPISecret, serverUrl, text, aue, vcn, tte, fileName string,
) (string, error) {
	// 创建新的客户端
	client := NewV2TTSClient(aseAppid, aseAPIKey, aseAPISecret, serverUrl)

	// 设置默认值
	client.SetDefaultAue(aue)
	client.SetDefaultVoice(vcn)
	client.SetDefaultTextEncoding(tte)

	// 将文本转换为语音并写入文件
	err := client.TextToSpeechToFile(text, fileName)
	if err != nil {
		return "", err
	}

	// 返回会话ID（这里无法获取真实的sid，因为新接口不直接暴露它）
	return "session_id_not_available_in_new_api", nil
}
