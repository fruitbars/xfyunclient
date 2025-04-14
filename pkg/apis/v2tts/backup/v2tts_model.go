package v2tts

import "encoding/base64"

type V2TTSRequest struct {
	Common struct {
		AppID string `json:"app_id"`
	} `json:"common"`
	Business struct {
		Aue    string `json:"aue"`
		Sfl    int    `json:"sfl,omitempty"`
		Auf    string `json:"auf,omitempty"`
		Vcn    string `json:"vcn"`
		Speed  int    `json:"speed,omitempty"`
		Volume int    `json:"volume,omitempty"`

		Pitch int `json:"pitch,omitempty"`
		Bgs   int `json:"bgs,omitempty"`

		Tte string `json:"tte,omitempty"`
		Reg string `json:"reg,omitempty"`
		Rdn string `json:"rdn,omitempty"`
	} `json:"business"`
	Data struct {
		Status int    `json:"status"`
		Text   string `json:"text"`
	} `json:"data"`
}

type V2TTSResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Sid     string `json:"sid"`
	Data    struct {
		Audio  string `json:"audio"`
		Ced    string `json:"ced"`
		Status int    `json:"status"`
	} `json:"data"`
}

func NewV2TTSRequestDefault(appid, text, aue, vcn, tte string) V2TTSRequest {
	textBase64 := base64.StdEncoding.EncodeToString([]byte(text))
	return V2TTSRequest{
		Common: struct {
			AppID string `json:"app_id"`
		}{
			AppID: appid,
		},
		Business: struct {
			Aue    string `json:"aue"`
			Sfl    int    `json:"sfl,omitempty"`
			Auf    string `json:"auf,omitempty"`
			Vcn    string `json:"vcn"`
			Speed  int    `json:"speed,omitempty"`
			Volume int    `json:"volume,omitempty"`

			Pitch int `json:"pitch,omitempty"`
			Bgs   int `json:"bgs,omitempty"`

			Tte string `json:"tte,omitempty"`
			Reg string `json:"reg,omitempty"`
			Rdn string `json:"rdn,omitempty"`
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

func NewV2TTSRequest(appid, text, aue string, sfl int, auf, vcn string, speed, volume, pitch, bgs int, tte, reg, rdn string) V2TTSRequest {
	textBase64 := base64.StdEncoding.EncodeToString([]byte(text))
	return V2TTSRequest{
		Common: struct {
			AppID string `json:"app_id"`
		}{
			AppID: appid,
		},
		Business: struct {
			Aue    string `json:"aue"`
			Sfl    int    `json:"sfl,omitempty"`
			Auf    string `json:"auf,omitempty"`
			Vcn    string `json:"vcn"`
			Speed  int    `json:"speed,omitempty"`
			Volume int    `json:"volume,omitempty"`

			Pitch int `json:"pitch,omitempty"`
			Bgs   int `json:"bgs,omitempty"`

			Tte string `json:"tte,omitempty"`
			Reg string `json:"reg,omitempty"`
			Rdn string `json:"rdn,omitempty"`
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
