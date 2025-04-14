package ost

import (
	"encoding/json"
	"fmt"
	"strings"
)

// 解析原始查询响应为用户友好的结果
func FormatAsUserFriendlyResult(raw *RawQueryResponse) (*UserFriendlyRecognitionResult, error) {
	if raw == nil {
		return nil, fmt.Errorf("无效的输入响应")
	}

	result := &UserFriendlyRecognitionResult{
		Code:    raw.Code,
		Message: raw.Message,
	}

	result.Data.TaskID = raw.Data.TaskID
	result.Data.TaskStatus = raw.Data.TaskStatus

	if raw.Data.TaskStatus != StatusCompleted && raw.Data.TaskStatus != StatusCallbackComplete {
		return result, nil
	}

	result.Data.Result.FileLength = raw.Data.Result.FileLength

	// 提取第一组结果
	sentences1, fullText1 := extractSentences(raw.Data.Result.Lattice)
	result.Data.Result.Sentences = sentences1
	result.Data.Result.FullText = fullText1

	// 提取第二组结果
	sentences2, fullText2 := extractSentences(raw.Data.Result.Lattice2)
	result.Data.Result.Sentences2 = sentences2
	result.Data.Result.FullText2 = fullText2

	return result, nil
}

func extractSentences(lattices []RawLattice) ([]Sentence, string) {
	sentences := make([]Sentence, 0, len(lattices))
	var fullText strings.Builder

	for _, lattice := range lattices {
		sentence := Sentence{
			Begin:     lattice.Begin,
			End:       lattice.End,
			SpeakerID: lattice.Spk,
			Words:     make([]Word, 0),
		}

		for _, rt := range lattice.JSON1Best.St.Rt {
			for _, ws := range rt.Ws {
				for _, cw := range ws.Cw {
					word := Word{
						Text:         cw.W,
						Confidence:   cw.Sc,
						WordProperty: cw.Wp,
					}
					sentence.Words = append(sentence.Words, word)
					fullText.WriteString(cw.W)
				}
			}
		}

		sentences = append(sentences, sentence)
	}

	return sentences, fullText.String()
}

// 从识别结果中提取文本
func extractText(response *RawQueryResponse) string {
	if response == nil || len(response.Data.Result.Lattice) == 0 {
		return ""
	}

	var text strings.Builder

	for _, lattice := range response.Data.Result.Lattice {
		for _, rt := range lattice.JSON1Best.St.Rt {
			for _, ws := range rt.Ws {
				for _, cw := range ws.Cw {
					text.WriteString(cw.W)
				}
			}
		}
	}

	return text.String()
}

// 从字节数据解析查询响应
func parseQueryResponse(data []byte) (*RawQueryResponse, error) {
	var response RawQueryResponse
	err := json.Unmarshal(data, &response)
	if err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &response, nil
}

// 检查任务是否完成
func isTaskCompleted(status string) bool {
	return status != StatusInProgress && status != StatusProcessing
}

// 转换讯飞任务状态为描述性文本
func statusToString(status string) string {
	switch status {
	case StatusInProgress:
		return "处理中"
	case StatusProcessing:
		return "正在识别"
	case StatusCompleted:
		return "已完成"
	case StatusCallbackComplete:
		return "回调完成"
	default:
		return "未知状态: " + status
	}
}

// 获取讯飞状态码对应的错误信息
func getErrorMessage(code int) string {
	switch code {
	case 0:
		return "成功"
	case 10101:
		return "参数不合法"
	case 10102:
		return "系统繁忙"
	case 10103:
		return "未开通服务权限"
	case 10104:
		return "未授权或IP不在白名单内"
	case 10105:
		return "授权过期"
	case 10106:
		return "签名错误"
	case 10107:
		return "应用ID不存在"
	case 10108:
		return "语音文件为空或过大"
	case 10109:
		return "语音文件解码失败"
	case 10110:
		return "语音文件无有效数据"
	case 10111:
		return "格式有误"
	default:
		return fmt.Sprintf("未知错误码: %d", code)
	}
}

// 合并选项，用默认值填充
func mergeOptions(options *RecognitionOptions) *RecognitionOptions {
	if options == nil {
		return DefaultOptions()
	}

	defaultOpts := DefaultOptions()

	// 使用默认值填充空值
	if options.Language == "" {
		options.Language = defaultOpts.Language
	}
	if options.Accent == "" {
		options.Accent = defaultOpts.Accent
	}
	if options.Domain == "" {
		options.Domain = defaultOpts.Domain
	}
	if options.Encoding == "" {
		options.Encoding = defaultOpts.Encoding
	}
	if options.MaxRetries <= 0 {
		options.MaxRetries = defaultOpts.MaxRetries
	}
	if options.RetryInterval <= 0 {
		options.RetryInterval = defaultOpts.RetryInterval
	}
	if options.Timeout <= 0 {
		options.Timeout = defaultOpts.Timeout
	}

	return options
}
