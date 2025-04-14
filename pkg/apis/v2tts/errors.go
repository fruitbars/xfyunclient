package v2tts

// pkg/v2tts/errors.go

import "fmt"

// TTSError 表示TTS服务返回的错误
type TTSError struct {
	Code    int    // 错误码
	Message string // 错误信息
	Sid     string // 会话ID
}

// Error 实现error接口
func (e *TTSError) Error() string {
	if e.Sid != "" {
		return fmt.Sprintf("TTS error (sid: %s): code=%d, message=%s", e.Sid, e.Code, e.Message)
	}
	return fmt.Sprintf("TTS error: code=%d, message=%s", e.Code, e.Message)
}

// NewTTSError 创建一个新的TTS错误
func NewTTSError(code int, message, sid string) *TTSError {
	return &TTSError{
		Code:    code,
		Message: message,
		Sid:     sid,
	}
}

// 预定义错误码常量
const (
	ErrCodeInvalidParameter  = 10001 // 参数无效
	ErrCodeNetworkError      = 10002 // 网络错误
	ErrCodeServerError       = 10003 // 服务器错误
	ErrCodeUnauthorized      = 10004 // 未授权
	ErrCodeQuotaExceeded     = 10005 // 配额超限
	ErrCodeInvalidAppID      = 10006 // 无效的AppID
	ErrCodeInvalidSignature  = 10007 // 无效的签名
	ErrCodeInvalidEncoding   = 10008 // 无效的编码
	ErrCodeTextTooLong       = 10009 // 文本过长
	ErrCodeUnsupportedFormat = 10010 // 不支持的格式
)

// 错误代码映射到友好错误信息
var errorMessages = map[int]string{
	ErrCodeInvalidParameter:  "无效的参数",
	ErrCodeNetworkError:      "网络连接错误",
	ErrCodeServerError:       "服务器内部错误",
	ErrCodeUnauthorized:      "未授权，请检查API密钥",
	ErrCodeQuotaExceeded:     "API调用配额已超限",
	ErrCodeInvalidAppID:      "无效的应用ID",
	ErrCodeInvalidSignature:  "无效的签名",
	ErrCodeInvalidEncoding:   "无效的编码格式",
	ErrCodeTextTooLong:       "文本内容过长",
	ErrCodeUnsupportedFormat: "不支持的音频格式",
}

// GetErrorMessage 获取错误码对应的错误信息
func GetErrorMessage(code int) string {
	if msg, ok := errorMessages[code]; ok {
		return msg
	}
	return "未知错误"
}

// IsNetworkError 判断是否为网络错误
func IsNetworkError(err error) bool {
	if ttsErr, ok := err.(*TTSError); ok {
		return ttsErr.Code == ErrCodeNetworkError
	}
	return false
}

// IsServerError 判断是否为服务器错误
func IsServerError(err error) bool {
	if ttsErr, ok := err.(*TTSError); ok {
		return ttsErr.Code == ErrCodeServerError
	}
	return false
}

// IsAuthError 判断是否为认证错误
func IsAuthError(err error) bool {
	if ttsErr, ok := err.(*TTSError); ok {
		return ttsErr.Code == ErrCodeUnauthorized ||
			ttsErr.Code == ErrCodeInvalidAppID ||
			ttsErr.Code == ErrCodeInvalidSignature
	}
	return false
}
