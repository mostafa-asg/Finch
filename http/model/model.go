package model

type HashResponse struct {
	TinyUrl      string `json:"tiny,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

type HashRequest struct {
	Url string `json:"url"`
}

type GetResponse struct {
	Found bool   `json:"found"`
	Url   string `json:"url,omitempty"`
}
