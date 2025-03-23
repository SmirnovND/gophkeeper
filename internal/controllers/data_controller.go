package controllers

import (
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"github.com/SmirnovND/toolbox/pkg/paramsparser"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// DataController контроллер для работы с данными пользователя
type DataController struct {
	dataUseCase interfaces.DataUseCase
}

// NewDataController создает новый экземпляр DataController
func NewDataController(dataUseCase interfaces.DataUseCase) *DataController {
	return &DataController{
		dataUseCase: dataUseCase,
	}
}

// SaveCredential сохраняет учетные данные (логин/пароль)
// @Summary Сохранить учетные данные
// @Description Сохраняет пару логин/пароль пользователя
// @Tags data
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer токен"
// @Param label path string true "Метка для идентификации данных"
// @Param credential body domain.CredentialData true "Учетные данные"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/data/credential/{label} [post]
func (c *DataController) SaveCredential(w http.ResponseWriter, r *http.Request) {
	// Получаем метку из URL
	label := chi.URLParam(r, "label")
	if label == "" {
		http.Error(w, "метка не предоставлена", http.StatusBadRequest)
		return
	}

	// Получаем данные из тела запроса
	credentialData, err := paramsparser.JSONParse[domain.CredentialData](w, r)
	if err != nil {
		return
	}

	c.dataUseCase.SaveCredential(w, r, label, credentialData)
}

// GetCredential получает учетные данные (логин/пароль)
// @Summary Получить учетные данные
// @Description Получает пару логин/пароль пользователя по метке
// @Tags data
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer токен"
// @Param label path string true "Метка для идентификации данных"
// @Success 200 {object} domain.CredentialData
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/data/credential/{label} [get]
func (c *DataController) GetCredential(w http.ResponseWriter, r *http.Request) {
	// Получаем метку из URL
	label := chi.URLParam(r, "label")
	if label == "" {
		http.Error(w, "метка не предоставлена", http.StatusBadRequest)
		return
	}

	c.dataUseCase.GetCredential(w, r, label)
}

// SaveCard сохраняет данные кредитной карты
// @Summary Сохранить данные кредитной карты
// @Description Сохраняет данные кредитной карты пользователя
// @Tags data
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer токен"
// @Param label path string true "Метка для идентификации данных"
// @Param card body domain.CardData true "Данные кредитной карты"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/data/card/{label} [post]
func (c *DataController) SaveCard(w http.ResponseWriter, r *http.Request) {
	// Получаем метку из URL
	label := chi.URLParam(r, "label")
	if label == "" {
		http.Error(w, "метка не предоставлена", http.StatusBadRequest)
		return
	}

	// Получаем данные из тела запроса
	cardData, err := paramsparser.JSONParse[domain.CardData](w, r)
	if err != nil {
		return
	}

	c.dataUseCase.SaveCard(w, r, label, cardData)
}

// GetCard получает данные кредитной карты
// @Summary Получить данные кредитной карты
// @Description Получает данные кредитной карты пользователя по метке
// @Tags data
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer токен"
// @Param label path string true "Метка для идентификации данных"
// @Success 200 {object} domain.CardData
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/data/card/{label} [get]
func (c *DataController) GetCard(w http.ResponseWriter, r *http.Request) {
	// Получаем метку из URL
	label := chi.URLParam(r, "label")
	if label == "" {
		http.Error(w, "метка не предоставлена", http.StatusBadRequest)
		return
	}

	c.dataUseCase.GetCard(w, r, label)
}

// SaveText сохраняет текстовые данные
// @Summary Сохранить текстовые данные
// @Description Сохраняет произвольный текст пользователя
// @Tags data
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer токен"
// @Param label path string true "Метка для идентификации данных"
// @Param text body domain.TextData true "Текстовые данные"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/data/text/{label} [post]
func (c *DataController) SaveText(w http.ResponseWriter, r *http.Request) {
	// Получаем метку из URL
	label := chi.URLParam(r, "label")
	if label == "" {
		http.Error(w, "метка не предоставлена", http.StatusBadRequest)
		return
	}

	// Получаем данные из тела запроса
	textData, err := paramsparser.JSONParse[domain.TextData](w, r)
	if err != nil {
		return
	}

	c.dataUseCase.SaveText(w, r, label, textData)
}

// GetText получает текстовые данные
// @Summary Получить текстовые данные
// @Description Получает произвольный текст пользователя по метке
// @Tags data
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer токен"
// @Param label path string true "Метка для идентификации данных"
// @Success 200 {object} domain.TextData
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/data/text/{label} [get]
func (c *DataController) GetText(w http.ResponseWriter, r *http.Request) {
	// Получаем метку из URL
	label := chi.URLParam(r, "label")
	if label == "" {
		http.Error(w, "метка не предоставлена", http.StatusBadRequest)
		return
	}

	c.dataUseCase.GetText(w, r, label)
}
