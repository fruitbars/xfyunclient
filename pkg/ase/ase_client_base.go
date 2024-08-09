package ase

import (
	"net/http"
	"net/url"
	"time"
	"xfyunclient/pkg/utils"
)

type ASEClientBase struct {
	serverUrl    string
	timeOut      time.Duration
	appid        string
	apikey       string
	apisecret    string
	HttpProto    string
	ASEAlgorithm string
}

// NewASEClientBase initializes a new ASEClientBase instance
func NewASEClientBase(serverURL, appid, apikey, apisecret, httpProto, aseAlgorithm string) *ASEClientBase {

	if httpProto == "" {
		httpProto = DefaultASEHttpProto
	}

	if aseAlgorithm == "" {
		aseAlgorithm = DefaultASEAlgorithm
	}

	return &ASEClientBase{
		serverUrl:    serverURL,
		timeOut:      30 * time.Second,
		appid:        appid,
		apikey:       apikey,
		apisecret:    apisecret,
		HttpProto:    httpProto,
		ASEAlgorithm: aseAlgorithm,
	}
}

// getAuthServerURL generates the authenticated server URL and headers
func (c *ASEClientBase) getAuthServerURL(method string) (string, *http.Header, error) {
	currentTime := time.Now().UTC().Format(time.RFC1123)
	host, path, err := utils.ExtractHostAndPath(c.serverUrl)
	if err != nil {
		return "", nil, err
	}

	signature := GenerateSignature(host, currentTime, method, path, c.HttpProto, c.apisecret)
	authorizationBstr := GenerateAuthorization(c.apikey, c.ASEAlgorithm, signature)

	v := url.Values{}
	v.Add("authorization", authorizationBstr)
	v.Add("date", currentTime)
	v.Add("host", host)

	callURL := c.serverUrl + "?" + v.Encode()

	headers := &http.Header{}
	headers.Add("authorization", authorizationBstr)
	headers.Add("date", currentTime)
	headers.Add("host", host)

	return callURL, headers, nil
}
