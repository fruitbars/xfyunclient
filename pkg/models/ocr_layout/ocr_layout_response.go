package ocr_layout

type ASEOCRLayoutResponse struct {
	Header struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Sid     string `json:"sid"`
	} `json:"header"`
	Payload struct {
		OcrOutputText struct {
			Encoding string `json:"encoding"`
			Compress string `json:"compress"`
			Format   string `json:"format"`
			Text     string `json:"text"`
		} `json:"ocr_output_text"`
	} `json:"payload"`
}
