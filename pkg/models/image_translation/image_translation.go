package image_translation

import "encoding/base64"

type ASEImageTranslateRequest struct {
	Header struct {
		AppID  string `json:"app_id"`
		Status int    `json:"status"`
	} `json:"header"`
	Parameter struct {
		Td4D24Ede struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"td4d24ede"`
		T8F0377Ad struct {
			ReturnType int `json:"returnType"`
		} `json:"t8f0377ad"`
		Ocr struct {
			Language string `json:"language"`
		} `json:"ocr"`
	} `json:"parameter"`
	Payload struct {
		Image struct {
			Encoding string `json:"encoding"`
			Image    string `json:"image"`
			Status   int    `json:"status"`
		} `json:"image"`
	} `json:"payload"`
}

type ASEImageTranslateResponse struct {
	Header struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Sid     string `json:"sid"`
		Status  int    `json:"status"`
	} `json:"header"`
	Payload struct {
		ItsImage struct {
			Encoding string `json:"encoding"`
			Image    string `json:"image"`
		} `json:"its_image"`
		BlockText struct {
			Encoding string `json:"encoding"`
			Compress string `json:"compress"`
			Format   string `json:"format"`
			Text     string `json:"text"`
		} `json:"block_text"`
		ItsOutput struct {
			Compress string `json:"compress"`
			Encoding string `json:"encoding"`
			Format   string `json:"format"`
			Status   string `json:"status"`
			Text     string `json:"text"`
		} `json:"its_output"`
		Image struct {
			Encoding string `json:"encoding"`
			Image    string `json:"image"`
			Seq      string `json:"seq"`
			Status   string `json:"status"`
		} `json:"image"`
		JSON struct {
			Compress string `json:"compress"`
			Encoding string `json:"encoding"`
			Format   string `json:"format"`
			Seq      string `json:"seq"`
			Status   string `json:"status"`
			Text     string `json:"text"`
		} `json:"json"`
	} `json:"payload"`
}

type ASEImageTranslatePayload struct {
	ItsImage struct {
		Encoding string `json:"encoding"`
		Image    string `json:"image"`
	} `json:"its_image"`
	BlockText struct {
		Encoding string `json:"encoding"`
		Compress string `json:"compress"`
		Format   string `json:"format"`
		Text     string `json:"text"`
	} `json:"block_text"`
	ItsOutput struct {
		Compress string `json:"compress"`
		Encoding string `json:"encoding"`
		Format   string `json:"format"`
		Status   string `json:"status"`
		Text     string `json:"text"`
	} `json:"its_output"`
	Image struct {
		Encoding string `json:"encoding"`
		Image    string `json:"image"`
		Seq      string `json:"seq"`
		Status   string `json:"status"`
	} `json:"image"`
	JSON struct {
		Compress string `json:"compress"`
		Encoding string `json:"encoding"`
		Format   string `json:"format"`
		Seq      string `json:"seq"`
		Status   string `json:"status"`
		Text     string `json:"text"`
	} `json:"json"`
}

func NewASEImageTranslateRequest(appid, fromLang, toLang string, ocrLang string, imageData []byte) *ASEImageTranslateRequest {
	imageBase64 := base64.StdEncoding.EncodeToString(imageData)
	return &ASEImageTranslateRequest{
		Header: struct {
			AppID  string `json:"app_id"`
			Status int    `json:"status"`
		}{
			AppID:  appid,
			Status: 0,
		},
		Parameter: struct {
			Td4D24Ede struct {
				From string `json:"from"`
				To   string `json:"to"`
			} `json:"td4d24ede"`
			T8F0377Ad struct {
				ReturnType int `json:"returnType"`
			} `json:"t8f0377ad"`
			Ocr struct {
				Language string `json:"language"`
			} `json:"ocr"`
		}{
			Td4D24Ede: struct {
				From string `json:"from"`
				To   string `json:"to"`
			}{
				From: fromLang,
				To:   toLang,
			},
			T8F0377Ad: struct {
				ReturnType int `json:"returnType"`
			}{
				ReturnType: 3, // or other default value
			},
			Ocr: struct {
				Language string `json:"language"`
			}{
				Language: ocrLang, // assuming English or you can pass it as a parameter
			},
		},
		Payload: struct {
			Image struct {
				Encoding string `json:"encoding"`
				Image    string `json:"image"`
				Status   int    `json:"status"`
			} `json:"image"`
		}{
			Image: struct {
				Encoding string `json:"encoding"`
				Image    string `json:"image"`
				Status   int    `json:"status"`
			}{
				Encoding: "jpg",
				Image:    imageBase64,
				Status:   3,
			},
		},
	}
}
