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

type CountResponse struct {
	Count int `json:"count"`
}

type StatsResponse struct {
	Timeline     map[string]int `json:"timeline"`
	Referrals    map[string]int `json:"referrals"`
	Browsers     map[string]int `json:"browsers"`
	Countries    map[string]int `json:"countries"`
	OS           map[string]int `json:"os"`
	ErrorMessage string         `json:"errorMessage,omitempty"`
}
