package ocr_layout

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fruitbars/xfyunclient/pkg/ase"
	"github.com/fruitbars/xfyunclient/pkg/models/ocr_layout"
	"github.com/fruitbars/xfyunclient/pkg/utils"
	"log"
)

func OCRLayoutFile(aseAppid, aseAPIKey, aseAPISecret string, serverUrl string, fname string) (error, string) {
	format, imageData, err := utils.ReadImageFile(fname)
	log.Println(format)

	if err != nil {
		return err, ""
	}

	return OCRLayout(aseAppid, aseAPIKey, aseAPISecret, serverUrl, format, imageData)
}

func OCRLayout(aseAppid, aseAPIKey, aseAPISecret string, serverUrl string, format string, imageData []byte) (error, string) {

	client := ase.NewASEHttpClient(serverUrl, aseAppid, aseAPIKey, aseAPISecret, "", "")

	log.Println(client.ASEClientBase)

	req := ocr_layout.NewASEOCRLayoutRequest(aseAppid, format, imageData)

	respData, err := client.CallASEAPIJson(req)
	if err != nil {
		return err, ""
	}

	if respData != nil {
		//log.Println(string(respData))

		var ocrLayoutResponse ocr_layout.ASEOCRLayoutResponse
		err = json.Unmarshal(respData, &ocrLayoutResponse)
		if err != nil {
			log.Println(err)
			return err, ""
		}

		if ocrLayoutResponse.Header.Code != 0 {
			log.Println(ocrLayoutResponse.Header.Message)
			err := errors.New(fmt.Sprintf("code: %d, message: %s", ocrLayoutResponse.Header.Code, ocrLayoutResponse.Header.Message))
			return err, ""
		}

		decodedBytes, err := base64.StdEncoding.DecodeString(ocrLayoutResponse.Payload.OcrOutputText.Text)
		if err != nil {
			fmt.Println("Error decoding:", err)
			return err, ""
		}
		decodedString := string(decodedBytes)

		log.Println(decodedString)

		return err, decodedString
	}

	return nil, ""
}
