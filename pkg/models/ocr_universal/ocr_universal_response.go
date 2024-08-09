package ocr_universal

type RespHeader struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	SID     string `json:"sid"`
}

type RespResult struct {
	Compress string `json:"compress"`
	Encoding string `json:"encoding"`
	Format   string `json:"format"`
	Text     string `json:"text"`
}

type RespPayload struct {
	Result RespResult `json:"result"`
}

type Response struct {
	Header  RespHeader  `json:"header"`
	Payload RespPayload `json:"payload"`
}
