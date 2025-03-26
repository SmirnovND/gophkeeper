package usecase

import (
	"encoding/json"
	"github.com/SmirnovND/gophkeeper/internal/domain"
	"github.com/SmirnovND/gophkeeper/internal/interfaces"
	"net/http"
)

type DataUseCase struct {
	dataService interfaces.DataService
	jwtService  interfaces.JwtService
}

func NewDataUseCase(
	dataService interfaces.DataService,
	jwtService interfaces.JwtService,
) interfaces.DataUseCase {
	return &DataUseCase{
		dataService: dataService,
		jwtService:  jwtService,
	}
}
func (c *DataUseCase) GetCredential(w http.ResponseWriter, r *http.Request, label string) {
	login, err := c.jwtService.ExtractLoginFromToken(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Ошибка получения логина: "+err.Error(), http.StatusInternalServerError)
		return
	}

	credentialData, metadata, err := c.dataService.GetCredential(login, label)
	if err != nil {
		if err.Error() == "учетные данные не найдены" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Создаем структуру ответа с метаинформацией
	response := struct {
		CredentialData *domain.CredentialData `json:"credential_data"`
		Metadata       string                 `json:"metadata"`
	}{
		CredentialData: credentialData,
		Metadata:       metadata,
	}

	// Отправляем данные в ответе
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (c *DataUseCase) SaveCredential(w http.ResponseWriter, r *http.Request, label string, credentialData *domain.CredentialData, metadata string) {
	login, err := c.jwtService.ExtractLoginFromToken(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Ошибка получения логина: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Сохраняем данные
	err = c.dataService.SaveCredential(login, label, credentialData, metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "учетные данные успешно сохранены"})
}

func (c *DataUseCase) SaveCard(w http.ResponseWriter, r *http.Request, label string, cardData *domain.CardData, metadata string) {
	login, err := c.jwtService.ExtractLoginFromToken(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Ошибка получения логина: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Сохраняем данные
	err = c.dataService.SaveCard(login, label, cardData, metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "данные карты успешно сохранены"})
}

func (c *DataUseCase) GetCard(w http.ResponseWriter, r *http.Request, label string) {
	login, err := c.jwtService.ExtractLoginFromToken(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Ошибка получения логина: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем данные
	cardData, metadata, err := c.dataService.GetCard(login, label)
	if err != nil {
		if err.Error() == "данные карты не найдены" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Создаем структуру ответа с метаинформацией
	response := struct {
		CardData *domain.CardData `json:"card_data"`
		Metadata string           `json:"metadata"`
	}{
		CardData: cardData,
		Metadata: metadata,
	}

	// Отправляем данные в ответе
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (c *DataUseCase) SaveText(w http.ResponseWriter, r *http.Request, label string, textData *domain.TextData, metadata string) {
	login, err := c.jwtService.ExtractLoginFromToken(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Ошибка получения логина: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Сохраняем данные
	err = c.dataService.SaveText(login, label, textData, metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "текстовые данные успешно сохранены"})
}

func (c *DataUseCase) GetText(w http.ResponseWriter, r *http.Request, label string) {
	login, err := c.jwtService.ExtractLoginFromToken(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Ошибка получения логина: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем данные
	textData, metadata, err := c.dataService.GetText(login, label)
	if err != nil {
		if err.Error() == "текстовые данные не найдены" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Создаем структуру ответа с метаинформацией
	response := struct {
		TextData *domain.TextData `json:"text_data"`
		Metadata string           `json:"metadata"`
	}{
		TextData: textData,
		Metadata: metadata,
	}

	// Отправляем данные в ответе
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (c *DataUseCase) DeleteCredential(w http.ResponseWriter, r *http.Request, label string) {
	login, err := c.jwtService.ExtractLoginFromToken(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Ошибка получения логина: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Удаляем данные
	err = c.dataService.DeleteCredential(login, label)
	if err != nil {
		if err.Error() == "учетные данные не найдены" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "учетные данные успешно удалены"})
}

func (c *DataUseCase) DeleteCard(w http.ResponseWriter, r *http.Request, label string) {
	login, err := c.jwtService.ExtractLoginFromToken(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Ошибка получения логина: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Удаляем данные
	err = c.dataService.DeleteCard(login, label)
	if err != nil {
		if err.Error() == "данные карты не найдены" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "данные карты успешно удалены"})
}

func (c *DataUseCase) DeleteText(w http.ResponseWriter, r *http.Request, label string) {
	login, err := c.jwtService.ExtractLoginFromToken(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Ошибка получения логина: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Удаляем данные
	err = c.dataService.DeleteText(login, label)
	if err != nil {
		if err.Error() == "текстовые данные не найдены" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "текстовые данные успешно удалены"})
}
