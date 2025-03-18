package controllers

import (
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"github.com/SmirnovND/gophkeeper/internal/usecase"
	"github.com/SmirnovND/toolbox/pkg/paramsparser"
	"net/http"
)

type AuthController struct {
	AuthUseCase *usecase.AuthUseCase
}

func NewAuthController(AuthUseCase *usecase.AuthUseCase) *AuthController {
	return &AuthController{
		AuthUseCase: AuthUseCase,
	}
}

// HandleRegisterJSON godoc
// @Summary Регистрация нового пользователя
// @Description Регистрирует нового пользователя в системе и возвращает токен авторизации
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body domain.Credentials true "Учетные данные пользователя"
// @Success 200 {object} map[string]string "Успешная регистрация, возвращает статус и токен в заголовке"
// @Failure 400 {object} map[string]string "Ошибка в формате запроса"
// @Failure 409 {object} map[string]string "Пользователь с таким логином уже существует"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /register [post]
func (a *AuthController) HandleRegisterJSON(w http.ResponseWriter, r *http.Request) {
	credentials, err := paramsparser.JSONParse[domain.Credentials](w, r)
	if err != nil {
		return
	}
	a.AuthUseCase.Register(w, credentials)
}

func (a *AuthController) HandleLoginJSON(w http.ResponseWriter, r *http.Request) {
	credentials, err := paramsparser.JSONParse[domain.Credentials](w, r)
	if err != nil {
		return
	}
	a.AuthUseCase.Login(w, credentials)
}
