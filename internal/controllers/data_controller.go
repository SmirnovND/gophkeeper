package controllers

import (
	"encoding/json"
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
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
// @Param credential body object true "Учетные данные с метаинформацией"
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
	var requestData struct {
		CredentialData *domain.CredentialData `json:"credential_data"`
		Metadata       string                 `json:"metadata"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestData); err != nil {
		http.Error(w, "ошибка при декодировании JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if requestData.CredentialData == nil {
		http.Error(w, "учетные данные не предоставлены", http.StatusBadRequest)
		return
	}

	c.dataUseCase.SaveCredential(w, r, label, requestData.CredentialData, requestData.Metadata)
}

// GetCredential получает учетные данные (логин/пароль)
// @Summary Получить учетные данные
// @Description Получает пару логин/пароль пользователя по метке
// @Tags data
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer токен"
// @Param label path string true "Метка для идентификации данных"
// @Success 200 {object} object "Учетные данные с метаинформацией"
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
// @Param card body object true "Данные кредитной карты с метаинформацией"
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
	var requestData struct {
		CardData *domain.CardData `json:"card_data"`
		Metadata string           `json:"metadata"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestData); err != nil {
		http.Error(w, "ошибка при декодировании JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if requestData.CardData == nil {
		http.Error(w, "данные карты не предоставлены", http.StatusBadRequest)
		return
	}

	c.dataUseCase.SaveCard(w, r, label, requestData.CardData, requestData.Metadata)
}

// GetCard получает данные кредитной карты
// @Summary Получить данные кредитной карты
// @Description Получает данные кредитной карты пользователя по метке
// @Tags data
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer токен"
// @Param label path string true "Метка для идентификации данных"
// @Success 200 {object} object "Данные кредитной карты с метаинформацией"
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
// @Param text body object true "Текстовые данные с метаинформацией"
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
	var requestData struct {
		TextData *domain.TextData `json:"text_data"`
		Metadata string           `json:"metadata"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestData); err != nil {
		http.Error(w, "ошибка при декодировании JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if requestData.TextData == nil {
		http.Error(w, "текстовые данные не предоставлены", http.StatusBadRequest)
		return
	}

	c.dataUseCase.SaveText(w, r, label, requestData.TextData, requestData.Metadata)
}

// GetText получает текстовые данные
// @Summary Получить текстовые данные
// @Description Получает произвольный текст пользователя по метке
// @Tags data
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer токен"
// @Param label path string true "Метка для идентификации данных"
// @Success 200 {object} object "Текстовые данные с метаинформацией"
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

// DeleteCredential удаляет учетные данные (логин/пароль)
// @Summary Удалить учетные данные
// @Description Удаляет пару логин/пароль пользователя по метке
// @Tags data
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer токен"
// @Param label path string true "Метка для идентификации данных"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/data/credential/{label} [delete]
func (c *DataController) DeleteCredential(w http.ResponseWriter, r *http.Request) {
	// Получаем метку из URL
	label := chi.URLParam(r, "label")
	if label == "" {
		http.Error(w, "метка не предоставлена", http.StatusBadRequest)
		return
	}

	c.dataUseCase.DeleteCredential(w, r, label)
}

// DeleteCard удаляет данные кредитной карты
// @Summary Удалить данные кредитной карты
// @Description Удаляет данные кредитной карты пользователя по метке
// @Tags data
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer токен"
// @Param label path string true "Метка для идентификации данных"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/data/card/{label} [delete]
func (c *DataController) DeleteCard(w http.ResponseWriter, r *http.Request) {
	// Получаем метку из URL
	label := chi.URLParam(r, "label")
	if label == "" {
		http.Error(w, "метка не предоставлена", http.StatusBadRequest)
		return
	}

	c.dataUseCase.DeleteCard(w, r, label)
}

// DeleteText удаляет текстовые данные
// @Summary Удалить текстовые данные
// @Description Удаляет произвольный текст пользователя по метке
// @Tags data
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer токен"
// @Param label path string true "Метка для идентификации данных"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/data/text/{label} [delete]
func (c *DataController) DeleteText(w http.ResponseWriter, r *http.Request) {
	// Получаем метку из URL
	label := chi.URLParam(r, "label")
	if label == "" {
		http.Error(w, "метка не предоставлена", http.StatusBadRequest)
		return
	}

	c.dataUseCase.DeleteText(w, r, label)
}
