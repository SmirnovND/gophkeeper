package controllers

import (
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"github.com/SmirnovND/toolbox/pkg/paramsparser"
	"net/http"
)

type AuthController struct {
	AuthUseCase interfaces.AuthUseCase
}

func NewAuthController(AuthUseCase interfaces.AuthUseCase) *AuthController {
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
// @Router /api/user/register [post]
func (a *AuthController) HandleRegisterJSON(w http.ResponseWriter, r *http.Request) {
	credentials, err := paramsparser.JSONParse[domain.Credentials](w, r)
	if err != nil {
		return
	}
	a.AuthUseCase.Register(w, credentials)
}

// HandleLoginJSON godoc
// @Summary Авторизация пользователя
// @Description Авторизует пользователя в системе и возвращает токен авторизации
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body domain.Credentials true "Учетные данные пользователя"
// @Success 200 {object} map[string]string "Успешная авторизация, возвращает статус и токен в заголовке"
// @Failure 400 {object} map[string]string "Ошибка в формате запроса"
// @Failure 401 {object} map[string]string "Неверный логин или пароль"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/user/login [post]
func (a *AuthController) HandleLoginJSON(w http.ResponseWriter, r *http.Request) {
	credentials, err := paramsparser.JSONParse[domain.Credentials](w, r)
	if err != nil {
		return
	}
	a.AuthUseCase.Login(w, credentials)
}
