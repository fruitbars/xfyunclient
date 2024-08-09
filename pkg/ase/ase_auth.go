package ase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
)

func GenerateHMAC(data, key string) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return h.Sum(nil)
}

// GenerateSignature generates the HMAC signature for the given parameters
func GenerateSignature(host, date, httpMethod, requestUri, httpProto, secret string) string {
	var signatureStr string
	if len(host) != 0 {
		signatureStr = "host: " + host + "\n"
	}
	signatureStr += "date: " + date + "\n"
	signatureStr += httpMethod + " " + requestUri + " " + httpProto

	log.Println(signatureStr)

	sign := GenerateHMAC(signatureStr, secret)

	return base64.StdEncoding.EncodeToString(sign)
}

// GenerateAuthorization creates the authorization header based on the given parameters.
func GenerateAuthorization(apiKey, algorithm, signature string) string {
	str := fmt.Sprintf(`api_key="%s", algorithm="%s", headers="host date request-line", signature="%s"`, apiKey, algorithm, signature)

	return base64.StdEncoding.EncodeToString([]byte(str))
}

func GenerateAuthorizationOld(apiKey, algorithm, signature string) string {
	str := fmt.Sprintf(`hmac username="%s", algorithm="%s", headers="host date request-line", signature="%s"`, apiKey, algorithm, signature)

	log.Println(str)

	return base64.StdEncoding.EncodeToString([]byte(str))
}
