package ase

// ASEHeader represents the header of the ASE API response
type ASEHeader struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Sid     string `json:"sid"`
	Status  int    `json:"status"`
}

// ASEBaseResponse represents the base structure of the ASE API response
type ASEBaseResponse struct {
	Header  ASEHeader              `json:"header"`
	Payload map[string]interface{} `json:"payload"`
}
