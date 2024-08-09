package image_block

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
	"strings"
	"xfyunclient/pkg/ase"
	"xfyunclient/pkg/models/ocr_multi_lang"
	"xfyunclient/pkg/models/ocr_universal_2024/ocr_universal_2024_engine_text"
	"xfyunclient/pkg/models/ocr_universal_2024/ocr_universal_2024_request"
	"xfyunclient/pkg/models/ocr_universal_2024/ocr_universal_2024_response"
	"xfyunclient/pkg/utils"
)

func OcrBlockTest(aseAppid, aseAPIKey, aseAPISecret string, serverUrl string, language string, fname string) (string, error) {

	client := ase.NewASEClient(serverUrl, aseAppid, aseAPIKey, aseAPISecret, "", "")

	format, imageBase64, err := utils.ReadImageFile(fname)
	log.Println(format)

	req := ocr_multi_lang.NewOcrRequest(aseAppid, language, format, imageBase64)

	//utils.PrintJson(req)
	response, err := client.CallASEAPIJson(req)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var ocrResponse ocr_multi_lang.Response
	json.Unmarshal(response, &ocrResponse)

	textData, err := base64.StdEncoding.DecodeString(ocrResponse.Payload.OCROutputText.Text)
	if err != nil {
		log.Println(err)
		return "", err
	}

	//log.Println(textData)

	var ocrEngineText ocr_multi_lang.OCRResponseText
	err = json.Unmarshal(textData, &ocrEngineText)
	if err != nil {
		log.Println(err, textData)
		return "", err
	}

	s := strings.ReplaceAll(serverUrl, "/", "_")

	outFname := s + "_resp_" + language + "_" + ocrEngineText.Protoc + "_" + ocrEngineText.Version
	os.WriteFile(outFname, textData, 0666)

	log.Println(outFname)

	//log.Println(ocrEngineText)
	blockContents := ocr_multi_lang.GetBlockContents(&ocrEngineText)

	result := utils.MapToString(blockContents)

	log.Println(result)

	return result, nil
}

func OcrUniversal2024Test(aseAppid, aseAPIKey, aseAPISecret string, serverUrl string, language string, fname string) (string, error) {

	client := ase.NewASEClient(serverUrl, aseAppid, aseAPIKey, aseAPISecret, "", "")

	format, imageBase64, err := utils.ReadImageFile(fname)
	log.Println(format)
	//req := ocr_universal_2024_request.NewOcrUniversal2024Request(aseAppid, "ch_en_public_cloud", format, imageBase64)
	req := ocr_universal_2024_request.NewOcrUniversal2024Request(aseAppid, language, format, imageBase64)

	//utils.PrintJson(req)
	response, err := client.CallASEAPIJson(req)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	var ocrResponse ocr_universal_2024_response.OCRResponse
	json.Unmarshal(response, &ocrResponse)

	textData, err := base64.StdEncoding.DecodeString(ocrResponse.Payload.OCROutputText.Text)
	if err != nil {
		log.Println(err)
		return "", err
	}

	//	log.Println(textData)

	//os.WriteFile("resp.json", textData, 0666)
	var ocrEngineText ocr_universal_2024_engine_text.OCREngineText
	err = json.Unmarshal(textData, &ocrEngineText)
	if err != nil {
		log.Println(err)
		log.Println(string(response))
		return "", err
	}

	s := strings.ReplaceAll(serverUrl, "/", "_")

	outFname := s + "_resp_" + language + "_" + ocrEngineText.Protoc + "_" + ocrEngineText.Version
	os.WriteFile(outFname, textData, 0666)

	log.Println(outFname)

	//	log.Println(ocrEngineText)
	blockContents := ocr_multi_lang.GetUniversal2024BlockContents(&ocrEngineText)

	result := utils.MapToString(blockContents)

	log.Println(blockContents)

	return result, nil
}
