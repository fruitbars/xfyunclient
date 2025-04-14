package translation

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/fruitbars/xfyunclient/pkg/ase"
	"github.com/fruitbars/xfyunclient/pkg/models/translation_v1_its"
	"log"
)

func Translate(aseAppid, aseAPIKey, aseAPISecret string, serverUrl string, fromLang, toLang, text string) (string, error) {
	client := ase.NewASEHttpClient(serverUrl, aseAppid, aseAPIKey, aseAPISecret, "", "")

	log.Println(client)

	req := translation_v1_its.NewASETranslationRequest(aseAppid, fromLang, toLang, text)

	response, err := client.CallASEAPIJson(req)
	if err != nil {
		log.Println(err)
		return "", err
	}

	log.Println(string(response))

	var transRespone translation_v1_its.ASETranslationResult
	json.Unmarshal(response, &transRespone)
	if err != nil {
		log.Println(err, string(response))
		return "", err
	}

	if transRespone.Header.Code != 0 {
		return "", errors.New(transRespone.Header.Message)
	}

	dst, err := base64.StdEncoding.DecodeString(transRespone.Payload.Result.Text)
	if err != nil {
		return "", err
	}
	log.Println(string(dst))

	var engResult translation_v1_its.ASETranslationEngineResult
	err = json.Unmarshal(dst, &engResult)
	if err != nil {
		return "", err
	}
	log.Println(engResult.TransResult.Dst)

	return engResult.TransResult.Dst, nil
}
