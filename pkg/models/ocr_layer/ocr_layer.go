package ocr_layer

type ASEOCRLayerRequest struct {
	Header struct {
		AppID  string `json:"app_id"`
		Status int    `json:"status"`
	} `json:"header"`
	Parameter struct {
		Iocrld struct {
			JSON struct {
				Encoding string `json:"encoding"`
				Compress string `json:"compress"`
				Format   string `json:"format"`
			} `json:"json"`
			Image struct {
				Encoding string `json:"encoding"`
			} `json:"image"`
		} `json:"iocrld"`
	} `json:"parameter"`
	Payload struct {
		JSON struct {
			Encoding string `json:"encoding"`
			Compress string `json:"compress"`
			Format   string `json:"format"`
			Status   int    `json:"status"`
			Text     string `json:"text"`
		} `json:"json"`
		Image struct {
			Encoding string `json:"encoding"`
			Image    string `json:"image"`
			Status   int    `json:"status"`
		} `json:"image"`
	} `json:"payload"`
}

type ASEOCRLayerResponse struct {
	Header struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Sid     string `json:"sid"`
	} `json:"header"`
	Payload struct {
		JSON struct {
			Encoding string `json:"encoding"`
			Compress string `json:"compress"`
			Format   string `json:"format"`
			Text     string `json:"text"`
		} `json:"json"`
		Image struct {
			Encoding string `json:"encoding"`
			Image    string `json:"image"`
		} `json:"image"`
	} `json:"payload"`
}
