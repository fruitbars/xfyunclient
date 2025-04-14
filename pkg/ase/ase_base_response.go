package ase

// ASEBaseResponse represents the base structure of the ASE API response
type ASEBaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Sid     string `json:"sid"`
	Data    struct {
		Status int `json:"status"`
	} `json:"data"`
}
