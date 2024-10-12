package models

// swagger:model AudioAnalyzeResponse
type AudioAnalyzeResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

// swagger:model PhotoAnalyzeResponse
type PhotoAnalyzeResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

// swagger:model ErrorResponse
type ErrorResponse struct {
	Error string `json:"error"`
}
