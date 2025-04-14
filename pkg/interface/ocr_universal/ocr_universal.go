package ocr_universal

import (
	"encoding/base64"
	"encoding/json"
	"github.com/fruitbars/xfyunclient/pkg/ase"
	"github.com/fruitbars/xfyunclient/pkg/models/ocr_universal"
	"github.com/fruitbars/xfyunclient/pkg/utils"
	"log"
)

func OcrUniversal(aseAppid, aseAPIKey, aseAPISecret string, serverUrl string, language string, fname string) (*ocr_universal.Response, error) {
	client := ase.NewASEHttpClient(serverUrl, aseAppid, aseAPIKey, aseAPISecret, "", "")

	format, imageBase64, err := utils.ReadImageFile(fname)
	log.Println(format)
	//req := ocr_universal_2024_request.NewOcrUniversal2024Request(aseAppid, "ch_en_public_cloud", format, imageBase64)
	req := ocr_universal.NewOcrUniversalRequest(aseAppid, language, format, imageBase64)

	//utils.PrintJson(req)
	response, err := client.CallASEAPIJson(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	log.Println(string(response))

	var ocrResponse ocr_universal.Response
	json.Unmarshal(response, &ocrResponse)

	textData, err := base64.StdEncoding.DecodeString(ocrResponse.Payload.Result.Text)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println(string(textData))

	return nil, nil
}
