package v2tts

// pkg/v2tts/models.go

import "encoding/base64"

// V2TTSRequest 表示向语音合成API发送的请求
// 包含应用信息、业务参数和要合成的文本数据
type V2TTSRequest struct {
	Common struct {
		AppID string `json:"app_id"` // 应用ID
	} `json:"common"`
	Business struct {
		Aue    string `json:"aue"`              // 音频编码格式，如wav、mp3
		Sfl    int    `json:"sfl,omitempty"`    // 是否开启流式返回
		Auf    string `json:"auf,omitempty"`    // 音频采样率
		Vcn    string `json:"vcn"`              // 发音人
		Speed  int    `json:"speed,omitempty"`  // 语速，取值范围：[0,100]
		Volume int    `json:"volume,omitempty"` // 音量，取值范围：[0,100]
		Pitch  int    `json:"pitch,omitempty"`  // 音高，取值范围：[0,100]
		Bgs    int    `json:"bgs,omitempty"`    // 是否开启背景音
		Tte    string `json:"tte,omitempty"`    // 文本编码格式，如utf8
		Reg    string `json:"reg,omitempty"`    // 地方口音
		Rdn    string `json:"rdn,omitempty"`    // 数字发音方式
	} `json:"business"`
	Data struct {
		Status int    `json:"status"` // 状态，固定为2
		Text   string `json:"text"`   // base64编码的文本
	} `json:"data"`
}

// V2TTSResponse 表示从语音合成API接收的响应
type V2TTSResponse struct {
	Code    int    `json:"code"`    // 返回码，0表示成功
	Message string `json:"message"` // 返回信息
	Sid     string `json:"sid"`     // 会话ID
	Data    struct {
		Audio  string `json:"audio"`  // base64编码的音频数据
		Ced    string `json:"ced"`    // 合成的额外信息
		Status int    `json:"status"` // 状态
	} `json:"data"`
}

// CreateDefaultV2TTSRequest 创建默认的TTS请求
func CreateDefaultV2TTSRequest(appID, text, aue, vcn, tte string) V2TTSRequest {
	textBase64 := base64.StdEncoding.EncodeToString([]byte(text))
	return V2TTSRequest{
		Common: struct {
			AppID string `json:"app_id"`
		}{
			AppID: appID,
		},
		Business: struct {
			Aue    string `json:"aue"`
			Sfl    int    `json:"sfl,omitempty"`
			Auf    string `json:"auf,omitempty"`
			Vcn    string `json:"vcn"`
			Speed  int    `json:"speed,omitempty"`
			Volume int    `json:"volume,omitempty"`
			Pitch  int    `json:"pitch,omitempty"`
			Bgs    int    `json:"bgs,omitempty"`
			Tte    string `json:"tte,omitempty"`
			Reg    string `json:"reg,omitempty"`
			Rdn    string `json:"rdn,omitempty"`
		}{
			Aue: aue,
			Vcn: vcn,
			Tte: tte,
		},
		Data: struct {
			Status int    `json:"status"`
			Text   string `json:"text"`
		}{
			Status: 2,
			Text:   textBase64,
		},
	}
}

// CreateV2TTSRequest 创建自定义参数的TTS请求
func CreateV2TTSRequest(appID, text, aue string, sfl int, auf, vcn string, speed, volume, pitch, bgs int, tte, reg, rdn string) V2TTSRequest {
	textBase64 := base64.StdEncoding.EncodeToString([]byte(text))
	return V2TTSRequest{
		Common: struct {
			AppID string `json:"app_id"`
		}{
			AppID: appID,
		},
		Business: struct {
			Aue    string `json:"aue"`
			Sfl    int    `json:"sfl,omitempty"`
			Auf    string `json:"auf,omitempty"`
			Vcn    string `json:"vcn"`
			Speed  int    `json:"speed,omitempty"`
			Volume int    `json:"volume,omitempty"`
			Pitch  int    `json:"pitch,omitempty"`
			Bgs    int    `json:"bgs,omitempty"`
			Tte    string `json:"tte,omitempty"`
			Reg    string `json:"reg,omitempty"`
			Rdn    string `json:"rdn,omitempty"`
		}{
			Aue:    aue,
			Sfl:    sfl,
			Auf:    auf,
			Vcn:    vcn,
			Speed:  speed,
			Volume: volume,
			Pitch:  pitch,
			Bgs:    bgs,
			Tte:    tte,
			Reg:    reg,
			Rdn:    rdn,
		},
		Data: struct {
			Status int    `json:"status"`
			Text   string `json:"text"`
		}{
			Status: 2,
			Text:   textBase64,
		},
	}
}
