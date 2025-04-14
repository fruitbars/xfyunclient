package ocr_layout

import "encoding/base64"

type ASEOCRLayoutRequest struct {
	Header struct {
		AppID     string  `json:"app_id"`
		UID       string  `json:"uid,omitempty"`
		DID       string  `json:"did,omitempty"`
		Imei      string  `json:"imei,omitempty"`
		Imsi      string  `json:"imsi,omitempty"`
		Mac       string  `json:"mac,omitempty"`
		NetType   string  `json:"net_type,omitempty"`
		NetIsp    string  `json:"net_isp,omitempty"`
		Status    int     `json:"status"`
		RequestID *string `json:"request_id,omitempty"`
		ResID     string  `json:"res_id,omitempty"`
	} `json:"header"`
	Parameter struct {
		Ocr struct {
			Language      string `json:"language"`
			OcrOutputText struct {
				Encoding string `json:"encoding"`
				Compress string `json:"compress"`
				Format   string `json:"format"`
			} `json:"ocr_output_text"`
		} `json:"ocr"`
	} `json:"parameter"`
	Payload struct {
		Image struct {
			Encoding string `json:"encoding"`
			Image    string `json:"image"`
			Status   int    `json:"status"`
		} `json:"image"`
	} `json:"payload"`
}

// NewASEOCRLayoutRequest initializes and returns a pointer to an ASEOCRLayoutRequest structure.
func NewASEOCRLayoutRequest(appid string, imageEncoding string, imageData []byte) *ASEOCRLayoutRequest {
	imageBase64 := base64.StdEncoding.EncodeToString(imageData)
	return &ASEOCRLayoutRequest{
		Header: struct {
			AppID     string  `json:"app_id"`
			UID       string  `json:"uid,omitempty"`
			DID       string  `json:"did,omitempty"`
			Imei      string  `json:"imei,omitempty"`
			Imsi      string  `json:"imsi,omitempty"`
			Mac       string  `json:"mac,omitempty"`
			NetType   string  `json:"net_type,omitempty"`
			NetIsp    string  `json:"net_isp,omitempty"`
			Status    int     `json:"status"`
			RequestID *string `json:"request_id,omitempty"`
			ResID     string  `json:"res_id,omitempty"`
		}{
			AppID:  appid,
			Status: 3,
		},
		Parameter: struct {
			Ocr struct {
				Language      string `json:"language"`
				OcrOutputText struct {
					Encoding string `json:"encoding"`
					Compress string `json:"compress"`
					Format   string `json:"format"`
				} `json:"ocr_output_text"`
			} `json:"ocr"`
		}{
			Ocr: struct {
				Language      string `json:"language"`
				OcrOutputText struct {
					Encoding string `json:"encoding"`
					Compress string `json:"compress"`
					Format   string `json:"format"`
				} `json:"ocr_output_text"`
			}{
				Language: "layout", // Default language set to English
				OcrOutputText: struct {
					Encoding string `json:"encoding"`
					Compress string `json:"compress"`
					Format   string `json:"format"`
				}{
					Encoding: "utf8", // Default encoding set to UTF-8
					Compress: "raw",
					Format:   "json",
				},
			},
		},
		Payload: struct {
			Image struct {
				Encoding string `json:"encoding"`
				Image    string `json:"image"`
				Status   int    `json:"status"`
			} `json:"image"`
		}{
			Image: struct {
				Encoding string `json:"encoding"`
				Image    string `json:"image"`
				Status   int    `json:"status"`
			}{
				Encoding: imageEncoding,
				Image:    imageBase64,
				Status:   3,
			},
		},
	}
}
