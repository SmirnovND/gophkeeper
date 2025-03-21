package domain

type CredentialData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type FileData struct {
	Name      string `json:"name" binding:"required"`
	Extension string `json:"extension" binding:"required"`
}

type FileDataResponse struct {
	Url         string `json:"url" binding:"required"`
	Description string `json:"description" binding:"required"`
}
