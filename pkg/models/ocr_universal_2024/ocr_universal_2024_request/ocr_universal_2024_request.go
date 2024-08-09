package ocr_universal_2024_request

import "encoding/base64"

type Header struct {
	AppID  string `json:"app_id"`
	Status int    `json:"status"`
}

type OCROutputText struct {
	Encoding string `json:"encoding"`
	Compress string `json:"compress"`
	Format   string `json:"format"`
}

type OCRParameter struct {
	Language       string        `json:"language"`
	TableOption    string        `json:"table_option"`
	ElementOption  string        `json:"element_option"`
	CharOption     string        `json:"char_option"`
	DocumentOption string        `json:"document_option"`
	ResultOption   string        `json:"result_option"`
	OCROutputText  OCROutputText `json:"ocr_output_text"`
}

type Parameter struct {
	OCR OCRParameter `json:"ocr"`
}

type Image struct {
	Encoding string `json:"encoding"`
	Image    string `json:"image"`
	Status   int    `json:"status"`
}

type Payload struct {
	Image Image `json:"image"`
}

type OCRRequest struct {
	Header    Header    `json:"header"`
	Parameter Parameter `json:"parameter"`
	Payload   Payload   `json:"payload"`
}

func NewOcrUniversal2024Request(appID string, category, imageEncoding string, imageData []byte) OCRRequest {
	imageBase64 := base64.StdEncoding.EncodeToString(imageData)

	return OCRRequest{
		Header: Header{
			AppID:  appID,
			Status: 3,
		},
		Parameter: Parameter{
			OCR: OCRParameter{
				Language:       category,
				TableOption:    "1",
				ElementOption:  "1",
				CharOption:     "0",
				DocumentOption: "off",
				ResultOption:   "all",
				OCROutputText: OCROutputText{
					Encoding: "utf8",
					Compress: "raw",
					Format:   "json",
				},
			},
		},
		Payload: Payload{
			Image: Image{
				Encoding: imageEncoding,
				Image:    imageBase64,
				Status:   3,
			},
		},
	}
}
