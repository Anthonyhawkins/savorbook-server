package validators

type ErrorResponse struct {
	FailedField string `json:"field"`
	Tag         string `json:"tag"`
	value       string `json:"value"`
	Message     string `json:"message"`
}
