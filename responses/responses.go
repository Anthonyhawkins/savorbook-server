package responses

type ErrorResponse struct {
	FailedField string `json:"field"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
	Message     string `json:"message"`
}

type StandardResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  []string    `json:"errors"`
}
