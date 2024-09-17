package utils

type (
	LoginResponse struct {
		AccesssToken string `json:"access_token"`
	}

	SuccessResponse struct {
		Message string `json:"message"`
		Data    any    `json:"data"`
	}

	SuccessPaginationResponse struct {
		Message string `json:"message"`
		Data    any    `json:"data"`
		Total   int    `json:"total"`
	}

	ErrorResponse struct {
		Message string `json:"message"`
		Error   error  `json:"error"`
	}

	ValidateErrorResponse struct {
		Message string           `json:"message"`
		Error   []*ValidateError `json:"error"`
	}

	ValidateError struct {
		Field string `json:"field"`
		Tag   string `json:"tag"`
		Value string `json:"value"`
	}
)
