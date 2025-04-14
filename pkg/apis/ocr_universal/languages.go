// pkg/ocr_universal/languages.go

package ocr_universal

// Language 定义了支持的OCR识别语言
type Language string

// 支持的语言常量
const (
	LangChineseEnglish Language = "ch_en_public_cloud" // 中英文
	LangChinese        Language = "ch_public_cloud"    // 中文
	LangEnglish        Language = "en_public_cloud"    // 英文
	LangJapanese       Language = "ja_public_cloud"    // 日文
	LangKorean         Language = "ko_public_cloud"    // 韩文
	LangFrench         Language = "fr_public_cloud"    // 法文
	LangSpanish        Language = "es_public_cloud"    // 西班牙文
	LangPortuguese     Language = "pt_public_cloud"    // 葡萄牙文
	LangGerman         Language = "de_public_cloud"    // 德文
	LangItalian        Language = "it_public_cloud"    // 意大利文
	LangRussian        Language = "ru_public_cloud"    // 俄文
)

// IsValidLanguage 检查语言代码是否有效
func IsValidLanguage(lang string) bool {
	switch Language(lang) {
	case LangChineseEnglish, LangChinese, LangEnglish, LangJapanese, LangKorean,
		LangFrench, LangSpanish, LangPortuguese, LangGerman, LangItalian, LangRussian:
		return true
	default:
		return false
	}
}

// GetLanguageName 获取语言的友好名称
func GetLanguageName(lang string) string {
	switch Language(lang) {
	case LangChineseEnglish:
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
