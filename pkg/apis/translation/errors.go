// pkg/translation/errors.go

package translation

import "fmt"

// TranslationError 表示翻译服务返回的错误
type TranslationError struct {
	Code    int    // 错误码
	Message string // 错误信息
}

// Error 实现error接口
func (e *TranslationError) Error() string {
	return fmt.Sprintf("translation error: code=%d, message=%s", e.Code, e.Message)
}

// NewTranslationError 创建一个新的翻译错误
func NewTranslationError(code int, message string) *TranslationError {
	return &TranslationError{
		Code:    code,
		Message: message,
	}
}

// IsTranslationError 判断是否为翻译错误
func IsTranslationError(err error) bool {
	_, ok := err.(*TranslationError)
	return ok
}

// GetTranslationErrorCode 获取翻译错误码
func GetTranslationErrorCode(err error) (int, bool) {
	if transErr, ok := err.(*TranslationError); ok {
		return transErr.Code, true
	}
	return 0, false
}

// 预定义错误码常量
const (
	ErrInvalidParameter    = 10001 // 参数无效
	ErrNetworkError        = 10002 // 网络错误
	ErrServerError         = 10003 // 服务器错误
	ErrUnauthorized        = 10004 // 未授权
	ErrQuotaExceeded       = 10005 // 配额超限
	ErrUnsupportedLanguage = 10006 // 不支持的语言
	ErrTextTooLong         = 10007 // 文本过长
	ErrTranslationFailed   = 10008 // 翻译失败
)
