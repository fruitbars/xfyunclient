package ocr_multi_lang

import "encoding/base64"

type Header struct {
	AppID  string `json:"app_id"`
	Status int    `json:"status"`
}

type OCRParameter struct {
	Language      string `json:"language"`
	BlockOption   string `json:"block_option,omitempty"`
	OCROutputText struct {
		Encoding string `json:"encoding"`
		Compress string `json:"compress"`
		Format   string `json:"format"`
	} `json:"ocr_output_text"`
}

type ImagePayload struct {
	Encoding string `json:"encoding"`
	Image    string `json:"image"`
	Status   int    `json:"status"`
}

type Payload struct {
	Image ImagePayload `json:"image"`
}

type Request struct {
	Header    Header                  `json:"header"`
	Parameter map[string]OCRParameter `json:"parameter"`
	Payload   Payload                 `json:"payload"`
}

func NewOcrRequest(appID string, language string, imageEncoding string, imageData []byte) Request {
	imageBase64 := base64.StdEncoding.EncodeToString(imageData)
	ocrParam := OCRParameter{
		Language:    language,
		BlockOption: "1",
	}
	ocrParam.OCROutputText.Encoding = "utf8"
	ocrParam.OCROutputText.Compress = "raw"
	ocrParam.OCROutputText.Format = "json"

	imagePayload := ImagePayload{
		Encoding: imageEncoding,
		Image:    imageBase64,
		Status:   3,
	}

	payload := Payload{
		Image: imagePayload,
	}

	header := Header{
		AppID:  appID,
		Status: 3,
	}

	parameter := map[string]OCRParameter{
		"ocr": ocrParam,
	}

	return Request{
		Header:    header,
		Parameter: parameter,
		Payload:   payload,
	}
}
