package response

type OKResponse[T any] struct {
	Status string `json:"status"`
	Data   T      `json:"data,omitempty"`
}

type FailedResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
	Reason string `json:"reason,omitempty"`
}
