package v2tts

// pkg/v2tts/options.go

// TTSOption 用于配置TTSRequest的函数类型
type TTSOption func(*V2TTSRequest)

// WithAue 设置音频编码格式
func WithAue(aue string) TTSOption {
	return func(r *V2TTSRequest) {
		r.Business.Aue = aue
	}
}

// WithStreamMode 设置是否开启流式返回
func WithStreamMode(enable bool) TTSOption {
	sfl := 0
	if enable {
		sfl = 1
	}
	return func(r *V2TTSRequest) {
		r.Business.Sfl = sfl
	}
}

// WithAudioFormat 设置音频采样率
func WithAudioFormat(auf string) TTSOption {
	return func(r *V2TTSRequest) {
		r.Business.Auf = auf
	}
}

// WithVoice 设置发音人
func WithVoice(vcn string) TTSOption {
	return func(r *V2TTSRequest) {
		r.Business.Vcn = vcn
	}
}

// WithSpeed 设置语速
func WithSpeed(speed int) TTSOption {
	return func(r *V2TTSRequest) {
		if speed < 0 {
			speed = 0
		}
		if speed > 100 {
			speed = 100
		}
		r.Business.Speed = speed
	}
}

// WithVolume 设置音量
func WithVolume(volume int) TTSOption {
	return func(r *V2TTSRequest) {
		if volume < 0 {
			volume = 0
		}
		if volume > 100 {
			volume = 100
		}
		r.Business.Volume = volume
	}
}

// WithPitch 设置音高
func WithPitch(pitch int) TTSOption {
	return func(r *V2TTSRequest) {
		if pitch < 0 {
			pitch = 0
		}
		if pitch > 100 {
			pitch = 100
		}
		r.Business.Pitch = pitch
	}
}

// WithBackgroundSound 设置是否开启背景音
func WithBackgroundSound(enable bool) TTSOption {
	bgs := 0
	if enable {
		bgs = 1
	}
	return func(r *V2TTSRequest) {
		r.Business.Bgs = bgs
	}
}

// WithTextEncoding 设置文本编码格式
func WithTextEncoding(tte string) TTSOption {
	return func(r *V2TTSRequest) {
		r.Business.Tte = tte
	}
}

// WithRegionalAccent 设置地方口音
func WithRegionalAccent(reg string) TTSOption {
	return func(r *V2TTSRequest) {
		r.Business.Reg = reg
	}
}

// WithNumberReading 设置数字发音方式
func WithNumberReading(rdn string) TTSOption {
	return func(r *V2TTSRequest) {
		r.Business.Rdn = rdn
	}
}

// ApplyOptions 应用多个选项到请求
func ApplyOptions(req *V2TTSRequest, options ...TTSOption) {
	for _, opt := range options {
		opt(req)
	}
}
