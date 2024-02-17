package telegram

type Update struct {
	Id      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"`
}

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type IncomingMessage struct {
	Text  string `json:"text"`
	From  User   `json:"from"`
	Chat  Chat   `json:"chat"`
	Voice Voice  `json:"voice"`
}

type Voice struct {
	FileId string `json:"file_id"`
}

type User struct {
	Username string `json:"username"`
	Id       int    `json:"id"`
}

type Chat struct {
	ID int `json:"id"`
}

type File struct {
	FileId   string `json:"file_id"`
	FilePath string `json:"file_path,omitempty"`
}

type GetFileResponse struct {
	Ok     bool   `json:"ok"`
	Error  string `json:"error"`
	Result File   `json:"result"`
}
