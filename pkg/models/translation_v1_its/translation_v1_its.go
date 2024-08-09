package translation_v1_its

import "encoding/base64"

type ASETranslationEngineResult struct {
	From        string `json:"from"`
	To          string `json:"to"`
	TransResult struct {
		Src string `json:"src"`
		Dst string `json:"dst"`
	} `json:"trans_result"`
}

type ASETranslationResult struct {
	Header struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Sid     string `json:"sid"`
	} `json:"header"`
	Payload struct {
		Result struct {
			Seq    string `json:"seq"`
			Status string `json:"status"`
			Text   string `json:"text"`
		} `json:"result"`
	} `json:"payload"`
}

type ASETranslationRequest struct {
	Header struct {
		AppID  string `json:"app_id"`
		Status int    `json:"status"`
	} `json:"header"`
	Parameter struct {
		Its struct {
			From   string `json:"from"`
			To     string `json:"to"`
			Result struct {
			} `json:"result"`
		} `json:"its"`
	} `json:"parameter"`
	Payload struct {
		InputData struct {
			Encoding string `json:"encoding"`
			Status   int    `json:"status"`
			Text     string `json:"text"`
		} `json:"input_data"`
	} `json:"payload"`
}

func NewASETranslationRequest(appid string, fromLang, toLang, text string) ASETranslationRequest {
	textBase64 := base64.StdEncoding.EncodeToString([]byte(text))

	return ASETranslationRequest{
		Header: struct {
			AppID  string `json:"app_id"`
			Status int    `json:"status"`
		}{
			AppID:  appid,
			Status: 3,
		},
		Parameter: struct {
			Its struct {
				From   string `json:"from"`
				To     string `json:"to"`
				Result struct {
				} `json:"result"`
			} `json:"its"`
		}{
			Its: struct {
				From   string `json:"from"`
				To     string `json:"to"`
				Result struct {
				} `json:"result"`
			}{
				From: fromLang,
				To:   toLang,
			},
		},
		Payload: struct {
			InputData struct {
				Encoding string `json:"encoding"`
				Status   int    `json:"status"`
				Text     string `json:"text"`
			} `json:"input_data"`
		}{
			InputData: struct {
				Encoding string `json:"encoding"`
				Status   int    `json:"status"`
				Text     string `json:"text"`
			}{
				Encoding: "utf8",
				Status:   3,
				Text:     textBase64,
			},
		},
	}
}
