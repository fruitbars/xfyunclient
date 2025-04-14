// pkg/image_block/languages.go

package image_block

// Language 定义了支持的OCR识别语言
type Language string

// 支持的语言常量
const (
	LangChineseEnglish Language = "ch_en" // 中英文
	LangChinese        Language = "ch"    // 中文
	LangEnglish        Language = "en"    // 英文
	LangJapanese       Language = "ja"    // 日文
	LangKorean         Language = "ko"    // 韩文
	LangFrench         Language = "fr"    // 法文
	LangSpanish        Language = "es"    // 西班牙文
	LangPortuguese     Language = "pt"    // 葡萄牙文
	LangGerman         Language = "de"    // 德文
	LangItalian        Language = "it"    // 意大利文
	LangRussian        Language = "ru"    // 俄文

	// 通用OCR 2024版语言代码
	LangUniversalChineseEnglish Language = "ch_en_public_cloud" // 通用OCR 2024中英文
)

// IsValidLanguage 检查语言代码是否有效
func IsValidLanguage(lang string) bool {
	switch Language(lang) {
	case LangChineseEnglish, LangChinese, LangEnglish, LangJapanese, LangKorean,
		LangFrench, LangSpanish, LangPortuguese, LangGerman, LangItalian, LangRussian,
		LangUniversalChineseEnglish:
		return true
	default:
		return false
	}
}

// GetLanguageName 获取语言的友好名称
func GetLanguageName(lang string) string {
	switch Language(lang) {
	case LangChineseEnglish, LangUniversalChineseEnglish:
		return "Chinese and English"
	case LangChinese:
		return "Chinese"
	case LangEnglish:
		return "English"
	case LangJapanese:
		return "Japanese"
	case LangKorean:
		return "Korean"
	case LangFrench:
		return "French"
	case LangSpanish:
		return "Spanish"
	case LangPortuguese:
		return "Portuguese"
	case LangGerman:
		return "German"
	case LangItalian:
		return "Italian"
	case LangRussian:
		return "Russian"
	default:
		return "Unknown"
	}
}

// SupportedLanguages 返回所有支持的语言代码
func SupportedLanguages() []Language {
	return []Language{
		LangChineseEnglish, LangChinese, LangEnglish, LangJapanese, LangKorean,
		LangFrench, LangSpanish, LangPortuguese, LangGerman, LangItalian, LangRussian,
	}
}

// SupportedUniversalLanguages 返回支持的通用OCR 2024语言代码
func SupportedUniversalLanguages() []Language {
	return []Language{
		LangUniversalChineseEnglish,
		// 添加其他支持的通用OCR 2024语言...
	}
}

// ConvertToUniversalLanguage 将标准语言代码转换为通用OCR 2024语言代码
func ConvertToUniversalLanguage(lang Language) Language {
	switch lang {
	case LangChineseEnglish:
		return LangUniversalChineseEnglish
	// 添加其他语言映射...
	default:
		return lang
	}
}
