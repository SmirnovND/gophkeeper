package domain

type Credentials struct {
	Id       int
	Login    string `json:"login"`
	Password string `json:"password"`
	PassHash string
}
