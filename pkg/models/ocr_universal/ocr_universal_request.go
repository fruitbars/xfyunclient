package ocr_universal

import "encoding/base64"

type Header struct {
	AppID  string `json:"app_id"`
	Status int    `json:"status"`
}

type Result struct {
	Encoding string `json:"encoding"`
	Compress string `json:"compress"`
	Format   string `json:"format"`
}

type Parameter struct {
	SF8E6ACA1 struct {
		Category string `json:"category"`
		Result   Result `json:"result"`
	} `json:"sf8e6aca1"`
}

type Payload struct {
	SF8E6ACA1Data1 struct {
		Encoding string `json:"encoding"`
		Status   int    `json:"status"`
		Image    string `json:"image"`
	} `json:"sf8e6aca1_data_1"`
}

type DataStruct struct {
	Header    Header    `json:"header"`
	Parameter Parameter `json:"parameter"`
	Payload   Payload   `json:"payload"`
}

func NewOcrUniversalRequest(appID string, category, imageEncoding string, imageData []byte) DataStruct {
	imageBase64 := base64.StdEncoding.EncodeToString(imageData)
	return DataStruct{
		Header: Header{
			AppID:  appID,
			Status: 3,
		},
		Parameter: Parameter{
			SF8E6ACA1: struct {
				Category string `json:"category"`
				Result   Result `json:"result"`
			}{
				Category: category,
				Result: Result{
					Encoding: "utf8",
					Compress: "raw",
					Format:   "json",
				},
			},
		},
		Payload: Payload{
			SF8E6ACA1Data1: struct {
				Encoding string `json:"encoding"`
				Status   int    `json:"status"`
				Image    string `json:"image"`
			}{
				Encoding: imageEncoding,
				Status:   3,
				Image:    imageBase64,
			},
		},
	}
}
