package salutespeech

type Response struct {
	StatusCode  int      `json:"status,omitempty"`
	Message     string   `json:"message,omitempty"`
	Result      []string `json:"result,omitempty"`
	AccessToken string   `json:"access_token,omitempty"`
}
