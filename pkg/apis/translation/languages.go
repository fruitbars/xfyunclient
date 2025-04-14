// pkg/translation/languages.go

package translation

// Language 定义了支持的语言代码
type Language string

// 支持的语言常量
const (
	LangAuto Language = "auto" // 自动检测
	LangZH   Language = "zh"   // 中文
	LangEN   Language = "en"   // 英文
	LangJA   Language = "ja"   // 日语
	LangKO   Language = "ko"   // 韩语
	LangFR   Language = "fr"   // 法语
	LangES   Language = "es"   // 西班牙语
	LangIT   Language = "it"   // 意大利语
	LangDE   Language = "de"   // 德语
	LangRU   Language = "ru"   // 俄语
	LangPT   Language = "pt"   // 葡萄牙语
	LangVI   Language = "vi"   // 越南语
	LangID   Language = "id"   // 印尼语
	LangTH   Language = "th"   // 泰语
	LangAR   Language = "ar"   // 阿拉伯语
)

// IsValidLanguage 检查语言代码是否有效
func IsValidLanguage(lang string) bool {
	switch Language(lang) {
	case LangAuto, LangZH, LangEN, LangJA, LangKO, LangFR, LangES, LangIT, LangDE, LangRU, LangPT, LangVI, LangID, LangTH, LangAR:
		return true
	default:
		return false
	}
}

// GetLanguageName 获取语言的名称
func GetLanguageName(lang string) string {
	switch Language(lang) {
	case LangAuto:
		return "Auto Detect"
	case LangZH:
		return "Chinese"
	case LangEN:
		return "English"
	case LangJA:
		return "Japanese"
	case LangKO:
		return "Korean"
	case LangFR:
		return "French"
	case LangES:
		return "Spanish"
	case LangIT:
		return "Italian"
	case LangDE:
		return "German"
	case LangRU:
		return "Russian"
	case LangPT:
		return "Portuguese"
	case LangVI:
		return "Vietnamese"
	case LangID:
		return "Indonesian"
	case LangTH:
		return "Thai"
	case LangAR:
		return "Arabic"
	default:
		return "Unknown"
	}
}

// SupportedLanguages 返回所有支持的语言代码
func SupportedLanguages() []Language {
	return []Language{
		LangZH, LangEN, LangJA, LangKO, LangFR, LangES,
		LangIT, LangDE, LangRU, LangPT, LangVI, LangID,
		LangTH, LangAR,
	}
}

// SupportedLanguagePairs 返回所有支持的语言对
func SupportedLanguagePairs() [][2]Language {
	// 这里仅列出部分常用语言对
	// 实际支持的语言对可能根据服务提供商有所不同
	return [][2]Language{
		{LangZH, LangEN},
		{LangEN, LangZH},
		{LangZH, LangJA},
		{LangJA, LangZH},
		{LangZH, LangKO},
		{LangKO, LangZH},
		{LangEN, LangJA},
		{LangJA, LangEN},
		{LangEN, LangES},
		{LangES, LangEN},
		// 添加更多支持的语言对...
	}
}
