package ocr_multi_lang

type RespHeader struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Sid     string `json:"sid"`
}

type OCROutputText struct {
	Compress string `json:"compress"`
	Encoding string `json:"encoding"`
	Format   string `json:"format"`
	Seq      string `json:"seq"`
	Status   string `json:"status"`
	Text     string `json:"text"`
}

type RespPayload struct {
	OCROutputText OCROutputText `json:"ocr_output_text"`
}

type Response struct {
	Header  RespHeader  `json:"header"`
	Payload RespPayload `json:"payload"`
}
