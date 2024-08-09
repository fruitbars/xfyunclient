package language_detect

type ASELanguageDetectRequest struct {
	Header struct {
		AppID  string `json:"app_id"`
		Status int    `json:"status"`
	} `json:"header"`
	Parameter struct {
		Cnen struct {
			Outfmt string `json:"outfmt"`
			Result struct {
				Encoding string `json:"encoding"`
				Compress string `json:"compress"`
				Format   string `json:"format"`
			} `json:"result"`
		} `json:"cnen"`
	} `json:"parameter"`
	Payload struct {
		Request struct {
			Encoding string `json:"encoding"`
			Compress string `json:"compress"`
			Format   string `json:"format"`
			Status   int    `json:"status"`
			Text     string `json:"text"`
		} `json:"request"`
	} `json:"payload"`
}

type ASELanguageDetectResponse struct {
	Header struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Sid     string `json:"sid"`
	} `json:"header"`
	Payload struct {
		Result struct {
			Encoding string `json:"encoding"`
			Compress string `json:"compress"`
			Format   string `json:"format"`
			Text     string `json:"text"`
		} `json:"result"`
	} `json:"payload"`
}

type ASELanguageDetectTranResult struct {
	TransResult []struct {
		Src      string `json:"src"`
		LanProbs string `json:"lan_probs"`
	} `json:"trans_result"`
}
