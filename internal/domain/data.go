package domain

// CredentialData представляет собой структуру для хранения пары логин/пароль
type CredentialData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// CardData представляет собой структуру для хранения данных кредитной карты
type CardData struct {
	Number     string `json:"number"`
	Holder     string `json:"holder"`
	ExpiryDate string `json:"expiry_date"`
	CVV        string `json:"cvv"`
}

// TextData представляет собой структуру для хранения произвольного текста
type TextData struct {
	Content string `json:"content"`
}

type FileData struct {
	Name      string `json:"name" binding:"required"`
	Extension string `json:"extension" binding:"required"`
}

type FileDataResponse struct {
	Url         string `json:"url" binding:"required"`
	Description string `json:"description" binding:"required"`
}
