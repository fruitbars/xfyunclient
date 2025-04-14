package v2tts

// pkg/v2tts/config.go

// TTSConfig 包含TTS服务的配置
type TTSConfig struct {
	// API配置
	AppID     string // 应用ID
	ASEAppID  string // ASE应用ID
	ASEAPIKey string // ASE API密钥
	ASESecret string // ASE密钥
	ServerURL string // 服务器URL

	// 默认TTS参数
	DefaultAue string // 默认音频格式
	DefaultVcn string // 默认发音人
	DefaultTte string // 默认文本编码

	// 高级参数
	Timeout        int  // 超时时间(秒)
	RetryCount     int  // 重试次数
	UseCompression bool // 是否使用压缩
}

// DefaultConfig 返回默认配置
func DefaultConfig() TTSConfig {
	return TTSConfig{
		// API默认配置留空，需要用户提供
		AppID:     "",
		ASEAppID:  "",
		ASEAPIKey: "",
		ASESecret: "",
		ServerURL: "wss://tts-api.xfyun.cn/v2/tts",

		// 默认TTS参数
		DefaultAue: "raw",     // 原始PCM格式
		DefaultVcn: "xiaoyan", // 默认发音人
		DefaultTte: "utf8",    // 默认文本编码

		// 高级参数
		Timeout:        30,    // 30秒超时
		RetryCount:     3,     // 最多重试3次
		UseCompression: false, // 默认不使用压缩
	}
}

// NewClientWithConfig 使用配置创建新的TTS客户端
func NewClientWithConfig(config TTSConfig) *V2TTSClient {
	client := NewV2TTSClient(
		config.AppID,
		config.ASEAPIKey,
		config.ASESecret,
		config.ServerURL,
	)

	// 设置默认参数
	client.SetDefaultAue(config.DefaultAue)
	client.SetDefaultVoice(config.DefaultVcn)
	client.SetDefaultTextEncoding(config.DefaultTte)

	return client
}

// 预定义的语音角色常量
const (
	VoiceXiaoyan   = "xiaoyan"   // 小燕，女声，普通话
	VoiceAimei     = "aimei"     // 艾梅，女声，普通话
	VoiceXiaoyu    = "xiaoyu"    // 小宇，男声，普通话
	VoiceXiaowu    = "xiaowu"    // 小武，男声，普通话
	VoiceXiaomei   = "xiaomei"   // 小梅，女声，粤语
	VoiceXiaolin   = "xiaolin"   // 小林，女声，台湾普通话
	VoiceXiaoqian  = "xiaoqian"  // 小倩，女声，东北话
	VoiceXiaorong  = "xiaorong"  // 小蓉，女声，四川话
	VoiceXiaokun   = "xiaokun"   // 小坤，男声，河南话
	VoiceXiaoqiang = "xiaoqiang" // 小强，男声，湖南话
)

// 音频格式常量
const (
	AudioFormatRaw  = "raw"  // 原始PCM
	AudioFormatWav  = "wav"  // WAV格式
	AudioFormatMp3  = "mp3"  // MP3格式
	AudioFormatAAC  = "aac"  // AAC格式
	AudioFormatSilk = "silk" // SILK格式
)
