package utils

import (
	"github.com/saintfish/chardet"
)

// encodingCategoryMapping 映射 chardet 检测到的编码名 → 自定义业务分类名
var encodingCategoryMapping = map[string]string{
	"UTF-8":    "UTF8",
	"UTF-16LE": "UNICODE",
	"UTF-16BE": "UNICODE",
	"UTF-32LE": "UNICODE",
	"UTF-32BE": "UNICODE",
	"GB-18030": "GB18030",
	"GBK":      "GBK",
	"GB2312":   "GB2312",
	"Big5":     "BIG5",
}

// DetectEncoding 返回最可能的原始编码名和置信度
func DetectEncoding(input string) (charset string, confidence int, err error) {
	result, err := detect(input)
	if err != nil {
		return "", 0, err
	}
	return result.Charset, result.Confidence, nil
}

// DetectEncodingCategory 返回你的业务分类名和置信度
func DetectV2TTSEncodingCategory(input string) (category string, confidence int, err error) {
	result, err := detect(input)
	if err != nil {
		return "", 0, err
	}

	category, ok := encodingCategoryMapping[result.Charset]
	if !ok {
		category = "未知编码"
	}
	return category, result.Confidence, nil
}

// 内部复用函数，避免重复 detector 创建和转换逻辑
func detect(input string) (*chardet.Result, error) {
	detector := chardet.NewTextDetector()
	return detector.DetectBest([]byte(input))
}
