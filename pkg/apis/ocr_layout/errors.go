// pkg/ocr_layout/errors.go

package ocr_layout

import "fmt"

// OCRError 表示OCR服务返回的错误
type OCRError struct {
	Code    int    // 错误码
	Message string // 错误信息
}

// Error 实现error接口
func (e *OCRError) Error() string {
	return fmt.Sprintf("OCR error: code=%d, message=%s", e.Code, e.Message)
}

// NewOCRError 创建一个新的OCR错误
func NewOCRError(code int, message string) *OCRError {
	return &OCRError{
		Code:    code,
		Message: message,
	}
}

// IsOCRError 判断是否为OCR错误
func IsOCRError(err error) bool {
	_, ok := err.(*OCRError)
	return ok
}

// GetOCRErrorCode 获取OCR错误码
func GetOCRErrorCode(err error) (int, bool) {
	if ocrErr, ok := err.(*OCRError); ok {
		return ocrErr.Code, true
	}
	return 0, false
}

// 预定义错误码常量
const (
	ErrInvalidParameter = 10001 // 参数无效
	ErrNetworkError     = 10002 // 网络错误
	ErrServerError      = 10003 // 服务器错误
	ErrUnauthorized     = 10004 // 未授权
	ErrQuotaExceeded    = 10005 // 配额超限
	ErrUnsupportedImage = 10006 // 不支持的图像格式
	ErrImageTooLarge    = 10007 // 图像过大
	ErrNoTextDetected   = 10008 // 未检测到文本
)
