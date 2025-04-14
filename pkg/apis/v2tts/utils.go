package v2tts

import (
	"bytes"
	"encoding/base64"
	"io"
	"os"
)

// DecodeAudioFromResponse 从TTS响应中解码音频数据
func DecodeAudioFromResponse(resp V2TTSResponse) ([]byte, error) {
	if resp.Data.Audio == "" {
		return nil, nil
	}

	return base64.StdEncoding.DecodeString(resp.Data.Audio)
}

// SplitTextByLength 按长度分割文本，避免超出TTS服务的长度限制
// 尽量在标点符号处分割
func SplitTextByLength(text string, maxLength int) []string {
	if len(text) <= maxLength {
		return []string{text}
	}

	var result []string
	var current string

	// 标点符号集合
	punctuations := []rune{'.', '。', '!', '！', '?', '？', ',', '，', ';', '；', ':', '：', '\n'}

	runes := []rune(text)
	lastPunctuationIndex := 0

	for i := 0; i < len(runes); i++ {
		current += string(runes[i])

		// 检查是否为标点符号
		isPunctuation := false
		for _, p := range punctuations {
			if runes[i] == p {
				isPunctuation = true
				lastPunctuationIndex = i
				break
			}
		}

		// 如果达到最大长度或是最后一个字符
		if len([]rune(current)) >= maxLength || i == len(runes)-1 {
			// 如果有标点符号，在标点符号处分割
			if lastPunctuationIndex > 0 && lastPunctuationIndex < i {
				splitPos := lastPunctuationIndex + 1
				result = append(result, string(runes[:splitPos]))

				// 剩余部分递归处理
				remainingText := string(runes[splitPos:])
				subResults := SplitTextByLength(remainingText, maxLength)
				result = append(result, subResults...)

				return result
			}

			// 如果没有合适的标点符号，直接按长度分割
			result = append(result, current)

			if i < len(runes)-1 {
				remainingText := string(runes[i+1:])
				subResults := SplitTextByLength(remainingText, maxLength)
				result = append(result, subResults...)
			}

			return result
		}

		// 更新最后标点符号位置
		if isPunctuation {
			lastPunctuationIndex = i
		}
	}

	return result
}

// ConcatAudioFiles 合并多个音频文件
func ConcatAudioFiles(outputFile string, inputFiles ...string) error {
	outFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	for _, file := range inputFiles {
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		_, err = outFile.Write(data)
		if err != nil {
			return err
		}
	}

	return nil
}

// SaveAudioToFile 保存音频数据到文件
func SaveAudioToFile(audioData []byte, filePath string) error {
	return os.WriteFile(filePath, audioData, 0644)
}

// AudioBufferToFile 将音频缓冲区写入文件
func AudioBufferToFile(buf *bytes.Buffer, filePath string) error {
	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, buf)
	return err
}

// GetTextLength 获取文本长度（字符数）
func GetTextLength(text string) int {
	return len([]rune(text))
}

// DetectTextEncoding 检测文本编码
// 这里简化处理，实际应用中可能需要更复杂的检测逻辑
func DetectTextEncoding(text string) string {
	// 检查是否包含中文字符
	for _, r := range text {
		if r > 0x4E00 && r < 0x9FFF {
			return "utf8" // 中文使用UTF-8
		}
	}
	return "gb2312" // 默认使用GB2312
}

// ValidateText 验证文本是否符合TTS要求
func ValidateText(text string) bool {
	// 检查文本长度
	if len([]rune(text)) == 0 {
		return false
	}

	// 最大支持2000个汉字
	if len([]rune(text)) > 2000 {
		return false
	}

	return true
}

// CleanText 清理文本，去除不必要的字符
func CleanText(text string) string {
	var result bytes.Buffer

	for _, r := range text {
		// 过滤掉一些控制字符
		if r >= 32 && r != 127 {
			result.WriteRune(r)
		}
	}

	return result.String()
}
