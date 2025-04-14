package ase

import (
	"github.com/fruitbars/xfyunclient/pkg/utils"
	"log"
	"net/http"
	"net/url"
	"time"
)

type ASEClientBase struct {
	ServerUrl    string
	TimeOut      time.Duration
	Appid        string
	Apikey       string
	Apisecret    string
	HttpProto    string
	ASEAlgorithm string
}

// NewASEClientBase initializes a new ASEClientBase instance
func NewASEClientBase(serverURL, appid, apikey, apisecret, httpProto, aseAlgorithm string) ASEClientBase {

	if httpProto == "" {
		httpProto = DefaultASEHttpProto
	}

	if aseAlgorithm == "" {
		aseAlgorithm = DefaultASEAlgorithm
	}

	return ASEClientBase{
		ServerUrl:    serverURL,
		TimeOut:      30 * time.Second,
		Appid:        appid,
		Apikey:       apikey,
		Apisecret:    apisecret,
		HttpProto:    httpProto,
		ASEAlgorithm: aseAlgorithm,
	}
}

// getAuthServerURL generates the authenticated server URL and headers
func (c *ASEClientBase) getAuthServerURL(method string) (string, *http.Header, error) {
	currentTime := time.Now().UTC().Format(time.RFC1123)
	log.Println(c.ServerUrl)
	host, path, err := utils.ExtractHostAndPath(c.ServerUrl)
	if err != nil {
		return "", nil, err
	}

	log.Println(host, path)
	signature := GenerateSignature(host, currentTime, method, path, c.HttpProto, c.Apisecret)
	authorizationBstr := GenerateAuthorization(c.Apikey, c.ASEAlgorithm, signature)

	v := url.Values{}
	v.Add("authorization", authorizationBstr)
	v.Add("date", currentTime)
	v.Add("host", host)

	callURL := c.ServerUrl + "?" + v.Encode()

	headers := &http.Header{}
	headers.Add("authorization", authorizationBstr)
	headers.Add("date", currentTime)
	headers.Add("host", host)

	log.Println(callURL, headers)

	return callURL, headers, nil
}
