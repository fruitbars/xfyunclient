package ost

// const.go - 常量定义
// ====================================

// 上传相关常量
const (
	// UploadHost 文件上传主机地址
	UploadHost = "http://upload-ost-api.xfyun.cn/file"

	// ApiInit 初始化分片上传接口
	ApiInit = "/mpupload/init"

	// ApiUpload 普通上传接口
	ApiUpload = "/upload"

	// ApiCut 分片上传接口
	ApiCut = "/mpupload/upload"

	// ApiCutComplete 完成分片上传接口
	ApiCutComplete = "/mpupload/complete"

	// ApiCutCancel 取消分片上传接口
	ApiCutCancel = "/mpupload/cancel"

	// FilePieceSize 文件分片大小(5MB)
	FilePieceSize = 5 * 1024 * 1024

	// MaxUploadRetries 上传最大重试次数
	MaxUploadRetries = 3

	// MaxUploadSize 普通上传最大文件大小(30MB)
	MaxUploadSize = 30 * 1024 * 1024
)

// 识别相关常量
const (
	// OstApiHost OST API主机地址
	OstApiHost = "ost-api.xfyun.cn"

	// RequestUriCreate 创建任务接口
	RequestUriCreate = "/v2/ost/pro_create"

	// RequestUriQuery 查询任务接口
	RequestUriQuery = "/v2/ost/query"
)

// 任务状态常量
const (
	// StatusInProgress 任务进行中状态1
	StatusInProgress = "1"

	// StatusProcessing 任务进行中状态2
	StatusProcessing = "2"

	// StatusCompleted 任务已完成状态
	StatusCompleted = "3"

	// StatusCallbackComplete 任务回调完成状态
	StatusCallbackComplete = "4"

	// StatusFailed 任务失败状态
	StatusFailed = "10"

	// StatusUnauthorized 未授权状态
	StatusUnauthorized = "11"

	// StatusInvalidParameter 参数无效状态
	StatusInvalidParameter = "12"
)

// 语言常量
const (
	// LangZhCN 中文普通话
	LangZhCN = "zh_cn"

	// LangEnUS 英语（美国）
	LangEnUS = "en_us"

	// LangJaJP 日语
	LangJaJP = "ja_jp"

	// LangKoKR 韩语
	LangKoKR = "ko_kr"

	// LangRuRU 俄语
	LangRuRU = "ru_ru"

	// LangFrFR 法语
	LangFrFR = "fr_fr"

	// LangEsES 西班牙语
	LangEsES = "es_es"
)

// 方言常量
const (
	// AccentMandarin 普通话
	AccentMandarin = "mandarin"

	// AccentCantonese 粤语
	AccentCantonese = "cantonese"

	// AccentSichuan 四川话
	AccentSichuan = "sichuan"
)

// 领域常量
const (
	// DomainOST 通用领域
	DomainOST = "pro_ost"

	// DomainOSTEd 通用增强版
	DomainOSTEd = "pro_ost_ed"

	// DomainHealthCare 医疗领域
	DomainHealthCare = "healthcare"

	// DomainAuto 汽车领域
	DomainAuto = "auto"

	// DomainFinance 金融领域
	DomainFinance = "finance"
)

// 编码常量
const (
	// EncodingRaw 原始编码
	EncodingRaw = "raw"

	// EncodingLAME LAME编码
	EncodingLAME = "lame"

	// EncodingAMR AMR编码
	EncodingAMR = "amr"
)

// 音频格式常量
const (
	// FormatPCM PCM格式
	FormatPCM = "audio/L16;rate=16000"

	// FormatWAV WAV格式
	FormatWAV = "audio/wav;rate=16000"

	// FormatMP3 MP3格式
	FormatMP3 = "audio/mpeg"
)

// 错误码常量
const (
	// ErrCodeSuccess 成功
	ErrCodeSuccess = 0

	// ErrCodeInvalidParam 参数不合法
	ErrCodeInvalidParam = 10101

	// ErrCodeBusy 系统繁忙
	ErrCodeBusy = 10102

	// ErrCodeNoPermission 未开通服务权限
	ErrCodeNoPermission = 10103

	// ErrCodeUnauthorized 未授权或IP白名单错误
	ErrCodeUnauthorized = 10104

	// ErrCodeExpired 授权过期
	ErrCodeExpired = 10105

	// ErrCodeSignatureError 签名错误
	ErrCodeSignatureError = 10106

	// ErrCodeAppIDNotExist 应用ID不存在
	ErrCodeAppIDNotExist = 10107

	// ErrCodeAudioEmpty 音频文件为空或过大
	ErrCodeAudioEmpty = 10108

	// ErrCodeAudioDecodeFailed 音频文件解码失败
	ErrCodeAudioDecodeFailed = 10109

	// ErrCodeAudioInvalid 音频文件无有效音频数据
	ErrCodeAudioInvalid = 10110
)

// 其他常量
const (
	// DefaultTimeout 默认超时时间(秒)
	DefaultTimeout = 30

	// DefaultRetryInterval 默认重试间隔(秒)
	DefaultRetryInterval = 1

	// DefaultMaxRetries 默认最大重试次数
	DefaultMaxRetries = 100

	// MaxAudioDuration 最大支持音频时长(小时)
	MaxAudioDuration = 8

	// Version 库版本
	Version = "1.0.0"
)
