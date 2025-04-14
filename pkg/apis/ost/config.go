package ost

// config.go - 配置相关
// ====================================

// Config 配置结构
type Config struct {
	AppID      string
	APIKey     string
	APISecret  string
	Timeout    int // 超时时间（秒）
	HTTPClient interface{}
}

// DefaultOptions 返回默认选项
func DefaultOptions() *RecognitionOptions {
	return &RecognitionOptions{
		Language:      "zh_cn",
		Accent:        "mandarin",
		Domain:        "pro_ost_ed",
		Encoding:      "raw",
		MaxRetries:    100,
		RetryInterval: 1,
		Timeout:       300,
	}
}

// Clone 克隆一个配置
func (c *Config) Clone() Config {
	return Config{
		AppID:      c.AppID,
		APIKey:     c.APIKey,
		APISecret:  c.APISecret,
		Timeout:    c.Timeout,
		HTTPClient: c.HTTPClient,
	}
}

// WithTimeout 设置超时时间并返回新的配置
func (c Config) WithTimeout(timeout int) Config {
	c.Timeout = timeout
	return c
}

// CreateOptions 创建识别选项
func CreateOptions() *RecognitionOptions {
	return DefaultOptions()
}

// WithLanguage 设置语言
func (o *RecognitionOptions) WithLanguage(language string) *RecognitionOptions {
	o.Language = language
	return o
}

// WithAccent 设置方言
func (o *RecognitionOptions) WithAccent(accent string) *RecognitionOptions {
	o.Accent = accent
	return o
}

// WithDomain 设置领域
func (o *RecognitionOptions) WithDomain(domain string) *RecognitionOptions {
	o.Domain = domain
	return o
}

// WithEncoding 设置编码
func (o *RecognitionOptions) WithEncoding(encoding string) *RecognitionOptions {
	o.Encoding = encoding
	return o
}

// WithFormat 设置格式
func (o *RecognitionOptions) WithFormat(format string) *RecognitionOptions {
	o.Format = format
	return o
}

// WithCallbackURL 设置回调URL
func (o *RecognitionOptions) WithCallbackURL(callbackURL string) *RecognitionOptions {
	o.CallbackURL = callbackURL
	return o
}

// WithStatusCallback 设置状态回调
func (o *RecognitionOptions) WithStatusCallback(callback StatusCallback) *RecognitionOptions {
	o.StatusCallback = callback
	return o
}

// WithResultCallback 设置结果回调
func (o *RecognitionOptions) WithResultCallback(callback ResultCallback) *RecognitionOptions {
	o.ResultCallback = callback
	return o
}

// WithTimeout 设置超时时间
func (o *RecognitionOptions) WithTimeout(timeout int) *RecognitionOptions {
	o.Timeout = timeout
	return o
}

// WithMaxRetries 设置最大重试次数
func (o *RecognitionOptions) WithMaxRetries(maxRetries int) *RecognitionOptions {
	o.MaxRetries = maxRetries
	return o
}

// WithRetryInterval 设置重试间隔
func (o *RecognitionOptions) WithRetryInterval(retryInterval int) *RecognitionOptions {
	o.RetryInterval = retryInterval
	return o
}
