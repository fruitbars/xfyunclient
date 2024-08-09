package ocr_universal_2024_response

type Header struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	SID     string `json:"sid"`
}

type OCROutputText struct {
	Encoding string `json:"encoding"`
	Compress string `json:"compress"`
	Format   string `json:"format"`
	Text     string `json:"text"`
}

type Payload struct {
	OCROutputText OCROutputText `json:"ocr_output_text"`
}

type OCRResponse struct {
	Header  Header  `json:"header"`
	Payload Payload `json:"payload"`
}
