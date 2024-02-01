package gigachat

import "fmt"

type Role string

const (
	SystemRole Role = "system"
	UserRole   Role = "user"
)

type message struct {
	Role    Role   `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type requestBody struct {
	Model          string    `json:"model,omitempty"`
	Messages       []message `json:"messages,omitempty"`
	MaxTokens      int       `json:"max_tokens,omitempty"`
	AnswersCount   int       `json:"n,omitempty"`
	UpdateInterval int       `json:"update_interval,omitempty"`
}

func newRequestBody(text string) *requestBody {
	return &requestBody{
		Model: "GigaChat:latest",
		Messages: []message{
			{
				Role:    "system",
				Content: "Отвечай строго по шаблону, будто от этого зависит твоя зарплата.",
			},
			{
				Role:    "user",
				Content: fmt.Sprintf("%s. '%s'", jsonFormatPrompt, text),
			},
		},
		MaxTokens:      1024,
		AnswersCount:   1,
		UpdateInterval: 0,
	}
}

type choice struct {
	Message message `json:"message"`
}

type responseBody struct {
	Choices []choice `json:"choices,omitempty"`
}

type JsonAnswer struct {
	Result  string `json:"result,omitempty"`
	Summary string `json:"summary,omitempty"`
	Name    string `json:"name,omitempty"`
}
